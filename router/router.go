package router

import (
	"net/http"
	"os"

	"github.com/rysmaadit/go-template/handler"
	"github.com/rysmaadit/go-template/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(dependencies service.Dependencies) http.Handler {
	r := mux.NewRouter()

	setHomeRouter(r)
	setAuthRouter(r, dependencies.AuthService)
	setCheckRouter(r, dependencies.CheckService)
	setUserRouter(r, dependencies.UserService)
	setCustomerRouter(r, dependencies.CustomerService)
	setEmployeeRouter(r, dependencies.EmployeeService)

	loggedRouter := handlers.LoggingHandler(os.Stderr, r)
	return loggedRouter
}

func setHomeRouter(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/").Handler(handler.Home())
}

func setAuthRouter(router *mux.Router, dependencies service.AuthServiceInterface) {
	router.Methods(http.MethodGet).Path("/auth/token").Handler(handler.GetToken(dependencies))
	router.Methods(http.MethodPost).Path("/auth/token/validate").Handler(handler.ValidateToken(dependencies))
}

func setCheckRouter(router *mux.Router, checkService service.CheckService) {
	router.Methods(http.MethodGet).Path("/check/mysql").Handler(handler.CheckMysql(checkService))
}

func setUserRouter(router *mux.Router, dependencies service.UserServiceInterface) {
	router.Methods(http.MethodPost).Path("/signup").Handler(handler.Create(dependencies))
	router.Methods(http.MethodPost).Path("/login").Handler(handler.Login(dependencies))
}

func setCustomerRouter(router *mux.Router, dependencies service.CustomerServiceInterface) {
	router.Methods(http.MethodGet).Path("/checkapply").Handler(handler.GetCheckApply(dependencies))
	router.Methods(http.MethodPost).Path("/createidentity").Handler(handler.CreateIdentity(dependencies))
	router.Methods(http.MethodPost).Path("/createsubmission").Handler(handler.CreateSubmission(dependencies))
	router.Methods(http.MethodGet).Path("/submission/getstatus").Handler(handler.GetSubmissionStatus(dependencies))
	router.Methods(http.MethodGet).Path("/getsubmission").Handler(handler.GetSubmissionCustomer(dependencies))
	router.Methods(http.MethodGet).Path("/getidentity").Handler(handler.GetIdentityCustomer(dependencies))
	router.Methods(http.MethodPost).Path("/uploadfilektp").Handler(handler.UploadFileKtp(dependencies))
	router.Methods(http.MethodPost).Path("/uploadfilegaji").Handler(handler.UploadFileGaji(dependencies))
	router.Methods(http.MethodPost).Path("/uploadfilependukung").Handler(handler.UploadFilePendukung(dependencies))
	router.Methods(http.MethodGet).Path("/cgetfilektp/{bukti_ktp}").Handler(handler.GetFileKtpCustomer(dependencies))
	router.Methods(http.MethodGet).Path("/cgetfilegaji/{bukti_gaji}").Handler(handler.GetFileBuktiGajiCustomer(dependencies))
	router.Methods(http.MethodGet).Path("/cgetfilependukung/{dokumen_pendukung}").Handler(handler.GetFilePendukungCustomer(dependencies))
	router.Methods(http.MethodPost).Path("/identity/update").Handler(handler.UpdateIdentityCustomer(dependencies))
}

func setEmployeeRouter(router *mux.Router, dependencies service.EmployeeServiceInterface) {
	router.Methods(http.MethodGet).Path("/numberofpage").Handler(handler.GetNumberOfPage(dependencies))
	router.Methods(http.MethodGet).Path("/listsubmission/{page}").Handler(handler.GetListSubmission(dependencies))
	router.Methods(http.MethodGet).Path("/searchbyname/{name}").Handler(handler.GetListByName(dependencies))
	router.Methods(http.MethodGet).Path("/submission/{id}").Handler(handler.GetSubmissionEmployee(dependencies))
	router.Methods(http.MethodPost).Path("/submission/status/{id_cust}").Handler(handler.PostSubmissionStatus(dependencies))
	router.Methods(http.MethodPost).Path("/identity/status/{id_cust}").Handler(handler.PostIdentityStatus(dependencies))
	router.Methods(http.MethodGet).Path("/pgetfilektp/{bukti_ktp}").Handler(handler.GetFileKtpEmployee(dependencies))
	router.Methods(http.MethodGet).Path("/pgetfilegaji/{bukti_gaji}").Handler(handler.GetFileBuktiGajiEmployee(dependencies))
	router.Methods(http.MethodGet).Path("/pgetfilependukung/{dokumen_pendukung}").Handler(handler.GetFilePendukungEmployee(dependencies))
	router.Methods(http.MethodGet).Path("/statustotal").Handler(handler.TotalIdentityUnconfirmed(dependencies))
	router.Methods(http.MethodGet).Path("/identitycustomer/{id_cust}").Handler(handler.GetIdentityEmployee(dependencies))
	router.Methods(http.MethodGet).Path("/downloadreport").Handler(handler.DownloadReport(dependencies))
	router.Methods(http.MethodGet).Path("/listsubmission").Handler(handler.GetListSubmissionParam(dependencies)).Queries("page", "{page}", "per_page", "{per_page}", "name", "{name}")
}
