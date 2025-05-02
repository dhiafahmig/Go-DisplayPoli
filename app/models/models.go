package models

import (
	"time"
)

// Display mewakili model untuk tabel bw_display_poli
type Display struct {
	KdDisplay   string `json:"kd_display" gorm:"column:kd_display;primaryKey"`
	NamaDisplay string `json:"nama_display" gorm:"column:nama_display"`
}

// TableName menentukan nama tabel untuk model Display
func (Display) TableName() string {
	return "bw_display_poli"
}

// RuangPoli mewakili model untuk tabel bw_ruang_poli
type RuangPoli struct {
	KdRuangPoli       string   `json:"kd_ruang_poli" gorm:"column:kd_ruang_poli;primaryKey"`
	NamaRuangPoli     string   `json:"nama_ruang_poli" gorm:"column:nama_ruang_poli"`
	KdDisplay         string   `json:"kd_display" gorm:"column:kd_display"`
	PosisiDisplayPoli int      `json:"posisi_display_poli" gorm:"column:posisi_display_poli"`
	Display           Display  `json:"display" gorm:"foreignKey:KdDisplay;references:KdDisplay"`
	Pasien            []Pasien `json:"pasien,omitempty" gorm:"-"`
}

// TableName menentukan nama tabel untuk model RuangPoli
func (RuangPoli) TableName() string {
	return "bw_ruang_poli"
}

// Dokter mewakili model untuk tabel dokter
type Dokter struct {
	KdDokter string `json:"kd_dokter" gorm:"column:kd_dokter;primaryKey"`
	NmDokter string `json:"nm_dokter" gorm:"column:nm_dokter"`
	Status   string `json:"status" gorm:"column:status"`
}

// TableName menentukan nama tabel untuk model Dokter
func (Dokter) TableName() string {
	return "dokter"
}

// RuangPoliDokter mewakili model untuk tabel bw_ruangpoli_dokter
type RuangPoliDokter struct {
	KdDokter    string    `json:"kd_dokter" gorm:"column:kd_dokter;primaryKey"`
	NamaDokter  string    `json:"nama_dokter" gorm:"column:nama_dokter"`
	KdRuangPoli string    `json:"kd_ruang_poli" gorm:"column:kd_ruang_poli"`
	RuangPoli   RuangPoli `json:"ruang_poli" gorm:"foreignKey:KdRuangPoli;references:KdRuangPoli"`
}

// TableName menentukan nama tabel untuk model RuangPoliDokter
func (RuangPoliDokter) TableName() string {
	return "bw_ruangpoli_dokter"
}

// JadwalDokter mewakili model untuk tabel bw_jadwal_dokter
type JadwalDokter struct {
	KdDokter   string     `json:"kd_dokter" gorm:"column:kd_dokter"`
	HariKerja  string     `json:"hari_kerja" gorm:"column:hari_kerja"`
	JamMulai   string     `json:"jam_mulai" gorm:"column:jam_mulai"`
	JamSelesai string     `json:"jam_selesai" gorm:"column:jam_selesai"`
	KdPoli     string     `json:"kd_poli" gorm:"column:kd_poli"`
	Dokter     Dokter     `json:"dokter" gorm:"foreignKey:KdDokter;references:KdDokter"`
	Poliklinik Poliklinik `json:"poliklinik" gorm:"foreignKey:KdPoli;references:KdPoli"`
}

// TableName menentukan nama tabel untuk model JadwalDokter
func (JadwalDokter) TableName() string {
	return "bw_jadwal_dokter"
}

// Poliklinik mewakili model untuk tabel poliklinik
type Poliklinik struct {
	KdPoli string `json:"kd_poli" gorm:"column:kd_poli;primaryKey"`
	NmPoli string `json:"nm_poli" gorm:"column:nm_poli"`
	Status string `json:"status" gorm:"column:status"`
}

// TableName menentukan nama tabel untuk model Poliklinik
func (Poliklinik) TableName() string {
	return "poliklinik"
}

// Pasien mewakili model untuk tabel pasien
type Pasien struct {
	NoRkmMedis string `json:"no_rkm_medis" gorm:"column:no_rkm_medis;primaryKey"`
	NmPasien   string `json:"nm_pasien" gorm:"column:nm_pasien"`
}

// TableName menentukan nama tabel untuk model Pasien
func (Pasien) TableName() string {
	return "pasien"
}

// RegPeriksa mewakili model untuk tabel reg_periksa
type RegPeriksa struct {
	NoReg         string     `json:"no_reg" gorm:"column:no_reg"`
	NoRawat       string     `json:"no_rawat" gorm:"column:no_rawat;primaryKey"`
	TglRegistrasi time.Time  `json:"tgl_registrasi" gorm:"column:tgl_registrasi"`
	JamReg        string     `json:"jam_reg" gorm:"column:jam_reg"`
	KdDokter      string     `json:"kd_dokter" gorm:"column:kd_dokter"`
	NoRkmMedis    string     `json:"no_rkm_medis" gorm:"column:no_rkm_medis"`
	KdPoli        string     `json:"kd_poli" gorm:"column:kd_poli"`
	KdPj          string     `json:"kd_pj" gorm:"column:kd_pj"`
	Pasien        Pasien     `json:"pasien" gorm:"foreignKey:NoRkmMedis;references:NoRkmMedis"`
	Dokter        Dokter     `json:"dokter" gorm:"foreignKey:KdDokter;references:KdDokter"`
	Poliklinik    Poliklinik `json:"poliklinik" gorm:"foreignKey:KdPoli;references:KdPoli"`
}

// TableName menentukan nama tabel untuk model RegPeriksa
func (RegPeriksa) TableName() string {
	return "reg_periksa"
}

// LogAntrianPoli mewakili model untuk tabel bw_log_antrian_poli
type LogAntrianPoli struct {
	NoRawat     string `json:"no_rawat" gorm:"column:no_rawat;primaryKey"`
	KdRuangPoli string `json:"kd_ruang_poli" gorm:"column:kd_ruang_poli"`
	Status      string `json:"status" gorm:"column:status"`
}

// TableName menentukan nama tabel untuk model LogAntrianPoli
func (LogAntrianPoli) TableName() string {
	return "bw_log_antrian_poli"
}

// Penjab mewakili model untuk tabel penjab (penjamin)
type Penjab struct {
	KdPj     string `json:"kd_pj" gorm:"column:kd_pj;primaryKey"`
	PngJawab string `json:"png_jawab" gorm:"column:png_jawab"`
}

// TableName menentukan nama tabel untuk model Penjab
func (Penjab) TableName() string {
	return "penjab"
}
