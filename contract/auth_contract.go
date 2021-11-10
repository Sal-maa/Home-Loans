package contract

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/common/util"
	log "github.com/sirupsen/logrus"
)

type JWTMapClaim struct {
	Authorized bool   `json:"authorized"`
	RequestID  string `json:"requestID"`
	IdUser     uint   `json:"id_user"`
	Username   string `json:"username"`
	LoginAs    uint   `json:"login_as"`
	jwt.StandardClaims
}

type GetTokenResponseContract struct {
	Token string `json:"token"`
}

type ValidateTokenRequestContract struct {
	Token string `json:"token"`
}

func NewValidateTokenRequest(r *http.Request) (*ValidateTokenRequestContract, error) {
	validateTokenContract := new(ValidateTokenRequestContract)
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(validateTokenContract); err != nil {
		log.Warning(err)
		return nil, errors.NewBadRequestError(err)
	}

	validate := validator.New()
	util.UseJsonFieldValidation(validate)

	if err := validate.Struct(validate); err != nil {
		log.Error(err)
		return nil, errors.NewValidationError(errors.ValidateErrToMapString(err.(validator.ValidationErrors)))
	}

	return validateTokenContract, nil
}

func NewValidateTokenRequestViaCookie(r *http.Request) (*ValidateTokenRequestContract, error) {
	validateTokenContract := new(ValidateTokenRequestContract)
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Warning(err)
			return nil, errors.NewUnauthorizedError("error no token in cookie")
		}
		return nil, errors.NewBadRequestError(err)
	}

	validateTokenContract.Token = cookie.Value

	validate := validator.New()
	util.UseJsonFieldValidation(validate)

	if err := validate.Struct(validate); err != nil {
		log.Error(err)
		return nil, errors.NewValidationError(errors.ValidateErrToMapString(err.(validator.ValidationErrors)))
	}

	return validateTokenContract, nil
}
