package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/dhiafahmig/Go-DisplayPoli/app/handlers"
)

var (
	db       *gorm.DB
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Log origin untuk debugging
			origin := r.Header.Get("Origin")
			log.Printf("WebSocket connection request from Origin: %s", origin)
			return true // Allow all origins for WebSocket
		},
	}
	clients     = make(map[*websocket.Conn]bool)
	broadcaster = make(chan handlers.PanggilPoliMessage)
)

func init() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Setup database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func main() {
	// Aktifkan mode debug untuk log yang lebih detail
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// Tambahkan template functions sebelum load template
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	r.Static("/assets", "./assets")

	// Konfigurasi CORS dengan gin-contrib/cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inisialisasi handlers
	displayPoliHandler := handlers.NewDisplayPoliHandler(db)
	settingDisplayPoliHandler := handlers.NewSettingDisplayPoliHandler(db)
	settingPoliHandler := handlers.NewSettingPoliHandler(db)
	settingPosisiDokterHandler := handlers.NewSettingPosisiDokterHandler(db)
	jadwalDokterHandler := handlers.NewJadwalDokterHandler(db)
	panggilPoliHandler := handlers.NewPanggilPoliHandler(db)
	panggilPoliHandler.SetBroadcaster(broadcaster)

	// Memulai broadcaster
	go handleMessages()

	// Rutekan API Halaman
	r.GET("/ws/:kd_display", handleWebsocket)
	r.GET("/ws/antrian/:kd_ruang_poli", panggilPoliHandler.HandleAntrianWebSocket)
	r.GET("/display/:kd_display", displayPoliHandler.HandleDisplay)
	r.GET("/settings/display", settingDisplayPoliHandler.HandleSettings)
	r.GET("/settings/poli", settingPoliHandler.HandleSettings)
	r.GET("/settings/dokter", settingPosisiDokterHandler.HandleSettings)
	r.GET("/jadwal/dokter", jadwalDokterHandler.HandleJadwal)
	r.GET("/panggilpoli/:kd_ruang_poli/:kd_display", panggilPoliHandler.HandlePanggil)

	// API untuk aplikasi React
	r.GET("/api/display/poli/:kd_display", displayPoliHandler.GetPoliListByDisplay)
	r.GET("/api/panggil/:kd_ruang_poli", panggilPoliHandler.HandlePanggilAPI)
	r.GET("/api/antrian/poli/:kd_ruang_poli", panggilPoliHandler.HandleAntrianPoliAPI)
	r.GET("/api/jadwal/dokter/all", jadwalDokterHandler.GetJadwalDokter) // Route baru untuk mendapatkan jadwal dokter

	// API untuk pengaturan display
	displayGroup := r.Group("/api/display")
	{
		displayGroup.GET("/", settingDisplayPoliHandler.GetAllDisplay)
		displayGroup.POST("/", settingDisplayPoliHandler.AddDisplay)
		displayGroup.PUT("/", settingDisplayPoliHandler.EditDisplay)
		displayGroup.DELETE("/:kd_display", settingDisplayPoliHandler.DeleteDisplay)
	}

	// API untuk pengaturan poli
	poliGroup := r.Group("/api/poli")
	{
		poliGroup.GET("/", settingPoliHandler.GetAllPoli)
		poliGroup.POST("/", settingPoliHandler.AddPoli)
		poliGroup.PUT("/", settingPoliHandler.EditPoli)
		poliGroup.DELETE("/:kd_ruang_poli", settingPoliHandler.DeletePoli)
		poliGroup.GET("/dokter/:kd_ruang_poli", settingPoliHandler.GetDokterPoli)
	}

	// API untuk pengaturan posisi dokter
	dokterPoliGroup := r.Group("/api/dokterpoli")
	{
		dokterPoliGroup.POST("/", settingPosisiDokterHandler.EditPoliDokter)
	}

	// API untuk jadwal dokter
	jadwalGroup := r.Group("/api/jadwal")
	{
		jadwalGroup.GET("/dokter", jadwalDokterHandler.CariDokter)
		jadwalGroup.POST("/", jadwalDokterHandler.TambahJadwalDokter)
		jadwalGroup.PUT("/", jadwalDokterHandler.UbahJadwalDokter)
		jadwalGroup.DELETE("/", jadwalDokterHandler.HapusJadwalDokter)
	}

	// API untuk antrian pasien
	r.POST("/api/panggilpoli", panggilPoliHandler.PanggilPasien)
	r.POST("/api/panggilpasien", panggilPoliHandler.PanggilPasien)
	r.POST("/api/antrian/panggil", panggilPoliHandler.PanggilPasienAPI)

	r.POST("/api/log", panggilPoliHandler.HandleLog)
	r.POST("/api/log/reset/:no_rawat", panggilPoliHandler.ResetLog)
	r.POST("/api/antrian/log", panggilPoliHandler.HandleLogAPI)
	r.POST("/api/antrian/log/reset/:no_rawat", panggilPoliHandler.ResetLogAPI)

	// Serve aplikasi React
	r.NoRoute(func(c *gin.Context) {
		// Coba gunakan app build React jika tersedia
		path := "./react-app/build" + c.Request.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.File("./react-app/build/index.html")
			return
		}
		c.File(path)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handleWebsocket(c *gin.Context) {
	kdDisplay := c.Param("kd_display")

	// Log headers untuk debugging
	log.Printf("WebSocket connection attempt for display: %s", kdDisplay)
	log.Printf("Headers: %+v", c.Request.Header)

	// Mencoba upgrade dengan error handling yang lebih baik
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Log successful connection
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("WebSocket connection established from %s for display: %s", remoteAddr, kdDisplay)

	defer func() {
		conn.Close()
		log.Printf("WebSocket connection closed for %s", remoteAddr)
	}()

	clients[conn] = true

	// Send initial message to confirm connection
	initialMsg := handlers.PanggilPoliMessage{
		KdDisplay:   kdDisplay,
		KdRuangPoli: "INITIAL",   // Tambahkan field yang diperlukan
		NmPasien:    "Connected", // Tambahkan informasi koneksi
		NoReg:       "0",
	}

	if err := conn.WriteJSON(initialMsg); err != nil {
		log.Printf("Error sending initial message: %v", err)
	}

	for {
		// Keep connection alive
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed for %s: %v", remoteAddr, err)
			delete(clients, conn)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcaster
		log.Printf("Broadcasting message for display %s, poli %s", msg.KdDisplay, msg.KdRuangPoli)

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
