package auth

import (
	"context"
	"errors"
	"fmt"
	jwt "todo/internal/app/api/middleware"
	"todo/internal/app/constants"
	onboardDbmodel "todo/internal/app/db/dto/onboard_user"

	"time"
	"todo/internal/app/db/repository/user"
	"todo/internal/app/service/graph/model"
	"todo/internal/app/service/logger"
	"todo/internal/app/util"

	"todo/internal/app/util/totp"

	"github.com/jinzhu/gorm"
)

type IAuthController interface {
	SignUp(ctx context.Context, data model.UserInput) (*model.LoginResponse, error)
	SignIn(ctx context.Context, data model.Login) (*model.LoginResponse, error)
	SendCode(ctx context.Context, data model.SendCode) (*model.Response, error)
	VerifyCode(ctx context.Context, data model.Code) (*model.Success, error)
	ResetPassword(ctx context.Context, data model.ResetPassword) (*model.Success, error)
}

type AuthUpController struct {
	User user.IUserRepository
	JWT  jwt.IJwtService
	TOTP totp.ITOTP
}

func NewAuthController(
	userDBClient user.IUserRepository,
	authService *jwt.JWTServiceImpl,
	totpService *totp.TOTP) IAuthController {
	return &AuthUpController{
		User: userDBClient,
		JWT:  authService,
		TOTP: totpService,
	}
}

func (s *AuthUpController) SignUp(ctx context.Context, data model.UserInput) (*model.LoginResponse, error) {
	log := logger.Logger(ctx)
	var response model.LoginResponse

	_, found, err := s.User.GetOnboardUser(ctx, data.Email, util.Bool(true))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("internal server error")
	}

	if !found {
		user, err := s.User.SignUp(ctx, data)
		if err != nil {
			log.Error(err.Error())
			return &response, err
		}

		token, err := s.JWT.CreateToken(ctx, user)
		if err != nil {
			log.Error(err.Error())
			return &response, err
		}

		response.User = util.UserDetailToUser(user)
		response.JwtToken = token

		return &response, nil
	}

	return nil, fmt.Errorf("internal server error")
}

func (s *AuthUpController) SignIn(ctx context.Context, data model.Login) (*model.LoginResponse, error) {
	log := logger.Logger(ctx)
	var response model.LoginResponse

	user, err := s.User.GetUserWithEmail(ctx, *data.Email)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if !util.ValidatePassword(ctx, *data.Password, *user.Password) {
		log.Error(err.Error())
		return nil, err
	}
	token, err := s.JWT.CreateToken(ctx, user)
	if err != nil {
		log.Error(err.Error())
		return &response, err
	}

	response.User = util.UserDetailToUser(user)
	response.JwtToken = token

	return &response, nil
}

func (s *AuthUpController) SendCode(ctx context.Context, data model.SendCode) (*model.Response, error) {
	log := logger.Logger(ctx)

	user, found, err := s.User.GetOnboardUser(ctx, data.Email, nil)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err)
			return nil, fmt.Errorf("internal server error")
		}
	}

	var validator string
	if data.Service == model.SendCodeServiceSignUp {

		if found && *user.Verified {
			_, err := s.User.GetUserWithEmail(ctx, data.Email)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					var patch = make(map[string]interface{})
					patch[onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__Verified] = false

					err := s.User.UpdateOnboardUser(ctx, data.Email, patch)
					if err != nil {
						log.Error(err)
						return nil, fmt.Errorf("internal server error")
					}
				} else {
					log.Error("user already have an account")
					return nil, fmt.Errorf("user already have an account")
				}
			}

		} else {
			err := s.User.CreateOnboardUser(ctx, onboardDbmodel.OnboardUser{
				Email:    &data.Email,
				Verified: util.Bool(false),
			})
			if err != nil {
				log.Error("user already have an account")
				return nil, fmt.Errorf("user already have an account")
			}
		}
		validator = model.SendCodeServiceSignUp.String()
	} else {
		validator = model.SendCodeServiceForgetPassword.String()
	}

	code, err := s.TOTP.GetCode(ctx, data.Email+validator, &constants.CODEVALIDITY)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	err = util.SendMail(ctx, "Account Verification", "", data.Email, constants.TEMPLATE_ID__ACCOUNT_VERIFICATION, code)
	if err != nil {
		log.Error("unable to send email")
		return nil, fmt.Errorf("unable to send email")
	}

	return &model.Response{Success: true, Message: "code send successfully"}, nil

}

func (s *AuthUpController) VerifyCode(ctx context.Context, data model.Code) (*model.Success, error) {
	log := logger.Logger(ctx)

	user, found, err := s.User.GetOnboardUser(ctx, data.Email, util.Bool(false))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("internal server error")
	}
	if found {
		valid, err := s.TOTP.VerifyCode(ctx, data.Email+data.Service.String(), data.Code, nil)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		if valid {
			var patch = make(map[string]interface{})
			patch[onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__Verified] = true
			patch[onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__VerifiedAt] = time.Now()
			err = s.User.UpdateOnboardUser(ctx, *user.Email, patch)
			if err != nil {
				log.Error(err)
				return nil, fmt.Errorf("internal server error")
			}
		}

		return &model.Success{Success: valid}, nil

	}
	return &model.Success{Success: false}, fmt.Errorf("user not found")
}

func (s *AuthUpController) ResetPassword(ctx context.Context, data model.ResetPassword) (*model.Success, error) {
	log := logger.Logger(ctx)

	user, err := s.User.GetUserWithEmail(ctx, data.Email)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	valid, err := s.TOTP.VerifyCode(ctx, data.Email+model.SendCodeServiceForgetPassword.String(), data.Code, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if valid {
		encryptedPass, err := util.GenerateHash(ctx, data.Password)
		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf("internal server error")
		}

		err = s.User.UpdatePassword(ctx, *user.Email, encryptedPass)
		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf("internal server error")
		}
		return &model.Success{Success: true}, nil

	}

	return &model.Success{Success: false}, fmt.Errorf("invalid code")
}

// func (s *AuthUpController) ForgetPassword(ctx context.Context, data model.ForgetPassword) (*model.Success, error) {
// 	log := logger.Logger(ctx)
// 	log.Info("get profile")
// 	userCtx, err := appContext.UserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unauthorized")
// 	}

// 	user, err := s.User.GetUserWithEmail(ctx, *userCtx.Email)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return nil, err
// 	}

// 	return util.UserDetailToUser(user), nil
// }
