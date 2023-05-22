package service

import (
	"context"
	"os"
	jwt "todo/internal/app/api/middleware"
	"todo/internal/app/constants"
	"todo/internal/app/db"
	"todo/internal/app/db/repository/user"
	"todo/internal/app/mutation/auth"
	userController "todo/internal/app/mutation/user"

	"todo/internal/app/service/logger"
	"todo/internal/app/util/totp"
	"todo/internal/config"
)

// Services is the interface for enclosing all the services
type Services interface {
	Auth() auth.IAuthController
	User() userController.IUserController
}

type services struct {
	authService auth.IAuthController
	userService userController.IUserController
}

func (svc *services) Auth() auth.IAuthController {
	return svc.authService
}

func (svc *services) User() userController.IUserController {
	return svc.userService
}

func Init() Services {
	serviceName := "user-management-parser"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "dev"
	}
	config.Init(serviceName, environment)
	logger.InitLogger()

	db.Init()

	jwt := jwt.NewJWT()
	user := user.NewUserRepository(jwt)
	totpService := totp.New(context.Background(), constants.ProjectEnvironment(environment))

	return &services{
		userService: userController.NewUserController(user, jwt),
		authService: auth.NewAuthController(user, jwt, totpService),
	}
}
