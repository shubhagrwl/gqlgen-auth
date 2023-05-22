package jwt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"todo/internal/app/service/graph/model"
	"todo/internal/app/service/logger"

	userDBmodel "todo/internal/app/db/dto/user"
	appContext "todo/internal/app/util/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

type IJwtService interface {
	CreateToken(ctx context.Context, user userDBmodel.UserDetail) (*model.TokenDetails, error)
}

type JWTServiceImpl struct {
}

func NewJWT() *JWTServiceImpl {
	return &JWTServiceImpl{}
}

var JWTAccessSigningKey = "TODO_USER_ACCESS_SECRET"
var JWTRefreshSigningKey = "TODO_USER_REFRESH_SECRET"

func (j *JWTServiceImpl) CreateToken(ctx context.Context, user userDBmodel.UserDetail) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}

	atExpires := int(time.Now().Add(time.Minute * 15).Unix())
	td.AtExpires = &atExpires
	newUUID := uuid.New().String()
	rtExpires := int(time.Now().Add(time.Hour * 24 * 7).Unix())

	td.AccessUUID = &newUUID
	td.RtExpires = &rtExpires
	refreshUUID := *td.AccessUUID + "++" + user.ID.String()
	td.RefreshUUID = &refreshUUID

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = user.ID
	atClaims["email"] = user.Email
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(os.Getenv(JWTAccessSigningKey)))
	if err != nil {
		return nil, err
	}

	td.AccessToken = &accessToken

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(os.Getenv(JWTRefreshSigningKey)))
	if err != nil {
		return nil, err
	}

	td.RefreshToken = &refreshToken
	return td, nil
}

func (j *JWTServiceImpl) Auth(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		log := logger.Logger(ctx)
		if authHeader != "" {
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 {
				tokenString := tokenParts[1]
				if len(tokenString) == 0 {
					log.Error("unauthorized")
					return
				}

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(os.Getenv(JWTAccessSigningKey)), nil
				})
				if err != nil {
					log.Error("unauthorized")
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if ok && !token.Valid {
					log.Error("unauthorized")
					return

				}

				userID, ok := claims["user_id"].(string)
				if !ok {
					log.Error("No user id present in context...")
					return
				}

				ctx := appContext.WithUserID(c.Request.Context(), userID)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
			}
		}
	}
}

type JwtKey struct {
	Claims      jwt.MapClaims
	TokenString string
}

// ParseJWT method will parse the token string and extract claims fields
func (j *JwtKey) ParseJWT() error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(j.TokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTAccessSigningKey), nil
	})
	if !token.Valid || err != nil {
		return err
	}
	j.Claims = claims
	return nil
}
