package contract

import (
	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	IdCust            uint    `gorm:"not null" json:"id_cust"`
	IdPengajuan       uint    `gorm:"not null" json:"id_pengajuan"`
	AlamatRumah       string  `gorm:"not null" json:"alamat_rumah" validate:"required"`
	LuasTanah         float64 `gorm:"not null" json:"luas_tanah" validate:"required"`
	HargaRumah        float64 `gorm:"not null" json:"harga_rumah" validate:"required"`
	JangkaPembayaran  uint    `gorm:"not null" json:"jangka_pembayaran" validate:"required"`
	DokumenPendukung  string  `gorm:"not null" json:"dokumen_pendukung" validate:"required"`
	StatusKelengkapan string  `gorm:"not null" json:"status_kelengkapan"`
}

type SubmissionReturn struct {
	IdKelengkapan     uint    `gorm:"not null" json:"id_kelengkapan"`
	IdCust            uint    `gorm:"not null" json:"id_cust"`
	IdPengajuan       uint    `gorm:"not null" json:"id_pengajuan"`
	AlamatRumah       string  `gorm:"not null" json:"alamat_rumah" validate:"required"`
	LuasTanah         float64 `gorm:"not null" json:"luas_tanah" validate:"required"`
	HargaRumah        float64 `gorm:"not null" json:"harga_rumah" validate:"required"`
	JangkaPembayaran  uint    `gorm:"not null" json:"jangka_pembayaran" validate:"required"`
	DokumenPendukung  string  `gorm:"not null" json:"dokumen_pendukung" validate:"required"`
	StatusKelengkapan string  `gorm:"not null" json:"status_kelengkapan"`
}

type ListAccepted struct {
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir"`
	Alamat             string  `gorm:"not null" json:"alamat"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp"`
	BuktiGaji          string  `gorm:"not null" json:"bukti_gaji"`
	Status             string  `gorm:"not null" json:"status"`
	AlamatRumah        string  `gorm:"not null" json:"alamat_rumah"`
	LuasTanah          float64 `gorm:"not null" json:"luas_tanah"`
	HargaRumah         float64 `gorm:"not null" json:"harga_rumah"`
	JangkaPembayaran   uint    `gorm:"not null" json:"jangka_pembayaran"`
	DokumenPendukung   string  `gorm:"not null" json:"dokumen_pendukung"`
	StatusKelengkapan  string  `gorm:"not null" json:"status_kelengkapan"`
}

type ListAll struct {
	gorm.Model
	IdCust             uint    `gorm:"not null;unique" json:"id_cust"`
	Nik                string  `gorm:"not null;unique" json:"nik"`
	NamaLengkap        string  `gorm:"not null" json:"nama_lengkap"`
	TempatLahir        string  `gorm:"not null" json:"tempat_lahir"`
	TanggalLahir       string  `gorm:"not null" json:"tanggal_lahir"`
	Alamat             string  `gorm:"not null" json:"alamat"`
	Pekerjaan          string  `gorm:"not null" json:"pekerjaan"`
	PendapatanPerbulan float64 `gorm:"not null" json:"pendapatan_perbulan"`
	BuktiKtp           string  `gorm:"not null" json:"bukti_ktp"`
	BuktiGaji          string  `gorm:"not null" json:"bukti_gaji"`
	Status             string  `gorm:"not null" json:"status"`
	AlamatRumah        string  `gorm:"not null" json:"alamat_rumah"`
	LuasTanah          float64 `gorm:"not null" json:"luas_tanah"`
	HargaRumah         float64 `gorm:"not null" json:"harga_rumah"`
	JangkaPembayaran   uint    `gorm:"not null" json:"jangka_pembayaran"`
	DokumenPendukung   string  `gorm:"not null" json:"dokumen_pendukung"`
	StatusKelengkapan  string  `gorm:"not null" json:"status_kelengkapan"`
}

type StatusKelengkapanReturn struct {
	Id                uint   `gorm:"not null;unique" json:"id"`
	IdCust            uint   `gorm:"not null;unique" json:"id_cust"`
	StatusKelengkapan string `gorm:"not null" json:"status_kelengkapan"`
}

type UploadDokumenPendukungReturn struct {
	DokumenPendukung string `gorm:"not null" json:"dokumen_pendukung"`
}
