package handler

import (
	"net/http"

	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/service"
)

func CheckMysql(service service.CheckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := service.CheckMysql()
		if err != nil {
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, resp, nil)
	}
}
