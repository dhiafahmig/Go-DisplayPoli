package handlers

// PanggilPoliMessage mewakili struktur pesan untuk memanggil pasien
type PanggilPoliMessage struct {
	NmPasien    string `json:"nm_pasien"`
	KdRuangPoli string `json:"kd_ruang_poli"`
	NmPoli      string `json:"nm_poli"`
	NoReg       string `json:"no_reg"`
	KdDisplay   string `json:"kd_display"`
}
