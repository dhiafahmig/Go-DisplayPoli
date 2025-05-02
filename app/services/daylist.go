package services

// DayList menyediakan mapping dari nama hari dalam bahasa Inggris ke format bahasa Indonesia
var DayList = map[string]string{
	"Sunday":    "AKHAD",
	"Monday":    "SENIN",
	"Tuesday":   "SELASA",
	"Wednesday": "RABU",
	"Thursday":  "KAMIS",
	"Friday":    "JUMAT",
	"Saturday":  "SABTU",
}

// GetDayList mengembalikan mapping hari
func GetDayList() map[string]string {
	return DayList
}
