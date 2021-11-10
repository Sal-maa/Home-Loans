package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetListSubmission(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]

		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		pages, err := strconv.Atoi(page)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetListSubmission(pages)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetNumberOfPage(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService := employeeService.SPGetNumberOfPage()

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetListByName(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService, err := employeeService.SPGetListByName(name)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmissionEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		custId := vars["id"]

		subIdint, err := strconv.Atoi(custId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetSubmission(uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func PostSubmissionStatus(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		// w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		subId := vars["id_cust"]

		subIdint, err := strconv.Atoi(subId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)
		var statusSubmission contract.Submission
		json.Unmarshal(payloads, &statusSubmission)

		dataService, err := employeeService.SPPostSubmissionStatus(&statusSubmission, uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func PostIdentityStatus(petugasService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// header := w.Header()
		// header.Add("Access-Control-Allow-Origin", "*")
		// header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		// header.Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := petugasService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		custId := vars["id_cust"]

		custIdint, err := strconv.Atoi(custId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)
		var statusDataDIri contract.Identity
		json.Unmarshal(payloads, &statusDataDIri)

		dataService, err := petugasService.SPPostIdentityStatus(&statusDataDIri, uint(custIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetFileKtpEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiKtp := vars["bukti_ktp"]

		dataService := employeeService.SPGetFileKtp(buktiKtp)

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
	}
}

func GetFileBuktiGajiEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiGaji := vars["bukti_gaji"]

		dataService := employeeService.SPGetFileBuktiGaji(buktiGaji)

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

	}
}

func GetFilePendukungEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiPendukung := vars["dokumen_pendukung"]

		dataService := employeeService.SPGetFileBuktiPendukung(buktiPendukung)

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
	}
}

func GetIdentityEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		subId := vars["id_cust"]

		subIdint, err := strconv.Atoi(subId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetIdentityEmployee(uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func TotalIdentityUnconfirmed(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService, _ := employeeService.SPGetStatusTotal()

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func DownloadReport(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		// vars := mux.Vars(r)
		// downloadReport := vars["download_report"]

		// dataService := employeeService.SPDownloadReport()
		f := employeeService.SPDownloadReport()

		// w.Header().Set("Content-Type", "application/excel")
		// w.WriteHeader(http.StatusOK)
		// w.Write(*dataService)
		// file, err := excelize.OpenFile("./Export.xlsx")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		w.Header().Set("Content-Disposition", "attachment; filename="+string("Report_pengajuan_yang_disetujui.xlsx")+"")
		f.Write(w)
	}
}

func GetListSubmissionParam(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]
		perPage := vars["per_page"]
		name := vars["name"]

		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}
		pages, err := strconv.Atoi(page)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		perPages, err := strconv.Atoi(perPage)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetListSubmissionParam(pages, perPages, name)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}
