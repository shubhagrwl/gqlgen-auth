package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"todo/internal/app/constants"
	db "todo/internal/app/db"
	onboardDbmodel "todo/internal/app/db/dto/onboard_user"
	userdDbmodel "todo/internal/app/db/dto/user"

	"todo/internal/app/service/graph/model"
	"todo/internal/app/service/logger"
	"todo/internal/app/util"

	jwt "todo/internal/app/api/middleware"
	utilTime "todo/internal/app/util/time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	SignUp(ctx context.Context, user model.UserInput) (userdDbmodel.UserDetail, error)
	GetUserWithEmail(ctx context.Context, email string) (userdDbmodel.UserDetail, error)
	GetUserWithId(ctx context.Context, id uuid.UUID) (userdDbmodel.UserDetail, error)
	UpdatePassword(ctx context.Context, email, password string) error

	GetOnboardUser(ctx context.Context, email string, status *bool) (onboardDbmodel.OnboardUser, bool, error)
	CreateOnboardUser(ctx context.Context, onboardUser onboardDbmodel.OnboardUser) error
	UpdateOnboardUser(ctx context.Context, email string, patch map[string]interface{}) error
}

type UserRepositoryImpl struct {
	JWT       jwt.IJwtService
	DBService db.DBService
}

func NewUserRepository(jwtService *jwt.JWTServiceImpl) IUserRepository {
	return &UserRepositoryImpl{
		JWT:       jwtService,
		DBService: db.DBService{},
	}
}

// //SignUp method adds a user details to the DB
func (u *UserRepositoryImpl) SignUp(ctx context.Context, user model.UserInput) (userdDbmodel.UserDetail, error) {
	log := logger.Logger(ctx)

	var usr userdDbmodel.UserDetail
	tx := u.DBService.GetDB().Begin()
	tx.LogMode(constants.DBLOGMODE)

	fullName := strings.Split(user.FullName, " ")
	hashedPassword, err := util.GenerateHash(ctx, user.Password)
	if err != nil {
		return usr, err
	}

	dob, err := utilTime.ParseTime(user.DateOfBirth)
	if err != nil {
		return usr, err
	}

	newID := uuid.New()
	usr = userdDbmodel.UserDetail{
		ID:        &newID,
		Password:  &hashedPassword,
		FirstName: &fullName[0],
		LastName:  &fullName[1],
		Email:     &user.Email,
		Dob:       &dob,
	}
	defer tx.Rollback()

	res := tx.Table(userdDbmodel.TABLEUSERDETAIL).
		Omit(userdDbmodel.TABLEUSERDETAIL_COLUMN__CreatedAt, userdDbmodel.TABLEUSERDETAIL_COLUMN__UpdatedAt).
		Create(&usr).Error
	if res != nil {
		log.Error(res)
		return usr, res
	}
	tx.Commit()

	return usr, nil
}

// GetUser method gets a user details to the DB
func (u *UserRepositoryImpl) GetUserWithEmail(ctx context.Context, email string) (userdDbmodel.UserDetail, error) {
	log := logger.Logger(ctx)

	var userDetailModel userdDbmodel.UserDetail
	tx := u.DBService.GetDB()
	tx.LogMode(constants.DBLOGMODE)

	result := tx.Table(userdDbmodel.TABLEUSERDETAIL).Where(userdDbmodel.TABLEUSERDETAIL_COLUMN__Email+" = ? ", email).First(&userDetailModel)
	err := result.Error
	if err != nil {
		log.Errorf("unable to get user", err)
		return userDetailModel, err
	}

	return userDetailModel, nil
}

// GetUser method gets a user details to the DB
func (u *UserRepositoryImpl) GetUserWithId(ctx context.Context, id uuid.UUID) (userdDbmodel.UserDetail, error) {
	log := logger.Logger(ctx)

	var userDetailModel userdDbmodel.UserDetail
	tx := u.DBService.GetDB()
	tx.LogMode(constants.DBLOGMODE)

	result := tx.Table(userdDbmodel.TABLEUSERDETAIL).Where(userdDbmodel.TABLEUSERDETAIL_COLUMN__ID+" = ? ", id).First(&userDetailModel)
	err := result.Error
	if err != nil {
		log.Errorf("unable to get user", err)
		return userDetailModel, err
	}

	return userDetailModel, nil
}

func (u *UserRepositoryImpl) UpdatePassword(ctx context.Context, email, password string) error {
	log := logger.Logger(ctx)
	tx := u.DBService.GetDB().Begin()
	tx.LogMode(constants.DBLOGMODE)
	defer tx.Rollback()

	if err := tx.Table(userdDbmodel.TABLEUSERDETAIL).
		Where(userdDbmodel.TABLEUSERDETAIL_COLUMN__Email+" = ? ", email).
		UpdateColumn(userdDbmodel.TABLEUSERDETAIL_COLUMN__Password, password); err.Error != nil {
		log.Errorf("error while updating password", err.Error)
		return err.Error
	}
	tx.Commit()
	return nil
}

func (u *UserRepositoryImpl) GetOnboardUser(ctx context.Context, email string, status *bool) (onboardDbmodel.OnboardUser, bool, error) {
	var onboardUserModel onboardDbmodel.OnboardUser
	tx := u.DBService.GetDB()
	tx.LogMode(constants.DBLOGMODE)

	whr := fmt.Sprintf("%s='%s' AND ", onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__Email, email)
	if status != nil {
		whr = fmt.Sprintf("%s %s='%s' AND ", whr, onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__Verified, strconv.FormatBool(*status))
	}
	whr = strings.TrimSuffix(whr, " AND ")
	err := tx.Table(onboardDbmodel.TABLE_ONBOARD_USER).
		Where(whr).
		First(&onboardUserModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return onboardUserModel, false, nil
		}
		return onboardUserModel, false, err
	}

	return onboardUserModel, true, nil
}

func (u *UserRepositoryImpl) CreateOnboardUser(ctx context.Context, onboardUser onboardDbmodel.OnboardUser) error {
	tx := u.DBService.GetDB().Begin()
	defer tx.Rollback()
	tx.LogMode(constants.DBLOGMODE)

	err := tx.Table(onboardDbmodel.TABLE_ONBOARD_USER).Omit(onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__ID).Create(&onboardUser).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *UserRepositoryImpl) UpdateOnboardUser(ctx context.Context, email string, patch map[string]interface{}) error {
	tx := u.DBService.GetDB().Begin()
	defer tx.Rollback()
	tx.LogMode(constants.DBLOGMODE)

	err := tx.Table(onboardDbmodel.TABLE_ONBOARD_USER).
		Where(onboardDbmodel.TABLE_ONBOARD_USER_COLUMN__Email+" = ?", email).
		UpdateColumns(patch).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
