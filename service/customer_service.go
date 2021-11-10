package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

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

type customerService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type CustomerServiceInterface interface {
	SCGetCheckApply(idCust uint) string
	VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error)
	SCCreateIdentity(identity *contract.Identity, idCust uint) (*contract.IdentityReturn, error)
	SCCreateSubmission(submission *contract.Submission, idCust uint) *contract.SubmissionReturn
	SCGetSubmissionStatus(id uint) string
	SCGetSubmission(id uint) (*contract.SubmissionReturn, error)
	SCUploadFileKTP(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadBuktiKtpReturn, error)
	SCUploadFileGaji(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadBuktiGajiReturn, error)
	SCUploadFilePendukung(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadDokumenPendukungReturn, error)
	SCGetFileKtpCustomer(buktiKtp string) *minio.Object
	SCGetFileBuktiGajiCustomer(buktiGaji string) *minio.Object
	SCGetFilePendukungCustomer(buktiFilependukung string) *minio.Object
	SCGetIdentity(id uint) (*contract.IdentityReturn, error)
	SCUpdateIdentityCustomer(identity *contract.Identity, id uint) (*contract.IdentityReturn, error)
}

func NewCustomerService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *customerService {
	return &customerService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}

func (s *customerService) SCGetCheckApply(idCust uint) string {
	var identity contract.Identity

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("identities").First(&identity, "id_cust = ?", idCust).Error

	if err != nil {
		return "Anda sedang tidak mengajukan KPR saat ini"
	}

	return "Anda sedang mengajukan KPR saat ini"
}

func (s *customerService) VerifyToken(req *contract.ValidateTokenRequestContract) (*contract.JWTMapClaim, error) {
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

func (s *customerService) SCCreateIdentity(identity *contract.Identity, idCust uint) (*contract.IdentityReturn, error) {
	identity.IdCust = idCust
	identity.Status = "Menunggu Verifikasi"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Create(&identity).Error
	if err != nil {
		log.Error("error connect db, %q", err)
	}

	pReturn := contract.IdentityReturn{
		IdCust:             identity.IdCust,
		Nik:                identity.Nik,
		NamaLengkap:        identity.NamaLengkap,
		TempatLahir:        identity.TempatLahir,
		TanggalLahir:       identity.TanggalLahir,
		Alamat:             identity.Alamat,
		Pekerjaan:          identity.Pekerjaan,
		PendapatanPerbulan: identity.PendapatanPerbulan,
		BuktiKtp:           identity.BuktiKtp,
		BuktiGaji:          identity.BuktiGaji,
		Status:             identity.Status,
	}
	return &pReturn, nil
}

func (s *customerService) SCCreateSubmission(submission *contract.Submission, id uint) *contract.SubmissionReturn {

	submission.IdCust = id
	submission.IdPengajuan = id
	submission.StatusKelengkapan = "Menunggu Persetujuan"

	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	db.DbConnection.Create(&submission)

	kReturn := contract.SubmissionReturn{
		IdKelengkapan:     submission.ID,
		IdCust:            submission.IdCust,
		IdPengajuan:       submission.IdPengajuan,
		AlamatRumah:       submission.AlamatRumah,
		LuasTanah:         submission.LuasTanah,
		HargaRumah:        submission.HargaRumah,
		JangkaPembayaran:  submission.JangkaPembayaran,
		DokumenPendukung:  submission.DokumenPendukung,
		StatusKelengkapan: submission.StatusKelengkapan,
	}
	return &kReturn

}

func (s *customerService) SCGetSubmission(id uint) (*contract.SubmissionReturn, error) {
	var getSubmission contract.Submission
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Table("submissions").Last(&getSubmission, "id_pengajuan = ?", id).Error
	if err != nil {
		return nil, err
	}
	kReturn := contract.SubmissionReturn{
		IdCust:            getSubmission.IdCust,
		IdPengajuan:       getSubmission.IdPengajuan,
		AlamatRumah:       getSubmission.AlamatRumah,
		LuasTanah:         getSubmission.LuasTanah,
		HargaRumah:        getSubmission.HargaRumah,
		JangkaPembayaran:  getSubmission.JangkaPembayaran,
		DokumenPendukung:  getSubmission.DokumenPendukung,
		StatusKelengkapan: getSubmission.StatusKelengkapan,
	}
	return &kReturn, nil
}

func (s *customerService) SCGetSubmissionStatus(id uint) string {
	var getStatusKelengkapan contract.Submission
	db := mysql.NewMysqlClient(*mysql.MysqlInit())
	err := db.DbConnection.Table("submissions").Last(&getStatusKelengkapan, "id_pengajuan = ?", id).Error
	if err != nil {
		return "Menu Submission invisible(Menu disable)"
	}
	return "Menu Submission visible(Menu able)"
}

func (s *customerService) SCUploadFileKTP(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadBuktiKtpReturn, error) {
	uploadTime := time.Now().Local().String()
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"ktp-", idString, "-", resp.Username, uploadTime[:10], "-", uploadTime[11:22], ".pdf"}, "")
	fileName := strings.Join([]string{"ktp/", fileLink}, "")
	link := strings.Join([]string{"http://backend-c-home-loans.digitalent.rakamin.com/cgetfilektp/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return nil, err
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	uploadLink := contract.UploadBuktiKtpReturn{
		BuktiKtp: link,
	}
	return &uploadLink, nil
}

func (s *customerService) SCUploadFileGaji(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadBuktiGajiReturn, error) {
	uploadTime := time.Now().Local().String()
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"gaji-", idString, "-", resp.Username, uploadTime[:10], "-", uploadTime[11:22], ".pdf"}, "")
	fileName := strings.Join([]string{"slip-gaji/", fileLink}, "")
	link := strings.Join([]string{"http://backend-c-home-loans.digitalent.rakamin.com/cgetfilegaji/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return nil, err
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	uploadLink := contract.UploadBuktiGajiReturn{
		BuktiGaji: link,
	}
	return &uploadLink, nil
}

func (s *customerService) SCUploadFilePendukung(file *multipart.File, handler *multipart.FileHeader, resp *contract.JWTMapClaim) (*contract.UploadDokumenPendukungReturn, error) {
	uploadTime := time.Now().Local().String()
	idString := strconv.Itoa(int(resp.IdUser))
	fileLink := strings.Join([]string{"pendukung-", idString, "-", resp.Username, uploadTime[:10], "-", uploadTime[11:22], ".pdf"}, "")
	fileName := strings.Join([]string{"bukti-pendukung/", fileLink}, "")
	link := strings.Join([]string{"http://backend-c-home-loans.digitalent.rakamin.com/cgetfilependukung/", fileLink}, "")

	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()

	fileReader := io.Reader(*file)
	uploadInfo, err := mi.MinioClient.PutObject(ctx, mi.BucketName, fileName, fileReader, handler.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("Error in uploading the file #%s: %v.", fileName, err)
		return nil, err
	}

	log.Printf("Uploading the file #%s succeeded!", fileName)
	fmt.Println("UploadInfo:")
	fmt.Printf("%+v\n", uploadInfo)

	uploadLink := contract.UploadDokumenPendukungReturn{
		DokumenPendukung: link,
	}
	return &uploadLink, nil
}

func (s *customerService) SCGetFileKtpCustomer(buktiKtp string) *minio.Object {
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

func (s *customerService) SCGetFileBuktiGajiCustomer(buktiGaji string) *minio.Object {
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

func (s *customerService) SCGetFilePendukungCustomer(buktiFilependukung string) *minio.Object {
	fileName := strings.Join([]string{"bukti-pendukung/", buktiFilependukung}, "")
	mi := miniopkg.NewMinioClient(*miniopkg.MinioInit())

	ctx := context.Background()
	obj, err := mi.MinioClient.GetObject(ctx, mi.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error in getting the object: %v.", err)
		return nil
	}
	return obj
}

func (s *customerService) SCGetIdentity(id uint) (*contract.IdentityReturn, error) {

	var getIdentity contract.Identity

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	err := db.DbConnection.Table("identities").Last(&getIdentity, "id_cust = ?", id).Error
	if err != nil {
		return nil, err
	}

	rgetIdentity := contract.IdentityReturn{
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

func (s *customerService) SCUpdateIdentityCustomer(identity *contract.Identity, id uint) (*contract.IdentityReturn, error) {
	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	var identityUpdates contract.Identity
	identityUpdates.Nik = identity.Nik
	identityUpdates.NamaLengkap = identity.NamaLengkap
	identityUpdates.TempatLahir = identity.TempatLahir
	identityUpdates.TanggalLahir = identity.TanggalLahir
	identityUpdates.Alamat = identity.Alamat
	identityUpdates.Pekerjaan = identity.Pekerjaan
	identityUpdates.PendapatanPerbulan = identity.PendapatanPerbulan
	identityUpdates.BuktiKtp = identity.BuktiKtp
	identityUpdates.BuktiGaji = identity.BuktiGaji
	identityUpdates.Status = "Menunggu Verifikasi"
	var identityUp contract.Identity

	err := db.DbConnection.Table("identities").Where("id_cust = ?", id).Find(&identityUp).Error

	if err != nil {
		return nil, err
	}
	err = db.DbConnection.Model(&identityUp).Updates(identityUpdates).Error
	if err != nil {
		return nil, err
	}

	ureturn := contract.IdentityReturn{
		Id:                 identityUp.ID,
		IdCust:             identityUp.IdCust,
		Nik:                identityUpdates.Nik,
		NamaLengkap:        identityUpdates.NamaLengkap,
		TempatLahir:        identityUpdates.TempatLahir,
		TanggalLahir:       identityUpdates.TanggalLahir,
		Alamat:             identityUpdates.Alamat,
		Pekerjaan:          identityUpdates.Pekerjaan,
		PendapatanPerbulan: identityUpdates.PendapatanPerbulan,
		BuktiKtp:           identityUpdates.BuktiKtp,
		BuktiGaji:          identityUpdates.BuktiGaji,
		Status:             identityUp.Status,
	}
	return &ureturn, nil
}
