package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/dhiafahmig/Go-DisplayPoli/app/services"
)

// JadwalDokterHandler menangani tampilan jadwal dokter
type JadwalDokterHandler struct {
	DB *gorm.DB
}

// NewJadwalDokterHandler membuat instance baru dari JadwalDokterHandler
func NewJadwalDokterHandler(db *gorm.DB) *JadwalDokterHandler {
	return &JadwalDokterHandler{DB: db}
}

// HandleJadwal menampilkan halaman jadwal dokter
func (h *JadwalDokterHandler) HandleJadwal(c *gin.Context) {
	hariKerja := c.DefaultQuery("hari", services.GetDayList()[time.Now().Format("Monday")])
	dokterList := h.getDokterByHari(hariKerja)
	poliList := h.getPoliList()

	c.HTML(http.StatusOK, "jadwaldokter.html", gin.H{
		"DokterList": dokterList,
		"PoliList":   poliList,
		"HariKerja":  hariKerja,
		"HariList":   services.GetDayList(),
	})
}

// UbahJadwalDokter mengedit jadwal dokter yang ada
func (h *JadwalDokterHandler) UbahJadwalDokter(c *gin.Context) {
	var input struct {
		KdDokter       string `json:"kd_dokter" binding:"required"`
		HariKerja      string `json:"hari_kerja" binding:"required"`
		JamMulai       string `json:"jam_mulai" binding:"required"`
		JamSelesai     string `json:"jam_selesai" binding:"required"`
		JamMulaiBaru   string `json:"jam_mulai_baru" binding:"required"`
		JamSelesaiBaru string `json:"jam_selesai_baru" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Table("bw_jadwal_dokter").
		Where("kd_dokter = ?", input.KdDokter).
		Where("hari_kerja = ?", input.HariKerja).
		Where("jam_mulai = ?", input.JamMulai).
		Where("jam_selesai = ?", input.JamSelesai).
		Updates(map[string]interface{}{
			"jam_mulai":   input.JamMulaiBaru,
			"jam_selesai": input.JamSelesaiBaru,
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil diubah"})
}

// HapusJadwalDokter menghapus jadwal dokter
func (h *JadwalDokterHandler) HapusJadwalDokter(c *gin.Context) {
	var input struct {
		KdDokter   string `json:"kd_dokter" binding:"required"`
		HariKerja  string `json:"hari_kerja" binding:"required"`
		JamMulai   string `json:"jam_mulai" binding:"required"`
		JamSelesai string `json:"jam_selesai" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Table("bw_jadwal_dokter").
		Where("kd_dokter = ?", input.KdDokter).
		Where("hari_kerja = ?", input.HariKerja).
		Where("jam_mulai = ?", input.JamMulai).
		Where("jam_selesai = ?", input.JamSelesai).
		Delete(nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil dihapus"})
}

// CariDokter mencari dokter berdasarkan kode atau nama
func (h *JadwalDokterHandler) CariDokter(c *gin.Context) {
	cariKode := c.Query("cari")

	var dokter []map[string]interface{}
	h.DB.Table("dokter").
		Select("dokter.kd_dokter, dokter.nm_dokter, dokter.status").
		Where("dokter.status = ?", "1").
		Where("dokter.kd_dokter LIKE ? OR dokter.nm_dokter LIKE ?", "%"+cariKode+"%", "%"+cariKode+"%").
		Limit(1).
		Find(&dokter)

	c.JSON(http.StatusOK, dokter)
}

// TambahJadwalDokter menambahkan jadwal dokter baru
func (h *JadwalDokterHandler) TambahJadwalDokter(c *gin.Context) {
	var input struct {
		KdDokter   string `json:"kd_dokter" binding:"required"`
		HariKerja  string `json:"hari_kerja" binding:"required"`
		JamMulai   string `json:"jam_mulai" binding:"required"`
		JamSelesai string `json:"jam_selesai" binding:"required"`
		KdPoli     string `json:"kd_poli" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Table("bw_jadwal_dokter").Create(map[string]interface{}{
		"kd_dokter":   input.KdDokter,
		"hari_kerja":  input.HariKerja,
		"jam_mulai":   input.JamMulai + ":00",
		"jam_selesai": input.JamSelesai + ":00",
		"kd_poli":     input.KdPoli,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil ditambahkan"})
}

// getDokterByHari mendapatkan daftar dokter berdasarkan hari
func (h *JadwalDokterHandler) getDokterByHari(hari string) []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("bw_jadwal_dokter").
		Select("dokter.nm_dokter, bw_jadwal_dokter.kd_dokter, bw_jadwal_dokter.hari_kerja, bw_jadwal_dokter.jam_mulai, bw_jadwal_dokter.jam_selesai, poliklinik.nm_poli").
		Joins("JOIN dokter ON bw_jadwal_dokter.kd_dokter = dokter.kd_dokter").
		Joins("JOIN poliklinik ON bw_jadwal_dokter.kd_poli = poliklinik.kd_poli").
		Where("bw_jadwal_dokter.hari_kerja = ?", hari).
		Find(&results)

	return results
}

// getPoliList mendapatkan daftar poli
func (h *JadwalDokterHandler) getPoliList() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("poliklinik").
		Select("poliklinik.kd_poli, poliklinik.nm_poli").
		Where("poliklinik.status = ?", "1").
		Find(&results)

	return results
}
