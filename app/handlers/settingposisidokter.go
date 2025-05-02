package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingPosisiDokterHandler menangani pengaturan posisi dokter
type SettingPosisiDokterHandler struct {
	DB *gorm.DB
}

// NewSettingPosisiDokterHandler membuat instance baru dari SettingPosisiDokterHandler
func NewSettingPosisiDokterHandler(db *gorm.DB) *SettingPosisiDokterHandler {
	return &SettingPosisiDokterHandler{DB: db}
}

// HandleSettings menampilkan halaman pengaturan posisi dokter
func (h *SettingPosisiDokterHandler) HandleSettings(c *gin.Context) {
	poliList := h.getPoli()
	dokterList := h.getListDokter()

	c.HTML(http.StatusOK, "settingposisidokter.html", gin.H{
		"PoliList":   poliList,
		"DokterList": dokterList,
	})
}

// EditPoliDokter memperbarui posisi dokter
func (h *SettingPosisiDokterHandler) EditPoliDokter(c *gin.Context) {
	var input struct {
		KdDokter    string `json:"kd_dokter" binding:"required"`
		NmDokter    string `json:"nm_dokter" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Menghapus entri yang ada jika kd_ruang_poli kosong
	if input.KdRuangPoli == "" {
		result := h.DB.Table("bw_ruangpoli_dokter").
			Where("kd_dokter = ?", input.KdDokter).
			Delete(nil)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Terjadi kesalahan saat menghapus posisi dokter",
				"color":   "danger",
				"icon":    "ban",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Posisi dokter berhasil dihapus",
			"color":   "success",
			"icon":    "check",
		})
		return
	}

	// Update atau Insert
	result := h.DB.Exec(`
		INSERT INTO bw_ruangpoli_dokter (kd_dokter, nama_dokter, kd_ruang_poli)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE nama_dokter = ?, kd_ruang_poli = ?
	`, input.KdDokter, input.NmDokter, input.KdRuangPoli, input.NmDokter, input.KdRuangPoli)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat update posisi dokter",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	// Ambil nama ruang poli
	var namaRuangPoli string
	h.DB.Table("bw_ruang_poli").
		Select("nama_ruang_poli").
		Where("kd_ruang_poli = ?", input.KdRuangPoli).
		Row().
		Scan(&namaRuangPoli)

	c.JSON(http.StatusOK, gin.H{
		"message": "Posisi Dokter Dipindahkan Ke " + input.KdRuangPoli + " - " + namaRuangPoli,
		"color":   "success",
		"icon":    "check",
	})
}

// getPoli mendapatkan daftar poli
func (h *SettingPosisiDokterHandler) getPoli() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("bw_ruang_poli").
		Select("bw_ruang_poli.kd_ruang_poli, bw_ruang_poli.nama_ruang_poli, bw_ruang_poli.kd_display, bw_ruang_poli.posisi_display_poli").
		Find(&results)

	return results
}

// getListDokter mendapatkan daftar dokter dengan posisi
func (h *SettingPosisiDokterHandler) getListDokter() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("dokter").
		Select("bw_ruangpoli_dokter.kd_ruang_poli, bw_ruang_poli.nama_ruang_poli, dokter.kd_dokter, dokter.nm_dokter").
		Joins("LEFT JOIN bw_ruangpoli_dokter ON dokter.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
		Joins("LEFT JOIN bw_ruang_poli ON bw_ruangpoli_dokter.kd_ruang_poli = bw_ruang_poli.kd_ruang_poli").
		Where("dokter.status = ?", "1").
		Find(&results)

	return results
}
