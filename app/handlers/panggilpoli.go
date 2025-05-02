package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/dhiafahmig/Go-DisplayPoli/app/services"
)

// PanggilPoliHandler menangani tampilan memanggil pasien di poli
type PanggilPoliHandler struct {
	DB *gorm.DB
}

// NewPanggilPoliHandler membuat instance baru dari PanggilPoliHandler
func NewPanggilPoliHandler(db *gorm.DB) *PanggilPoliHandler {
	return &PanggilPoliHandler{DB: db}
}

// HandlePanggil menampilkan halaman panggil poli
func (h *PanggilPoliHandler) HandlePanggil(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")
	kdDisplay := c.Param("kd_display")
	pasienList := h.getPasienList(kdRuangPoli)

	c.HTML(http.StatusOK, "panggilpoli.html", gin.H{
		"KdRuangPoli": kdRuangPoli,
		"KdDisplay":   kdDisplay,
		"PasienList":  pasienList,
	})
}

// HandleLog mengelola status log antrian pasien
func (h *PanggilPoliHandler) HandleLog(c *gin.Context) {
	var input struct {
		KdDokter    string `json:"kd_dokter"`
		NoRawat     string `json:"no_rawat" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli" binding:"required"`
		Type        string `json:"type" binding:"required"` // 'ada' atau 'tidak'
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := "1" // default: tidak ada
	if input.Type == "ada" {
		status = "0"
	}

	// Update or insert log
	result := h.DB.Exec(`
		INSERT INTO bw_log_antrian_poli (no_rawat, kd_ruang_poli, status)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE kd_ruang_poli = ?, status = ?
	`, input.NoRawat, input.KdRuangPoli, status, input.KdRuangPoli, status)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pasien berhasil diperbarui"})
}

// ResetLog menghapus log antrian pasien
func (h *PanggilPoliHandler) ResetLog(c *gin.Context) {
	noRawat := c.Param("no_rawat")

	result := h.DB.Table("bw_log_antrian_poli").Where("no_rawat = ?", noRawat).Delete(nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset log berhasil"})
}

// PanggilPasien mengirim event untuk memanggil pasien
func (h *PanggilPoliHandler) PanggilPasien(c *gin.Context) {
	var input struct {
		NmPasien    string `json:"nm_pasien" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli" binding:"required"`
		NmPoli      string `json:"nm_poli" binding:"required"`
		NoReg       string `json:"no_reg" binding:"required"`
		KdDisplay   string `json:"kd_display" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Di sini kita menggunakan channel broadcaster dari main.go
	// untuk mengirim pesan ke semua klien websocket

	c.JSON(http.StatusOK, gin.H{"message": "Panggilan berhasil dikirim"})
}

// getPasienList mendapatkan daftar pasien untuk poli tertentu
func (h *PanggilPoliHandler) getPasienList(kdRuangPoli string) []map[string]interface{} {
	hari := services.GetDayList()[time.Now().Format("Monday")]
	var results []map[string]interface{}

	h.DB.Table("reg_periksa").
		Select("reg_periksa.no_reg, reg_periksa.no_rawat, reg_periksa.no_rkm_medis, reg_periksa.kd_dokter, reg_periksa.kd_pj, jadwal.hari_kerja, jadwal.jam_mulai, bw_ruangpoli_dokter.kd_ruang_poli, bw_ruangpoli_dokter.nama_dokter, pasien.nm_pasien, bw_log_antrian_poli.status, penjab.png_jawab, poliklinik.nm_poli").
		Joins("LEFT JOIN bw_log_antrian_poli ON bw_log_antrian_poli.no_rawat = reg_periksa.no_rawat").
		Joins("JOIN bw_ruangpoli_dokter ON reg_periksa.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
		Joins("JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
		Joins("JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis").
		Joins("JOIN penjab ON reg_periksa.kd_pj = penjab.kd_pj").
		Joins("JOIN poliklinik ON reg_periksa.kd_poli = poliklinik.kd_poli").
		Where("reg_periksa.tgl_registrasi = ?", time.Now().Format("2006-01-02")).
		Where("jadwal.hari_kerja = ?", hari).
		Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", kdRuangPoli).
		Order("jadwal.jam_mulai asc").
		Order("reg_periksa.no_reg asc").
		Order("reg_periksa.jam_reg asc").
		Find(&results)

	return results
}
