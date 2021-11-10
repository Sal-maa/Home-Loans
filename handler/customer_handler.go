package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetCheckApply(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService := customerService.SCGetCheckApply(resp.IdUser)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreateIdentity(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}
		payloads, _ := ioutil.ReadAll(r.Body)

		var identity contract.Identity

		json.Unmarshal(payloads, &identity)

		validate := validator.New()
		error := validate.Struct(identity)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService, err := customerService.SCCreateIdentity(&identity, resp.IdUser)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreateSubmission(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)

		var submission contract.Submission
		json.Unmarshal(payloads, &submission)

		validate := validator.New()
		error := validate.Struct(submission)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService := customerService.SCCreateSubmission(&submission, resp.IdUser)
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmissionCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService, err := customerService.SCGetSubmission(resp.IdUser)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmissionStatus(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService := customerService.SCGetSubmissionStatus(resp.IdUser)
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func UploadFileKtp(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("ktp")

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		defer file.Close()

		dataService, err := customerService.SCUploadFileKTP(&file, handler, resp)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func UploadFileGaji(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("bukti_gaji")

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		defer file.Close()

		dataService, err := customerService.SCUploadFileGaji(&file, handler, resp)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func UploadFilePendukung(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("dokumen_pendukung")

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		defer file.Close()

		dataService, err := customerService.SCUploadFilePendukung(&file, handler, resp)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetFileKtpCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 1 && resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiKtp := vars["bukti_ktp"]

		dataService := customerService.SCGetFileKtpCustomer(buktiKtp)

		data, readErr := ioutil.ReadAll(dataService)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object ")
			return
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
		// responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetFileBuktiGajiCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 1 && resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiGaji := vars["bukti_gaji"]

		dataService := customerService.SCGetFileBuktiGajiCustomer(buktiGaji)

		data, readErr := ioutil.ReadAll(dataService)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object ")
			return
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
		// responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetFilePendukungCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 1 && resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}
		vars := mux.Vars(r)
		buktiFilependukung := vars["dokumen_pendukung"]

		dataService := customerService.SCGetFilePendukungCustomer(buktiFilependukung)

		data, readErr := ioutil.ReadAll(dataService)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object ")
			return
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
		// responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetIdentityCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService, err := customerService.SCGetIdentity(resp.IdUser)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func UpdateIdentityCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 1 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		payloads, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}
		var updateIdentity contract.Identity
		json.Unmarshal(payloads, &updateIdentity)

		validate := validator.New()
		error := validate.Struct(updateIdentity)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService, err := customerService.SCUpdateIdentityCustomer(&updateIdentity, resp.IdUser)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}
