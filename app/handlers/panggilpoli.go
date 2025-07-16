package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	htgotts "github.com/hegedustibor/htgo-tts"
	"gorm.io/gorm"

	"github.com/dhiafahmig/Go-DisplayPoli/app/services"
)

// PanggilPoliHandler menangani tampilan memanggil pasien di poli
type PanggilPoliHandler struct {
	DB          *gorm.DB
	Broadcaster chan<- PanggilPoliMessage // Channel untuk broadcast pesan
}

// NewPanggilPoliHandler membuat instance baru dari PanggilPoliHandler
func NewPanggilPoliHandler(db *gorm.DB) *PanggilPoliHandler {
	return &PanggilPoliHandler{DB: db}
}

// SetBroadcaster menetapkan channel broadcaster untuk handler ini
func (h *PanggilPoliHandler) SetBroadcaster(broadcaster chan<- PanggilPoliMessage) {
	h.Broadcaster = broadcaster
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

// generateTTS menghasilkan file audio dari teks dan mengembalikan URL relatif
func (h *PanggilPoliHandler) generateTTS(text, kdRuangPoli, noReg string) (string, error) {
	// Buat hash dari teks untuk nama file unik
	hasher := md5.New()
	hasher.Write([]byte(text + kdRuangPoli + noReg + time.Now().String()))
	filename := hex.EncodeToString(hasher.Sum(nil)) + ".mp3"

	// Pastikan direktori audio ada
	audioDir := "assets/audio"
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		if err := os.MkdirAll(audioDir, 0755); err != nil {
			return "", err
		}
	}

	// Buat instance TTS
	speech := htgotts.Speech{
		Folder:   audioDir,
		Language: "id",
	}

	// Generate audio file
	speech.CreateSpeechFile(text, filename)

	// Kembalikan URL relatif
	return "/assets/audio/" + filename, nil
}

// PanggilPasien mengirim event untuk memanggil pasien
func (h *PanggilPoliHandler) PanggilPasien(c *gin.Context) {
	var input struct {
		NmPasien    string `json:"nm_pasien" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli" binding:"required"`
		NmPoli      string `json:"nm_poli" binding:"required"`
		NoReg       string `json:"no_reg" binding:"required"`
		KdDisplay   string `json:"kd_display" binding:"required"`
		NoRawat     string `json:"no_rawat" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat teks untuk TTS
	ttsText := fmt.Sprintf("Nomor antrian %s, atas nama %s, silakan menuju %s",
		input.NoReg, input.NmPasien, input.NmPoli)

	// Generate file audio TTS
	audioUrl, err := h.generateTTS(ttsText, input.KdRuangPoli, input.NoReg)
	if err != nil {
		log.Printf("Error generating TTS: %v", err)
		// Lanjutkan meskipun TTS gagal
	}

	// Buat pesan untuk dikirim ke websocket
	msg := PanggilPoliMessage{
		NmPasien:    input.NmPasien,
		KdRuangPoli: input.KdRuangPoli,
		NmPoli:      input.NmPoli,
		NoReg:       input.NoReg,
		KdDisplay:   input.KdDisplay,
		AudioUrl:    audioUrl,
	}

	// Kirim ke broadcaster jika tersedia
	if h.Broadcaster != nil {
		log.Printf("Mengirim pesan panggil ke broadcaster: %+v", msg)
		h.Broadcaster <- msg
	} else {
		log.Printf("Broadcaster tidak tersedia, tidak bisa mengirim pesan: %+v", msg)
	}

	// Tambahkan: Update status pasien di bw_log_antrian_poli untuk menandai pasien yang dipanggil
	// Status "2" menandakan pasien sedang dipanggil
	if input.NoRawat != "" {
		// Hapus status "sedang dipanggil" (2) dari semua pasien di poli yang sama
		h.DB.Exec(`
			DELETE FROM bw_log_antrian_poli 
			WHERE kd_ruang_poli = ? AND status = '2'
		`, input.KdRuangPoli)

		// Tambahkan status baru untuk pasien yang dipanggil
		result := h.DB.Exec(`
			INSERT INTO bw_log_antrian_poli (no_rawat, kd_ruang_poli, status)
			VALUES (?, ?, '2')
			ON DUPLICATE KEY UPDATE kd_ruang_poli = ?, status = '2'
		`, input.NoRawat, input.KdRuangPoli, input.KdRuangPoli)

		if result.Error != nil {
			log.Printf("Error updating patient status: %v", result.Error)
		}

		// Jadwalkan reset status setelah 5 menit
		go func(noRawat string) {
			time.Sleep(5 * time.Minute)
			h.resetCallingStatus(noRawat)
		}(input.NoRawat)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Panggilan berhasil dikirim",
		"data":    msg,
	})
}

// resetCallingStatus mengembalikan status pasien dari "sedang dipanggil" (2) menjadi normal
func (h *PanggilPoliHandler) resetCallingStatus(noRawat string) {
	// Hapus status panggilan setelah 5 menit
	result := h.DB.Table("bw_log_antrian_poli").
		Where("no_rawat = ? AND status = '2'", noRawat).
		Delete(nil)

	if result.Error != nil {
		log.Printf("Error resetting calling status: %v", result.Error)
	} else {
		log.Printf("Successfully reset calling status for patient: %s", noRawat)
	}
}

// PanggilPasienAPI adalah API khusus untuk memanggil pasien di antrian
func (h *PanggilPoliHandler) PanggilPasienAPI(c *gin.Context) {
	var input struct {
		NmPasien    string `json:"nm_pasien" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli" binding:"required"`
		NmPoli      string `json:"nm_poli" binding:"required"`
		NoReg       string `json:"no_reg" binding:"required"`
		KdDisplay   string `json:"kd_display" binding:"required"`
		NoRawat     string `json:"no_rawat" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
		})
		return
	}

	// Buat teks untuk TTS
	ttsText := fmt.Sprintf("Nomor antrian %s, atas nama %s, silakan menuju %s",
		input.NoReg, input.NmPasien, input.NmPoli)

	// Generate file audio TTS
	audioUrl, err := h.generateTTS(ttsText, input.KdRuangPoli, input.NoReg)
	if err != nil {
		log.Printf("Error generating TTS: %v", err)
		// Lanjutkan meskipun TTS gagal
	}

	// Buat pesan untuk dikirim ke websocket
	msg := PanggilPoliMessage{
		NmPasien:    input.NmPasien,
		KdRuangPoli: input.KdRuangPoli,
		NmPoli:      input.NmPoli,
		NoReg:       input.NoReg,
		KdDisplay:   input.KdDisplay,
		AudioUrl:    audioUrl,
	}

	// Kirim ke broadcaster jika tersedia
	if h.Broadcaster != nil {
		log.Printf("Mengirim pesan panggil ke broadcaster: %+v", msg)
		h.Broadcaster <- msg
	} else {
		log.Printf("Broadcaster tidak tersedia, tidak bisa mengirim pesan: %+v", msg)
	}

	// Tambahkan: Update status pasien di bw_log_antrian_poli untuk menandai pasien yang dipanggil
	// Status "2" menandakan pasien sedang dipanggil
	if input.NoRawat != "" {
		// Hapus status "sedang dipanggil" (2) dari semua pasien di poli yang sama
		h.DB.Exec(`
			DELETE FROM bw_log_antrian_poli 
			WHERE kd_ruang_poli = ? AND status = '2'
		`, input.KdRuangPoli)

		// Tambahkan status baru untuk pasien yang dipanggil
		result := h.DB.Exec(`
			INSERT INTO bw_log_antrian_poli (no_rawat, kd_ruang_poli, status)
			VALUES (?, ?, '2')
			ON DUPLICATE KEY UPDATE kd_ruang_poli = ?, status = '2'
		`, input.NoRawat, input.KdRuangPoli, input.KdRuangPoli)

		if result.Error != nil {
			log.Printf("Error updating patient status: %v", result.Error)
		}

		// Jadwalkan reset status setelah 5 menit
		go func(noRawat string) {
			time.Sleep(5 * time.Minute)
			h.resetCallingStatus(noRawat)
		}(input.NoRawat)
	}

	// Mengembalikan response dalam format JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": msg,
		},
		"message": "Pasien berhasil dipanggil",
	})
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

// HandleAntrianWebSocket menangani koneksi WebSocket untuk pembaruan antrian
func (h *PanggilPoliHandler) HandleAntrianWebSocket(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")

	// Upgrade koneksi HTTP ke WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Log successful connection
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("WebSocket antrian connection established from %s for poli: %s", remoteAddr, kdRuangPoli)

	// Daftarkan koneksi ke channel broadcast
	// TODO: Implementasikan sistem broadcast

	// Send initial message to confirm connection
	initialMsg := PanggilPoliMessage{
		KdRuangPoli: kdRuangPoli,
		NmPasien:    "Connected",
		NoReg:       "0",
	}

	if err := conn.WriteJSON(initialMsg); err != nil {
		log.Printf("Error sending initial antrian message: %v", err)
	}

	// Tangani pesan yang masuk
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket antrian connection closed for %s: %v", remoteAddr, err)
			break
		}
	}
}

// HandlePanggilAPI menangani permintaan API dari frontend React untuk halaman panggil poli
func (h *PanggilPoliHandler) HandlePanggilAPI(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")
	pasienList := h.getPasienList(kdRuangPoli)

	// Dapatkan informasi poli
	var poliInfo map[string]interface{}
	h.DB.Table("bw_ruangpoli").
		Select("kd_ruang_poli, nama_ruang_poli").
		Where("kd_ruang_poli = ?", kdRuangPoli).
		First(&poliInfo)

	// Kirim sebagai JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"poli_info":   poliInfo,
			"pasien_list": pasienList,
		},
		"message": "Data pasien berhasil dimuat",
	})
}

// HandleAntrianPoliAPI menangani permintaan API untuk mendapatkan daftar antrian pada poli tertentu
func (h *PanggilPoliHandler) HandleAntrianPoliAPI(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")

	// Mendapatkan data dari database
	pasienList := h.getPasienList(kdRuangPoli)

	// Mendapatkan informasi poli
	var poliInfo map[string]interface{}
	h.DB.Table("bw_ruangpoli").
		Select("kd_ruang_poli, nama_ruang_poli").
		Where("kd_ruang_poli = ?", kdRuangPoli).
		First(&poliInfo)

	// Mengembalikan response dalam format JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"poli_info": poliInfo,
			"antrian":   pasienList,
		},
		"message": "Data antrian berhasil diambil",
	})
}

// HandleLogAPI menangani API untuk mengupdate status log antrian pasien
func (h *PanggilPoliHandler) HandleLogAPI(c *gin.Context) {
	var input struct {
		KdDokter    string `json:"kd_dokter"`
		NoRawat     string `json:"no_rawat" binding:"required"`
		KdRuangPoli string `json:"kd_ruang_poli" binding:"required"`
		Type        string `json:"type" binding:"required"` // 'ada' atau 'tidak'
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengupdate status: " + result.Error.Error(),
		})
		return
	}

	// Mengembalikan response dalam format JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"no_rawat":      input.NoRawat,
			"kd_ruang_poli": input.KdRuangPoli,
			"status":        status,
		},
		"message": "Status pasien berhasil diperbarui",
	})
}

// ResetLogAPI menangani API untuk menghapus log antrian pasien
func (h *PanggilPoliHandler) ResetLogAPI(c *gin.Context) {
	noRawat := c.Param("no_rawat")

	result := h.DB.Table("bw_log_antrian_poli").Where("no_rawat = ?", noRawat).Delete(nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mereset status: " + result.Error.Error(),
		})
		return
	}

	// Mengembalikan response dalam format JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"no_rawat": noRawat,
		},
		"message": "Reset log berhasil",
	})
}
