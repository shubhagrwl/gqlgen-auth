package user

import (
	"context"
	"fmt"
	jwt "todo/internal/app/api/middleware"
	"todo/internal/app/db/repository/user"
	"todo/internal/app/service/graph/model"
	"todo/internal/app/service/logger"
	"todo/internal/app/util"
	appContext "todo/internal/app/util/context"
)

type IUserController interface {
	GetProfile(ctx context.Context) (*model.User, error)
}

type UserController struct {
	User user.IUserRepository
	JWT  jwt.IJwtService
}

func NewUserController(
	userDBClient user.IUserRepository,
	authService *jwt.JWTServiceImpl) IUserController {
	return &UserController{
		User: userDBClient,
		JWT:  authService,
	}
}

func (s *UserController) GetProfile(ctx context.Context) (*model.User, error) {
	log := logger.Logger(ctx)
	log.Info("get profile")

	userID, err := appContext.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := s.User.GetUserWithId(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return util.UserDetailToUser(user), nil

}
