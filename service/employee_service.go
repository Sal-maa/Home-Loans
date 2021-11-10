package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
	"github.com/minio/minio-go/v7"
	"github.com/rysmaadit/go-template/common/errors"
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/jwt_client"
	miniopkg "github.com/rysmaadit/go-template/external/minio"
	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
)

type employeeService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type EmployeeServiceInterface interface {
	SPGetListSubmission(page int) (*[]contract.ListSubmission, error)
	SPGetListByName(name string) (*[]contract.ListSubmission, error)
	SPGetNumberOfPage() *contract.NumberOfPage
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SPGetSubmission(id uint) (*contract.Submission, error)
	SPPostSubmissionStatus(submissionStatus *contract.Submission, id uint) (*contract.StatusKelengkapanReturn, error)
	SPPostIdentityStatus(statusPengajuan *contract.Identity, id uint) (*contract.StatusReturn, error)
	SPGetFileKtp(buktiKtp string) *minio.Object
	SPGetFileBuktiGaji(buktiGaji string) *minio.Object
	SPGetFileBuktiPendukung(buktiPendukung string) *minio.Object
	SPGetIdentityEmployee(id uint) (*contract.IdentityReturn, error)
	SPGetStatusTotal() (*contract.StatusTotalIdentity, error)
	SPDownloadReport() *excelize.File
	SPGetListSubmissionParam(page int, perPage int, name string) (*[]contract.ListSubmission, error)
}

func NewEmployeeService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *employeeService {
	return &employeeService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *employeeService) SPGetListSubmission(page int) (*[]contract.ListSubmission, error) {
	var ListSubmission []contract.ListSubmission
	var List []contract.ListAll

	namePersen := strings.Join([]string{"%", "", "%"}, "")

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Raw("SELECT identities.id, identities.updated_at, identities.deleted_at, identities.id_cust, identities.nik, identities.nama_lengkap, identities.tempat_lahir, identities.tanggal_lahir, identities.alamat, identities.pekerjaan, identities.pendapatan_perbulan, identities.bukti_ktp, identities.bukti_gaji, identities.status, submissions.created_at, submissions.alamat_rumah, submissions.luas_tanah, submissions.harga_rumah, submissions.jangka_pembayaran, submissions.dokumen_pendukung, submissions.status_kelengkapan FROM identities LEFT JOIN submissions ON identities.id_cust = submissions.id_cust WHERE identities.nama_lengkap LIKE ? LIMIT ? offset ?", namePersen, 5, (page-1)*5).Find(&List).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i, v := range List {
		if v.StatusKelengkapan == "" {
			tanggal := v.UpdatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      "-",
			}
			ListSubmission = append(ListSubmission, lsubmission)
		} else {
			tanggal := v.CreatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.StatusKelengkapan,
				Rekomendasi:      Recommendation(&v.PendapatanPerbulan, &v.HargaRumah, &v.JangkaPembayaran),
			}
			ListSubmission = append(ListSubmission, lsubmission)
		}
	}
	return &ListSubmission, nil
}

func Recommendation(PendapatanPerbulan, HargaRumah *float64, JangkaPembayaran *uint) string {
	var kemampuanCicilanPerbulan float64
	var kenyataanCicilanPerbulan float64
	kemampuanCicilanPerbulan = (*PendapatanPerbulan / 3)
	kenyataanCicilanPerbulan = (*HargaRumah / float64(*JangkaPembayaran)) / 12
	if kemampuanCicilanPerbulan > kenyataanCicilanPerbulan {
		return "Boleh"
	}
	return "Tidak Boleh"
}

func (s *employeeService) VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error) {
	claims := jwt.MapClaims{}

	err := s.jwtClient.ParseTokenWithClaims(req.Token, claims, s.appConfig.JWTSecret)

	if err != nil {
		log.Errorln(err)
		return nil, errors.NewUnauthorizedError("invalid parse token with claims")
	}

	authorized := fmt.Sprintf("%v", claims["authorized"])
	requestID := fmt.Sprintf("%v", claims["requestID"])

	if authorized == "" || requestID == "" {
		return nil, errors.NewUnauthorizedError("invalid payload")
	}

	ok, err := strconv.ParseBool(authorized)

	if err != nil || !ok {
		log.Errorln(err)
		return nil, errors.NewUnauthorizedError("invalid payload")
	}

	id_user_uint := uint(claims["id_user"].(float64))
	login_as_uint := uint(claims["login_as"].(float64))

	resp := &contract.JWTMapClaim{
		Authorized:     claims["authorized"].(bool),
		RequestID:      claims["requestID"].(string),
		IdUser:         id_user_uint,
		Username:       claims["username"].(string),
		LoginAs:        login_as_uint,
		StandardClaims: jwt.StandardClaims{},
	}

	return resp, nil
}

func (s *employeeService) SPGetNumberOfPage() *contract.NumberOfPage {
	var NumberOfPage int64

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("identities").Count(&NumberOfPage)

	count := contract.NumberOfPage{
		NumberOfPage: int64(math.Ceil(float64(NumberOfPage) / 5.00)),
	}
	return &count
}

func (s *employeeService) SPGetListByName(name string) (*[]contract.ListSubmission, error) {
	var ListSubmission []contract.ListSubmission
	var List []contract.ListAll

	namePersen := strings.Join([]string{"%", name, "%"}, "")

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Raw("SELECT identities.id, identities.updated_at, identities.deleted_at, identities.id_cust, identities.nik, identities.nama_lengkap, identities.tempat_lahir, identities.tanggal_lahir, identities.alamat, identities.pekerjaan, identities.pendapatan_perbulan, identities.bukti_ktp, identities.bukti_gaji, identities.status, submissions.created_at, submissions.alamat_rumah, submissions.luas_tanah, submissions.harga_rumah, submissions.jangka_pembayaran, submissions.dokumen_pendukung, submissions.status_kelengkapan FROM identities LEFT JOIN submissions ON identities.id_cust = submissions.id_cust WHERE identities.nama_lengkap LIKE ? LIMIT ? offset ?", namePersen, 30, 0).Find(&List).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i, v := range List {
		if v.StatusKelengkapan == "" {
			tanggal := v.UpdatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      "-",
			}
			ListSubmission = append(ListSubmission, lsubmission)
		} else {
			tanggal := v.CreatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.StatusKelengkapan,
				Rekomendasi:      Recommendation(&v.PendapatanPerbulan, &v.HargaRumah, &v.JangkaPembayaran),
			}
			ListSubmission = append(ListSubmission, lsubmission)
		}
	}
	return &ListSubmission, nil
}

func (s *employeeService) SPGetSubmission(id uint) (*contract.Submission, error) {

	var getSubmission contract.Submission

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("submissions").Last(&getSubmission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &getSubmission, nil
}

func (s *employeeService) SPPostSubmissionStatus(submissionStatus *contract.Submission, id uint) (*contract.StatusKelengkapanReturn, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	var submissionUpdates contract.Submission
	submissionUpdates.StatusKelengkapan = submissionStatus.StatusKelengkapan
	var submission contract.Submission
	err := db.DbConnection.Table("submissions").Last(&submission, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&submission).Updates(submissionUpdates).Error
	if err != nil {
		return nil, err
	}
	sreturn := contract.StatusKelengkapanReturn{
		Id:                id,
		IdCust:            submission.IdCust,
		StatusKelengkapan: submissionUpdates.StatusKelengkapan,
	}
	return &sreturn, nil
}

func (s *employeeService) SPPostIdentityStatus(statusPengajuan *contract.Identity, id uint) (*contract.StatusReturn, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	var pengajuanUpdates contract.Identity
	pengajuanUpdates.Status = statusPengajuan.Status
	var pengajuan contract.Identity

	err := db.DbConnection.Table("identities").Last(&pengajuan, "id_cust = ?", id).Error

	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&pengajuan).Updates(pengajuanUpdates).Error
	if err != nil {
		return nil, err
	}
	sreturn := contract.StatusReturn{
		Id:     pengajuan.IdCust,
		IdCust: pengajuan.IdCust,
		Status: pengajuanUpdates.Status,
	}
	return &sreturn, nil
}

func (s *employeeService) SPGetFileKtp(buktiKtp string) *minio.Object {
	fileName := strings.Join([]string{"ktp/", buktiKtp}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *employeeService) SPGetFileBuktiGaji(buktiGaji string) *minio.Object {
	fileName := strings.Join([]string{"slip-gaji/", buktiGaji}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *employeeService) SPGetFileBuktiPendukung(buktiPendukung string) *minio.Object {
	fileName := strings.Join([]string{"bukti-pendukung/", buktiPendukung}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *employeeService) SPGetIdentityEmployee(id uint) (*contract.IdentityReturn, error) {
	var getIdentity contract.Identity

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("identities").Last(&getIdentity, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}

	rgetIdentity := contract.IdentityReturn{
		Id:                 getIdentity.IdCust,
		IdCust:             getIdentity.IdCust,
		Nik:                getIdentity.Nik,
		NamaLengkap:        getIdentity.NamaLengkap,
		TempatLahir:        getIdentity.TempatLahir,
		TanggalLahir:       getIdentity.TanggalLahir,
		Alamat:             getIdentity.Alamat,
		Pekerjaan:          getIdentity.Pekerjaan,
		PendapatanPerbulan: getIdentity.PendapatanPerbulan,
		BuktiKtp:           getIdentity.BuktiKtp,
		BuktiGaji:          getIdentity.BuktiGaji,
		Status:             getIdentity.Status,
	}
	return &rgetIdentity, nil
}

func (s *employeeService) SPGetStatusTotal() (*contract.StatusTotalIdentity, error) {
	var countMV int64
	var countT int64
	var countTT int64
	var countMP int64
	var countD int64
	var countTD int64

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.Table("identities").Where("status = ?", "Menunggu Verifikasi").Count(&countMV)

	db.DbConnection.Table("identities").Where("status = ?", "Terverifikasi").Count(&countT)

	db.DbConnection.Table("identities").Where("status = ?", "Tidak Terverifikasi").Count(&countTT)

	db.DbConnection.Table("submissions").Where("status_kelengkapan = ?", "Menunggu Persetujuan").Count(&countMP)

	db.DbConnection.Table("submissions").Where("status_kelengkapan = ?", "Disetujui").Count(&countD)

	db.DbConnection.Table("submissions").Where("status_kelengkapan = ?", "Tidak Disetujui").Count(&countTD)

	rgetStatusTotal := contract.StatusTotalIdentity{
		MenungguVerifikasi:  uint(countMV),
		Terverifikasi:       uint(countT),
		TidakTerverifikasi:  uint(countTT),
		MenungguPersetujuan: uint(countMP),
		Disetujui:           uint(countD),
		TidakDisetujui:      uint(countTD),
	}
	return &rgetStatusTotal, nil
}

func (s *employeeService) SPDownloadReport() *excelize.File {
	f := excelize.NewFile()
	var ListAccepted []contract.ListAccepted
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Raw("SELECT * FROM identities cross JOIN submissions WHERE identities.id_cust = submissions.id_cust AND submissions.status_kelengkapan = ?", "Disetujui").Find(&ListAccepted).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	f.SetCellValue("Sheet1", "a1", "No")
	f.SetCellValue("Sheet1", "b1", "NIK")
	f.SetCellValue("Sheet1", "c1", "Nama Lengkap")
	f.SetCellValue("Sheet1", "d1", "Tempat Lahir")
	f.SetCellValue("Sheet1", "e1", "Tanggal Lahir")
	f.SetCellValue("Sheet1", "f1", "Alamat")
	f.SetCellValue("Sheet1", "g1", "Pekerjaan")
	f.SetCellValue("Sheet1", "h1", "Pendapatan Perbulan")
	f.SetCellValue("Sheet1", "i1", "Status Data Diri")
	f.SetCellValue("Sheet1", "j1", "Alamat Rumah KPR")
	f.SetCellValue("Sheet1", "k1", "Luas Tanah KPR")
	f.SetCellValue("Sheet1", "l1", "Harga Rumah KPR")
	f.SetCellValue("Sheet1", "m1", "Jangka Pembayaran")
	f.SetCellValue("Sheet1", "n1", "Status Kelengkapan")

	for i, v := range ListAccepted {
		incstr := strconv.Itoa(i + 2)
		f.SetCellValue("Sheet1", "a"+incstr, i+1)
		f.SetCellValue("Sheet1", "b"+incstr, v.Nik)
		f.SetCellValue("Sheet1", "c"+incstr, v.NamaLengkap)
		f.SetCellValue("Sheet1", "d"+incstr, v.TempatLahir)
		f.SetCellValue("Sheet1", "e"+incstr, v.TanggalLahir)
		f.SetCellValue("Sheet1", "f"+incstr, v.Alamat)
		f.SetCellValue("Sheet1", "g"+incstr, v.Pekerjaan)
		f.SetCellValue("Sheet1", "h"+incstr, v.PendapatanPerbulan)
		f.SetCellValue("Sheet1", "i"+incstr, v.Status)
		f.SetCellValue("Sheet1", "j"+incstr, v.AlamatRumah)
		f.SetCellValue("Sheet1", "k"+incstr, v.LuasTanah)
		f.SetCellValue("Sheet1", "l"+incstr, v.HargaRumah)
		f.SetCellValue("Sheet1", "m"+incstr, v.JangkaPembayaran)
		f.SetCellValue("Sheet1", "n"+incstr, v.StatusKelengkapan)
	}
	// if err = f.SaveAs("./Export.xlsx"); err != nil {
	// 	println(err.Error())
	// }
	return f
}

func (s *employeeService) SPGetListSubmissionParam(pages int, perPage int, name string) (*[]contract.ListSubmission, error) {
	var ListSubmission []contract.ListSubmission
	var List []contract.ListAll

	namePersen := strings.Join([]string{"%", name, "%"}, "")

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Raw("SELECT identities.id, identities.updated_at, identities.deleted_at, identities.id_cust, identities.nik, identities.nama_lengkap, identities.tempat_lahir, identities.tanggal_lahir, identities.alamat, identities.pekerjaan, identities.pendapatan_perbulan, identities.bukti_ktp, identities.bukti_gaji, identities.status, submissions.created_at, submissions.alamat_rumah, submissions.luas_tanah, submissions.harga_rumah, submissions.jangka_pembayaran, submissions.dokumen_pendukung, submissions.status_kelengkapan FROM identities LEFT JOIN submissions ON identities.id_cust = submissions.id_cust WHERE identities.nama_lengkap LIKE ? LIMIT ? offset ?", namePersen, perPage, (pages-1)*perPage).Find(&List).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i, v := range List {
		if v.StatusKelengkapan == "" {
			tanggal := v.UpdatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.Status,
				Rekomendasi:      "-",
			}
			ListSubmission = append(ListSubmission, lsubmission)
		} else {
			tanggal := v.CreatedAt.Format(time.RFC1123)
			lsubmission := contract.ListSubmission{
				Id:               uint(i + 1),
				TanggalPengajuan: tanggal[:25],
				NamaLengkap:      v.NamaLengkap,
				Status:           v.StatusKelengkapan,
				Rekomendasi:      Recommendation(&v.PendapatanPerbulan, &v.HargaRumah, &v.JangkaPembayaran),
			}
			ListSubmission = append(ListSubmission, lsubmission)
		}
	}
	return &ListSubmission, nil
}
