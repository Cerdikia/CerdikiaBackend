package main

import (
	"coba1BE/config"
	"coba1BE/controllers"
	"coba1BE/middleware"
	"coba1BE/services"
	"fmt"
	"os"

	"time"

	// "fmt"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("PORT")

	r := gin.Default()

	r.Static("/uploads", "./uploads")
	r.Static("/tmp", "./tmp")

	// Atur middleware CORS
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"}, // sesuaikan origin frontend kamu
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = services.StartEnergyCron()
	if err != nil {
		log.Fatal("Failed to Start Cron:", err)
	}

	r.POST("/register/:role", controllers.CreateUser)

	r.POST("/login", controllers.LoginControler)
	r.POST("/refresh", controllers.RefreshToken)

	protected := r.Group("")
	protected.GET("/try", controllers.GetUsers)
	protected.Use(middleware.AuthMiddleware())
	// BUT HANDLER POST MAPEL
	// BUT HANDLER EDIT MAPEL
	// BUT HANDLER DELETE MAPEL
	protected.GET("/genericAllMapels", controllers.GetAllMapel)
	protected.GET("/genericMapel/:id", controllers.GetMapelByID)
	protected.GET("/genericMapels", controllers.CGenericMapels) // Query param : ?id_kelas=, ?finished= ambil semua mata pelajaran dan jumlah modulnya melalui kelas
	protected.POST("/genericMapels", controllers.CreateMapel)
	protected.PUT("/genericMapels/:id", controllers.UpdateMapel)
	protected.DELETE("/genericMapels/:id", controllers.DeleteMapel)

	// BUT HANDLER POST MODULE
	// BUT HANDLER EDIT MODULE
	// BUT HANDLER DELETE MODULE
	protected.GET("/genericModules", controllers.CGenericModules) // ambil semua module dengan acuan  mapel
	// r.GET("/", controllers.GetAllModules)
	protected.POST("/genericModules", controllers.CreateModule)
	protected.PUT("/genericModules/:id", controllers.UpdateModule)
	protected.DELETE("/genericModules/:id", controllers.DeleteModule)
	protected.GET("/genericModule/:id", controllers.GetModuleByID)
	protected.PUT("/togle-module/:id_module", controllers.ToggleModuleReady)
	protected.GET("/stats", controllers.GetStats)                        // Endpoint lama untuk kompatibilitas
	protected.GET("/all-stats", controllers.GetAllStats)                 // Endpoint baru yang menampilkan semua statistik
	protected.GET("/recent-activities", controllers.GetRecentActivities) // Endpoint untuk menampilkan aktivitas terakhir

	// ========= SOAL ==============================
	protected.GET("/genericSoal/:id_module", controllers.CGenericSoal) // ambil soal dari sebuah module dengan acuan id_module

	// ========= Guru - Mapel Relasi ===============
	protected.GET("/guru/:id_guru", controllers.GetMapelByGuru)
	protected.POST("/guru_mapel", controllers.AddGuruMapel)
	protected.PUT("/guru_mapel", controllers.UpdateGuruMapel)
	protected.DELETE("/guru_mapel/:id_guru/:id_mapel", controllers.DeleteGuruMapel)
	protected.DELETE("/guru_mapel/:id_guru", controllers.DeleteAllGuruMapel)
	// Batch operations for multiple teachers at once
	protected.POST("/guru_mapel/batch", controllers.BatchAddGuruMapel)
	protected.PUT("/guru_mapel/batch", controllers.BatchUpdateGuruMapel)
	protected.DELETE("/guru_mapel/batch", controllers.BatchDeleteGuruMapel)

	// ========= Actor ===============
	protected.GET("/getAllUsers", controllers.GetUsers)
	// protected.GET("/siswa", controllers.GetSiswa)
	protected.GET("/getDataActor/:role/:email", controllers.GetDataActor)
	protected.GET("/getDataUser", controllers.GetUser)
	protected.PUT("/editDataUser/:role", controllers.UpdateDataActor)
	protected.DELETE("/deleteDataUser", controllers.DeleteDataActor)
	protected.PATCH("/patchImageProfile/:role/:email", controllers.UpdateSiswaImageProfile)
	protected.POST("/changeUserRole", controllers.ChangeUserRole)

	// ========= Siawa Verified ===============
	protected.GET("/verified", controllers.Beingverified)
	protected.GET("/verifiedes", controllers.Beingverifieds)
	protected.PATCH("/verifiedes", controllers.UpdateUserVerifiedBatch)

	// ========= Soal ===============
	protected.POST("/upload-image", controllers.UploadImage)
	protected.POST("/upload-soal", controllers.UploadSoal)
	// if strings.Contains(strings.ToLower(message), "success") {
	protected.GET("/getDataSoal/:id_soal", controllers.GetDataSoal)
	protected.PUT("/editDataSoal/:id_soal", controllers.UpdateDataSoal)
	protected.DELETE("/deleteDataSoal/:id_soal", controllers.DeleteSoal)

	// ========= Point ===============
	protected.GET("/point", controllers.GetPoint)
	protected.PUT("/point", controllers.UpdatePoint)

	// ========= Ranking ===============
	// get all ranking
	// protected.GET("/ranking", controllers.GetRanking)
	// get ranking kelas
	protected.GET("/ranking", controllers.GetRankingByKelas)

	// mbil status module udah di kerjain atau belum
	// r.GET("/modules", controllers.GetModulesWithCompletion)

	// ========= Barang ===============
	// CRUD routes
	// protected.GET("/barang", controllers.GetAllBarang)
	// protected.GET("/barang/:id", controllers.GetBarangByID)
	// protected.POST("/barang", controllers.CreateBarang)
	// protected.PUT("/barang/:id", controllers.UpdateBarang)
	// protected.DELETE("/barang/:id", controllers.DeleteBarang)

	// Endpoints for gift management with image upload
	protected.POST("/gifts", controllers.CreateGift)
	protected.GET("/gifts", controllers.GetAllGifts)
	protected.GET("/gifts/:id", controllers.GetGiftByID)
	protected.PUT("/gifts/:id", controllers.UpdateGift)
	protected.DELETE("/gifts/:id", controllers.DeleteGift)

	// =================== Tukara Point ====================================
	// Route tukar barang
	protected.POST("/tukar-barang", controllers.TukarBarang) // USER POINT BLM DI KURANGIN!!!

	// =================== logs ====================================
	protected.GET("/logs/:email", controllers.GetLogsByEmail) // DEV MODE ONLY !!!
	protected.GET("/gegeralLogs", controllers.GetAllLogs)     // get all logs
	// r.GET("/:id", controllers.GetLogByID)
	protected.POST("/logs", controllers.CreateLog)
	// r.PUT("/:id", controllers.UpdateLog)
	// r.DELETE("/:id", controllers.DeleteLog)
	protected.GET("/logs", controllers.GetLogsBydEmailWithToken) // get logs where token email
	protected.GET("/logsBy", controllers.GetLogsBy)              // GET /logs/email/john@example.com/module/2
	protected.GET("/logs-periode", controllers.GetLogsByPeriod)

	// =================== Rekap Smester ====================================
	protected.POST("/rekap-semester", controllers.RekapSemester)            // <= KEMUNGKINAN GA DI PAKE
	protected.POST("/rekap-semester-siswa", controllers.RekapSemesterSiswa) // <= KEMUNGKINAN GA DI PAKE
	protected.POST("/rekap-semester-all-siswa", controllers.RekapSemesterAllSiswa)
	protected.POST("/edit-tahun-ajaran", controllers.EditTahunAjaran)
	protected.GET("/test-logs", controllers.TestLogsRetrieval) // Test endpoint for debugging logs retrieval
	protected.GET("/rekap-semester-all", controllers.GetAllDataSiswa)
	protected.GET("/rekap-semester/:id_data", controllers.GetDataSiswa)
	protected.DELETE("/rekap-semester/:id_data", controllers.DeleteDataSiswa)
	// {
	// 	"tahun_ajaran_lama": "2025/225",
	// 	"tahun_ajaran_baru": "2025/2026"
	// }

	// =================== KLEAS CRUD ====================================
	protected.GET("/kelas", controllers.GetAllKelas)
	protected.GET("/kelas/:id", controllers.GetKelasByID)
	protected.POST("/kelas", controllers.CreateKelas)
	protected.PUT("/kelas/:id", controllers.UpdateKelas)
	protected.DELETE("/kelas/:id", controllers.DeleteKelas)

	// =================== ENERGY CRUD ====================================
	protected.GET("/user-energy/:email", controllers.GetUserEnergy)
	protected.POST("/user-energy/:email", controllers.UseEnergyForAll)
	protected.POST("/add-energy/:email", controllers.AddEnergyForAll)

	// =================== GIFT REDEMPTION ====================================
	protected.POST("/redeem-gifts", controllers.RedeemGifts)
	protected.GET("/redemptions", controllers.GetAllRedemptions)
	protected.GET("/redemptions/:id", controllers.GetRedemptionByID)
	protected.GET("/redemptions/code/:code", controllers.GetRedemptionByCode)
	protected.PUT("/redemptions/:id/status", controllers.UpdateRedemptionStatus)
	protected.DELETE("/redemptions/:id", controllers.DeleteRedemption)
	protected.GET("/view-receipt/:code", controllers.GetRedemptionReceipt)
	protected.GET("/print-receipt/:code", controllers.PrintRedemptionReceipt)

	// =================== LAPORAN NILAI ====================================
	protected.GET("/score-report", controllers.GetScoreReport)
	protected.GET("/score-report-summary", controllers.GetScoreReportSummary)
	protected.GET("/student-score-comparison", controllers.GetStudentScoreComparison)
	protected.GET("/score-progress", controllers.GetScoreProgress)
	protected.GET("/all-students-report", controllers.GetAllStudentsReport)

	// =================== CHAT/MESSAGES ENDPOINTS ====================================
	protected.POST("/messages", controllers.CreateMessages)
	protected.GET("/messages", controllers.GetAllMessages) // New endpoint with query filtering
	protected.GET("/messages/recipient/:dest", controllers.GetMessagesByRecipient)
	protected.GET("/messages/sender/:form", controllers.GetMessagesBySender)
	protected.GET("/messages/subject/:subject", controllers.GetMessagesBySubject)
	protected.GET("/messages/unread/count", controllers.CountUnreadMessages)             // Count unread messages for a user
	protected.GET("/messages/unread/count/all", controllers.CountAllUnreadMessagesAdmin) // Count all unread messages for admin
	protected.GET("/messages/:id", controllers.GetChatMessageByIDAndMarkAsRead)          // Get message by ID and mark as read
	protected.PATCH("/messages/:id/status", controllers.UpdateMessageStatus)
	protected.POST("/messages/:id/read", controllers.MarkMessageAsRead) // New endpoint to mark messages as read

	// =================== IMPORT SOAL ================================================
	r.POST("/import", controllers.ImportSoalHandler)
	r.POST("/save-imported-data", controllers.SaveImportedDataHandler)

	// Menutup koneksi database saat aplikasi berhenti
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close() // Pastikan koneksi ditutup saat aplikasi keluar
	}

	// if err := r.Run(":4321"); err != nil {
	if err := r.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
