package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingPoliHandler menangani pengaturan poli
type SettingPoliHandler struct {
	DB *gorm.DB
}

// NewSettingPoliHandler membuat instance baru dari SettingPoliHandler
func NewSettingPoliHandler(db *gorm.DB) *SettingPoliHandler {
	return &SettingPoliHandler{DB: db}
}

// HandleSettings menampilkan halaman pengaturan poli
func (h *SettingPoliHandler) HandleSettings(c *gin.Context) {
	displays := h.getDisplays()
	poliList := h.getAllPoli()

	c.HTML(http.StatusOK, "settingpoli.html", gin.H{
		"Displays": displays,
		"PoliList": poliList,
	})
}

// AddPoli menambahkan poli baru
func (h *SettingPoliHandler) AddPoli(c *gin.Context) {
	var input struct {
		KdRuangPoli       string `form:"kd_ruang_poli" binding:"required"`
		NamaRuangPoli     string `form:"nama_ruang_poli" binding:"required"`
		KdDisplay         string `form:"kd_display" binding:"required"`
		PosisiDisplayPoli string `form:"posisi_display_poli" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Semua field harus diisi",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	result := h.DB.Table("bw_ruang_poli").Create(map[string]interface{}{
		"kd_ruang_poli":       input.KdRuangPoli,
		"nama_ruang_poli":     input.NamaRuangPoli,
		"kd_display":          input.KdDisplay,
		"posisi_display_poli": input.PosisiDisplayPoli,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat menambahkan poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poli berhasil ditambahkan!",
		"color":   "success",
		"icon":    "check",
	})
}

// EditPoli mengedit poli yang ada
func (h *SettingPoliHandler) EditPoli(c *gin.Context) {
	var input struct {
		KdRuangPoli       string `form:"kd_ruang_poli" binding:"required"`
		NamaRuangPoli     string `form:"nama_ruang_poli" binding:"required"`
		KdDisplay         string `form:"kd_display" binding:"required"`
		PosisiDisplayPoli string `form:"posisi_display_poli" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Semua field harus diisi",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	result := h.DB.Table("bw_ruang_poli").Where("kd_ruang_poli = ?", input.KdRuangPoli).Updates(map[string]interface{}{
		"nama_ruang_poli":     input.NamaRuangPoli,
		"kd_display":          input.KdDisplay,
		"posisi_display_poli": input.PosisiDisplayPoli,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat update poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poli berhasil diupdate!",
		"color":   "success",
		"icon":    "check",
	})
}

// DeletePoli menghapus poli yang ada
func (h *SettingPoliHandler) DeletePoli(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")

	result := h.DB.Table("bw_ruang_poli").Where("kd_ruang_poli = ?", kdRuangPoli).Delete(nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat menghapus poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poli berhasil dihapus!",
		"color":   "warning",
		"icon":    "check",
	})
}

// getDisplays mendapatkan daftar display poli
func (h *SettingPoliHandler) getDisplays() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("bw_display_poli").
		Select("bw_display_poli.kd_display, bw_display_poli.nama_display").
		Find(&results)

	return results
}

// getAllPoli mendapatkan daftar semua poli
func (h *SettingPoliHandler) getAllPoli() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("bw_ruang_poli").
		Select("bw_ruang_poli.kd_ruang_poli, bw_ruang_poli.nama_ruang_poli, bw_ruang_poli.kd_display, bw_ruang_poli.posisi_display_poli, bw_display_poli.nama_display").
		Joins("JOIN bw_display_poli ON bw_ruang_poli.kd_display = bw_display_poli.kd_display").
		Find(&results)

	return results
}
