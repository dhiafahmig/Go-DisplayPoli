package handlers

import (
	"log"
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

// GetAllPoli mengembalikan daftar semua poli dalam format JSON
func (h *SettingPoliHandler) GetAllPoli(c *gin.Context) {
	poliList := h.getAllPoli()
	c.JSON(http.StatusOK, poliList)
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

// GetDokterPoli menangani permintaan untuk mendapatkan daftar dokter pada poli tertentu
func (h *SettingPoliHandler) GetDokterPoli(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")

	// Log untuk debugging
	log.Printf("Mencari dokter untuk poli: %s", kdRuangPoli)

	// Tampilkan daftar semua poli yang tersedia (untuk debugging)
	var semuaPoli []map[string]interface{}
	h.DB.Table("bw_ruang_poli").Find(&semuaPoli)
	log.Printf("Daftar semua poli yang tersedia: %+v", semuaPoli)

	// Periksa apakah poli ada
	var poliCount int64
	h.DB.Table("bw_ruang_poli").Where("kd_ruang_poli = ?", kdRuangPoli).Count(&poliCount)
	log.Printf("Jumlah poli dengan kode %s: %d", kdRuangPoli, poliCount)

	// Coba tampilkan juga poli yang serupa (jika poli tidak ditemukan)
	if poliCount == 0 {
		var poliSimilar []map[string]interface{}
		h.DB.Table("bw_ruang_poli").Limit(5).Find(&poliSimilar)
		log.Printf("Poli tidak ditemukan. Beberapa poli yang tersedia: %+v", poliSimilar)

		// Tetap lanjutkan karena kita ingin mengembalikan data kosong, bukan error
	}

	// Mendapatkan informasi poli
	var poliInfo map[string]interface{}
	poliResult := h.DB.Table("bw_ruang_poli").
		Select("kd_ruang_poli, nama_ruang_poli").
		Where("kd_ruang_poli = ?", kdRuangPoli).
		First(&poliInfo)

	if poliResult.Error != nil {
		log.Printf("Error saat mengambil informasi poli: %v", poliResult.Error)
		// Tetap lanjutkan dengan nilai default
		poliInfo = map[string]interface{}{
			"kd_ruang_poli":   kdRuangPoli,
			"nama_ruang_poli": "Tidak ditemukan",
		}
	} else {
		log.Printf("Informasi poli: %+v", poliInfo)
	}

	// Mendapatkan daftar dokter dari database
	var dokters []map[string]interface{}

	// Tampilkan juga tabel bw_ruangpoli_dokter untuk debugging
	var allDokters []map[string]interface{}
	h.DB.Table("bw_ruangpoli_dokter").Limit(5).Find(&allDokters)
	log.Printf("Sampel data tabel bw_ruangpoli_dokter: %+v", allDokters)

	// Query untuk mendapatkan dokter dengan nm_poli dari tabel poliklinik
	result := h.DB.Table("bw_ruangpoli_dokter").
		Select("bw_ruangpoli_dokter.kd_dokter, bw_ruangpoli_dokter.nama_dokter, dokter.jk, bw_ruangpoli_dokter.kd_ruang_poli, bw_ruang_poli.nama_ruang_poli, poliklinik.nm_poli").
		Joins("LEFT JOIN dokter ON bw_ruangpoli_dokter.kd_dokter = dokter.kd_dokter").
		Joins("LEFT JOIN bw_ruang_poli ON bw_ruangpoli_dokter.kd_ruang_poli = bw_ruang_poli.kd_ruang_poli").
		Joins("LEFT JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
		Joins("LEFT JOIN poliklinik ON jadwal.kd_poli = poliklinik.kd_poli").
		Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", kdRuangPoli).
		Order("bw_ruangpoli_dokter.nama_dokter ASC").
		Find(&dokters)

	if result.Error != nil {
		log.Printf("Error saat mengambil data dokter: %v", result.Error)
		// Tetap lanjutkan dengan array kosong
		dokters = []map[string]interface{}{}
	}

	log.Printf("Jumlah dokter ditemukan: %d", len(dokters))
	if len(dokters) > 0 {
		log.Printf("Contoh data dokter pertama: %+v", dokters[0])
	}

	// Mengembalikan response dalam format JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"poli_info": poliInfo,
			"dokters":   dokters,
		},
		"message": "Data dokter berhasil diambil",
	})
}
