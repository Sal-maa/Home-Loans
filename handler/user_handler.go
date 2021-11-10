package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func Create(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

		payloads, _ := ioutil.ReadAll(r.Body)

		var user contract.User
		json.Unmarshal(payloads, &user)

		validate := validator.New()
		error := validate.Struct(user)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService := userService.SUCreate(&user)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func Login(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		payloads, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		var user contract.User
		json.Unmarshal(payloads, &user)

		validte := validator.New()
		error := validte.Struct(user)
		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService, err := userService.SULogin(&user)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 1)

		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   dataService.Token,
				Expires: expirationTime,
			})

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}
