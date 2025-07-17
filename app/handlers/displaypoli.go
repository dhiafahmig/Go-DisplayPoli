package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/dhiafahmig/Go-DisplayPoli/app/services"
)

// DisplayPoliHandler menangani tampilan display poli
type DisplayPoliHandler struct {
	DB *gorm.DB
}

// NewDisplayPoliHandler membuat instance baru dari DisplayPoliHandler
func NewDisplayPoliHandler(db *gorm.DB) *DisplayPoliHandler {
	return &DisplayPoliHandler{DB: db}
}

// HandleDisplay menangani permintaan untuk menampilkan display poli
func (h *DisplayPoliHandler) HandleDisplay(c *gin.Context) {
	kdDisplay := c.Param("kd_display")
	poliList := h.getPoliList(kdDisplay)

	c.HTML(http.StatusOK, "displaypoli.html", gin.H{
		"KdDisplay": kdDisplay,
		"PoliList":  poliList,
		"PusherKey": services.GetPusherKey(),
		"AppURL":    services.GetAppURL(),
	})
}

// GetPoliListByDisplay mendapatkan daftar poli untuk API
func (h *DisplayPoliHandler) GetPoliListByDisplay(c *gin.Context) {
	kdDisplay := c.Param("kd_display")
	poliList := h.getPoliList(kdDisplay)
	c.JSON(http.StatusOK, poliList)
}

// getPoliList mendapatkan daftar poli berdasarkan kode display
func (h *DisplayPoliHandler) getPoliList(kdDisplay string) []map[string]interface{} {
	var results []map[string]interface{}

	h.DB.Table("bw_ruang_poli").
		Select("bw_ruang_poli.kd_ruang_poli, bw_ruang_poli.nama_ruang_poli, bw_ruang_poli.kd_display, bw_ruang_poli.posisi_display_poli").
		Where("bw_ruang_poli.kd_display = ?", kdDisplay).
		Order("bw_ruang_poli.posisi_display_poli asc").
		Find(&results)

	// Mengambil pasien untuk setiap poli
	for i := range results {
		hari := services.GetDayList()[time.Now().Format("Monday")]
		var pasienList []map[string]interface{}
		var missedPatients []map[string]interface{}

		// Pertama cek apakah ada pasien yang sedang dipanggil (status=2)
		h.DB.Table("reg_periksa").
			Select("reg_periksa.no_reg, reg_periksa.no_rawat, bw_ruangpoli_dokter.nama_dokter, jadwal.hari_kerja, jadwal.jam_mulai, bw_ruangpoli_dokter.kd_ruang_poli, pasien.nm_pasien, reg_periksa.kd_pj").
			Joins("JOIN bw_ruangpoli_dokter ON reg_periksa.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
			Joins("JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
			Joins("JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis").
			Joins("JOIN bw_log_antrian_poli ON reg_periksa.no_rawat = bw_log_antrian_poli.no_rawat").
			Where("reg_periksa.tgl_registrasi = ?", time.Now().Format("2006-01-02")).
			Where("jadwal.hari_kerja = ?", hari).
			Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", results[i]["kd_ruang_poli"]).
			Where("bw_log_antrian_poli.status = '2'"). // Status 2 = sedang dipanggil
			Order("jadwal.jam_mulai asc").
			Order("reg_periksa.no_reg asc").
			Order("reg_periksa.jam_reg asc").
			Limit(1).
			Find(&pasienList)

		// Jika tidak ada pasien yang sedang dipanggil, ambil pasien berikutnya seperti biasa
		if len(pasienList) == 0 {
			h.DB.Table("reg_periksa").
				Select("reg_periksa.no_reg, reg_periksa.no_rawat, bw_ruangpoli_dokter.nama_dokter, jadwal.hari_kerja, jadwal.jam_mulai, bw_ruangpoli_dokter.kd_ruang_poli, pasien.nm_pasien, reg_periksa.kd_pj").
				Joins("JOIN bw_ruangpoli_dokter ON reg_periksa.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
				Joins("JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
				Joins("JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis").
				Where("reg_periksa.tgl_registrasi = ?", time.Now().Format("2006-01-02")).
				Where("jadwal.hari_kerja = ?", hari).
				Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", results[i]["kd_ruang_poli"]).
				Where("NOT EXISTS (SELECT 1 FROM bw_log_antrian_poli WHERE reg_periksa.no_rawat = bw_log_antrian_poli.no_rawat)").
				Order("jadwal.jam_mulai asc").
				Order("reg_periksa.no_reg asc").
				Order("reg_periksa.jam_reg asc").
				Limit(1).
				Find(&pasienList)
		}

		// Ambil daftar pasien yang terlewat (status=1)
		h.DB.Table("reg_periksa").
			Select("reg_periksa.no_reg, reg_periksa.no_rawat, bw_ruangpoli_dokter.nama_dokter, "+
				"jadwal.hari_kerja, jadwal.jam_mulai, bw_ruangpoli_dokter.kd_ruang_poli, "+
				"pasien.nm_pasien, reg_periksa.kd_pj, bw_log_antrian_poli.status").
			Joins("JOIN bw_ruangpoli_dokter ON reg_periksa.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
			Joins("JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
			Joins("JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis").
			Joins("JOIN bw_log_antrian_poli ON reg_periksa.no_rawat = bw_log_antrian_poli.no_rawat").
			Where("reg_periksa.tgl_registrasi = ?", time.Now().Format("2006-01-02")).
			Where("jadwal.hari_kerja = ?", hari).
			Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", results[i]["kd_ruang_poli"]).
			Where("bw_log_antrian_poli.status = '1'"). // Status 1 = terlewat
			Order("jadwal.jam_mulai asc").
			Order("reg_periksa.no_reg asc").
			Order("reg_periksa.jam_reg asc").
			Find(&missedPatients)

		if len(pasienList) > 0 {
			results[i]["getPasien"] = pasienList
		} else {
			results[i]["getPasien"] = []map[string]interface{}{}
		}

		// Tambahkan daftar pasien yang terlewat
		if len(missedPatients) > 0 {
			results[i]["missedPatients"] = missedPatients
		} else {
			results[i]["missedPatients"] = []map[string]interface{}{}
		}
	}

	return results
}

// GetMissedPatients mendapatkan daftar pasien yang terlewat untuk poli tertentu
func (h *DisplayPoliHandler) GetMissedPatients(c *gin.Context) {
	kdRuangPoli := c.Param("kd_ruang_poli")
	hari := services.GetDayList()[time.Now().Format("Monday")]

	var missedPatients []map[string]interface{}

	h.DB.Table("reg_periksa").
		Select("reg_periksa.no_reg, reg_periksa.no_rawat, bw_ruangpoli_dokter.nama_dokter, "+
			"jadwal.hari_kerja, jadwal.jam_mulai, bw_ruangpoli_dokter.kd_ruang_poli, "+
			"pasien.nm_pasien, reg_periksa.kd_pj, bw_log_antrian_poli.status").
		Joins("JOIN bw_ruangpoli_dokter ON reg_periksa.kd_dokter = bw_ruangpoli_dokter.kd_dokter").
		Joins("JOIN jadwal ON bw_ruangpoli_dokter.kd_dokter = jadwal.kd_dokter").
		Joins("JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis").
		Joins("JOIN bw_log_antrian_poli ON reg_periksa.no_rawat = bw_log_antrian_poli.no_rawat").
		Where("reg_periksa.tgl_registrasi = ?", time.Now().Format("2006-01-02")).
		Where("jadwal.hari_kerja = ?", hari).
		Where("bw_ruangpoli_dokter.kd_ruang_poli = ?", kdRuangPoli).
		Where("bw_log_antrian_poli.status = '1'"). // Status 1 = terlewat
		Order("jadwal.jam_mulai asc").
		Order("reg_periksa.no_reg asc").
		Order("reg_periksa.jam_reg asc").
		Find(&missedPatients)

	if len(missedPatients) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"missed_patients": missedPatients,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"missed_patients": []map[string]interface{}{},
			},
		})
	}
}
