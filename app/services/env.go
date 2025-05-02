package services

import "os"

// ValueENV menyediakan fungsi untuk mengakses nilai variabel lingkungan
type ValueENV struct{}

// GetPusherKey mengembalikan Pusher API key dari variabel lingkungan
func GetPusherKey() string {
	return os.Getenv("PUSHER_APP_KEY")
}

// GetPusherCluster mengembalikan Pusher cluster dari variabel lingkungan
func GetPusherCluster() string {
	return os.Getenv("PUSHER_APP_CLUSTER")
}

// GetAppURL mengembalikan URL aplikasi dari variabel lingkungan
func GetAppURL() string {
	return os.Getenv("APP_URL")
}
