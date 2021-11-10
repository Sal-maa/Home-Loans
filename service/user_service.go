package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/jwt_client"
	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type UserServiceInterface interface {
	SULogin(user *contract.User) (*contract.GetTokenResponseContract, error)
	SUCreate(user *contract.User) interface{}
	GetToken(*contract.User) (*contract.GetTokenResponseContract, error)
}

func NewUserService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *userService {
	return &userService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *userService) SUCreate(user *contract.User) interface{} {

	user.LoginAs = 1

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	password := user.Password
	hash, _ := HashPassword(password)
	user.Password = hash
	db.DbConnection.Create(&user)

	uReturn := contract.UserReturn{
		Username: user.Username,
		LoginAs:  user.LoginAs,
	}
	return &uReturn
}

func (s *userService) SULogin(user *contract.User) (*contract.GetTokenResponseContract, error) {
	var registeredUser *contract.User

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("users").First(&registeredUser, "username = ?", user.Username).Error

	if err != nil {
		return nil, errors.NewUnauthorizedError("error when accessing database")
	}

	if user.Username != registeredUser.Username {
		return nil, errors.NewUnauthorizedError("combination of username and password not match, username tidak ada")
	}

	if !CheckPasswordHash(user.Password, registeredUser.Password) {
		return nil, errors.NewUnauthorizedError("combination of username and password not match, password salah")
	}

	tk, err := s.GetToken(registeredUser)

	if err != nil {
		errMsg := fmt.Sprintf("error di get token: %v", err)
		log.Errorf(errMsg)
		return nil, errors.NewInternalError(err, errMsg)
	}

	return tk, nil
}

func (s *userService) GetToken(user *contract.User) (*contract.GetTokenResponseContract, error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	atClaims := contract.JWTMapClaim{
		Authorized: true,
		RequestID:  uuid.New().String(),
		IdUser:     user.ID,
		Username:   user.Username,
		LoginAs:    user.LoginAs,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := s.jwtClient.GenerateTokenStringWithClaims(atClaims, s.appConfig.JWTSecret)

	if err != nil {
		errMsg := fmt.Sprintf("error signed JWT credentials: %v", err)
		log.Errorf(errMsg)
		return nil, errors.NewInternalError(err, errMsg)
	}

	return &contract.GetTokenResponseContract{Token: token}, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
