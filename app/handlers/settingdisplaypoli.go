package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingDisplayPoliHandler menangani pengaturan display poli
type SettingDisplayPoliHandler struct {
	DB *gorm.DB
}

// NewSettingDisplayPoliHandler membuat instance baru dari SettingDisplayPoliHandler
func NewSettingDisplayPoliHandler(db *gorm.DB) *SettingDisplayPoliHandler {
	return &SettingDisplayPoliHandler{DB: db}
}

// HandleSettings menampilkan halaman pengaturan display poli
func (h *SettingDisplayPoliHandler) HandleSettings(c *gin.Context) {
	displays := h.getAllDisplay()

	c.HTML(http.StatusOK, "settingdisplaypoli.html", gin.H{
		"Displays": displays,
	})
}

// GetAllDisplay mengembalikan daftar semua display dalam format JSON
func (h *SettingDisplayPoliHandler) GetAllDisplay(c *gin.Context) {
	displays := h.getAllDisplay()
	c.JSON(http.StatusOK, displays)
}

// AddDisplay menambahkan display poli baru
func (h *SettingDisplayPoliHandler) AddDisplay(c *gin.Context) {
	var input struct {
		KdDisplay   string `form:"kd_display" binding:"required"`
		NamaDisplay string `form:"nama_display" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Kode Display dan Nama Display harus diisi",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	result := h.DB.Table("bw_display_poli").Create(map[string]interface{}{
		"kd_display":   input.KdDisplay,
		"nama_display": input.NamaDisplay,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat menambahkan display poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Display berhasil ditambahkan!",
		"color":   "success",
		"icon":    "check",
	})
}

// EditDisplay mengedit display poli yang ada
func (h *SettingDisplayPoliHandler) EditDisplay(c *gin.Context) {
	var input struct {
		KdDisplay   string `form:"kd_display" binding:"required"`
		NamaDisplay string `form:"nama_display" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Kode Display dan Nama Display harus diisi",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	result := h.DB.Table("bw_display_poli").Where("kd_display = ?", input.KdDisplay).Updates(map[string]interface{}{
		"nama_display": input.NamaDisplay,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat update display poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Display berhasil diupdate!",
		"color":   "success",
		"icon":    "check",
	})
}

// DeleteDisplay menghapus display poli yang ada
func (h *SettingDisplayPoliHandler) DeleteDisplay(c *gin.Context) {
	kdDisplay := c.Param("kd_display")

	result := h.DB.Table("bw_display_poli").Where("kd_display = ?", kdDisplay).Delete(nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat menghapus display poli",
			"color":   "danger",
			"icon":    "ban",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Display berhasil dihapus!",
		"color":   "warning",
		"icon":    "check",
	})
}

// getAllDisplay mendapatkan daftar semua display
func (h *SettingDisplayPoliHandler) getAllDisplay() []map[string]interface{} {
	var results []map[string]interface{}
	h.DB.Table("bw_display_poli").
		Select("bw_display_poli.kd_display, bw_display_poli.nama_display").
		Find(&results)

	return results
}
