package services

// GetDayList mengembalikan pemetaan nama hari dalam bahasa Inggris ke bahasa Indonesia
func GetDayList() map[string]string {
	return map[string]string{
		"Sunday":    "AKHAD",
		"Monday":    "SENIN",
		"Tuesday":   "SELASA",
		"Wednesday": "RABU",
		"Thursday":  "KAMIS",
		"Friday":    "JUMAT",
		"Saturday":  "SABTU",
	}
}
