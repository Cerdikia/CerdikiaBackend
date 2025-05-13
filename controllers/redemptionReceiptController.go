package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"

	"coba1BE/config"
	"coba1BE/models/redemption"

	"github.com/joho/godotenv"
)

// GetRedemptionReceipt returns HTML content for viewing a redemption receipt
func GetRedemptionReceipt(c *gin.Context) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	baseurl := os.Getenv("BASEURL")

	db := config.DB
	code := c.Param("code")

	// Get the redemption log
	var log redemption.LogsPenukaranPoint
	if err := db.Where("kode_penukaran = ?", code).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption with code %s not found", code),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption data",
				"error":   err.Error(),
			})
		}
		return
	}

	// Get item details
	var barang redemption.BarangInfo
	if err := db.Where("id_barang = ?", log.IDBarang).First(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving item details",
			"error":   err.Error(),
		})
		return
	}

	// Get student details
	type StudentInfo struct {
		Nama         string `gorm:"column:nama"`
		IDKelas      int    `gorm:"column:id_kelas"`
		ImageProfile string `gorm:"column:image_profile"`
	}

	var student StudentInfo
	if err := db.Table("siswa").Where("email = ?", log.Email).First(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving student details",
			"error":   err.Error(),
		})
		return
	}

	// Get class name
	var className string
	var class struct {
		Kelas string `gorm:"column:kelas"`
	}
	if err := db.Table("kelas").Select("kelas").Where("id_kelas = ?", student.IDKelas).First(&class).Error; err == nil {
		className = class.Kelas
	}

	// Format date
	formattedDate := log.TanggalPenukaran.Format("02 January 2006 15:04")

	// Generate HTML content
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Bukti Penukaran - %s</title>
		<style>
			.receipt {
				max-width: 800px;
				margin: 0 auto;
				background-color: white;
				padding: 30px;
				border-radius: 8px;
				box-shadow: 0 2px 10px rgba(0,0,0,0.1);
			}
			.header {
				text-align: center;
				margin-bottom: 30px;
				padding-bottom: 20px;
				border-bottom: 2px solid #f0f0f0;
			}
			.logo {
				font-size: 24px;
				font-weight: bold;
				color: #4a4a4a;
				margin-bottom: 10px;
			}
			.title {
				font-size: 28px;
				color: #333;
				margin: 10px 0;
			}
			.code {
				font-size: 20px;
				color: #e74c3c;
				font-weight: bold;
				margin: 10px 0;
				letter-spacing: 2px;
			}
			.student-info {
				display: flex;
				margin-bottom: 30px;
			}
			.student-photo {
				width: 100px;
				height: 100px;
				border-radius: 50%%;
				object-fit: cover;
				margin-right: 20px;
				border: 3px solid #f0f0f0;
			}
			.student-details {
				flex: 1;
			}
			.student-details h3 {
				margin: 0 0 10px 0;
				color: #333;
			}
			.student-details p {
				margin: 5px 0;
				color: #666;
			}
			.item-details {
				margin-bottom: 30px;
				padding: 20px;
				background-color: #f9f9f9;
				border-radius: 5px;
			}
			.item-details h3 {
				margin: 0 0 15px 0;
				color: #333;
			}
			.item-row {
				display: flex;
				justify-content: space-between;
				margin-bottom: 10px;
			}
			.item-label {
				font-weight: bold;
				color: #555;
			}
			.footer {
				text-align: center;
				margin-top: 30px;
				padding-top: 20px;
				border-top: 2px solid #f0f0f0;
				color: #888;
				font-size: 14px;
			}
			.status {
				display: inline-block;
				padding: 5px 10px;
				border-radius: 4px;
				font-weight: bold;
			}
			.status-waiting {
				background-color: #f39c12;
				color: white;
			}
			.status-completed {
				background-color: #2ecc71;
				color: white;
			}
			.status-cancelled {
				background-color: #e74c3c;
				color: white;
			}
			.print-button {
				display: block;
				width: 200px;
				margin: 20px auto;
				padding: 10px;
				background-color: #3498db;
				color: white;
				border: none;
				border-radius: 4px;
				font-size: 16px;
				cursor: pointer;
				text-align: center;
				text-decoration: none;
			}
			.print-button:hover {
				background-color: #2980b9;
			}
		</style>
	</head>
	<body>
		<div class="receipt">
			<div class="header">
				<div class="logo">Cerdikia</div>
				<h1 class="title">Bukti Penukaran Hadiah</h1>
				<div class="code">%s</div>
			</div>
			
			<div class="student-info">
				<img src="%s" alt="Foto Profil" class="student-photo" onerror="this.src='https://via.placeholder.com/100'">
				<div class="student-details">
					<h3>%s</h3>
					<p>Email: %s</p>
					<p>Kelas: %s</p>
					<p>Tanggal Penukaran: %s</p>
					<p>Status: <span class="status status-%s">%s</span></p>
				</div>
			</div>
			
			<div class="item-details">
				<h3>Detail Barang</h3>
				<div class="item-row">
					<span class="item-label">Nama Barang:</span>
					<span>%s</span>
				</div>
				<div class="item-row">
					<span class="item-label">Jumlah:</span>
					<span>%d</span>
				</div>
				<div class="item-row">
					<span class="item-label">Harga (Diamond):</span>
					<span>%d</span>
				</div>
				<div class="item-row">
					<span class="item-label">Total Diamond:</span>
					<span>%d</span>
				</div>
			</div>
			
			<!-- <a href="%s/print-receipt/%s" class="print-button" target="_blank">Cetak Bukti Penukaran</a> -->
			
			<div class="footer">
				<p>Bukti penukaran ini merupakan dokumen resmi. Silakan tunjukkan kode penukaran saat mengambil hadiah.</p>
				<p>&copy; 2025 Cerdikia. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`,
		code,                                // Title
		log.KodePenukaran,                   // Redemption code
		student.ImageProfile,                // Student photo
		student.Nama,                        // Student name
		log.Email,                           // Student email
		className,                           // Class name
		formattedDate,                       // Redemption date
		getStatusClass(log.StatusPenukaran), // Status class for CSS
		getStatusLabel(log.StatusPenukaran), // Status label
		barang.NamaBarang,                   // Item name
		log.Jumlah,                          // Quantity
		barang.Diamond,                      // Price per item
		barang.Diamond*int64(log.Jumlah),    // Total price
		baseurl,                             // baseurl
		log.KodePenukaran,                   // Redemption code for print URL
	)

	// Return HTML content
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// PrintRedemptionReceipt generates and returns a PDF receipt
func PrintRedemptionReceipt(c *gin.Context) {
	db := config.DB
	code := c.Param("code")

	// Get the redemption log
	var log redemption.LogsPenukaranPoint
	if err := db.Where("kode_penukaran = ?", code).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("Redemption with code %s not found", code),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error retrieving redemption data",
				"error":   err.Error(),
			})
		}
		return
	}

	// Get item details
	var barang redemption.BarangInfo
	if err := db.Where("id_barang = ?", log.IDBarang).First(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving item details",
			"error":   err.Error(),
		})
		return
	}

	// Get student details
	type StudentInfo struct {
		Nama         string `gorm:"column:nama"`
		IDKelas      int    `gorm:"column:id_kelas"`
		ImageProfile string `gorm:"column:image_profile"`
	}

	var student StudentInfo
	if err := db.Table("siswa").Where("email = ?", log.Email).First(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving student details",
			"error":   err.Error(),
		})
		return
	}

	// Get class name
	var className string
	var class struct {
		Kelas string `gorm:"column:kelas"`
	}
	if err := db.Table("kelas").Select("kelas").Where("id_kelas = ?", student.IDKelas).First(&class).Error; err == nil {
		className = class.Kelas
	}

	// Format date
	formattedDate := log.TanggalPenukaran.Format("02 January 2006 15:04")

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(190, 10, "BUKTI PENUKARAN HADIAH")
	pdf.Ln(15)

	// Redemption code
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Kode Penukaran: "+log.KodePenukaran)
	pdf.Ln(15)

	// Student information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, "Informasi Siswa")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Nama:")
	pdf.Cell(150, 8, student.Nama)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Email:")
	pdf.Cell(150, 8, log.Email)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Kelas:")
	pdf.Cell(150, 8, className)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Tanggal:")
	pdf.Cell(150, 8, formattedDate)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Status:")
	pdf.Cell(150, 8, getStatusLabel(log.StatusPenukaran))
	pdf.Ln(15)

	// Item information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, "Detail Barang")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Nama Barang:")
	pdf.Cell(150, 8, barang.NamaBarang)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Jumlah:")
	pdf.Cell(150, 8, strconv.Itoa(log.Jumlah))
	pdf.Ln(8)

	pdf.Cell(40, 8, "Harga (Diamond):")
	pdf.Cell(150, 8, strconv.FormatInt(barang.Diamond, 10))
	pdf.Ln(8)

	pdf.Cell(40, 8, "Total Diamond:")
	pdf.Cell(150, 8, strconv.FormatInt(barang.Diamond*int64(log.Jumlah), 10))
	pdf.Ln(20)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(190, 10, "Bukti penukaran ini merupakan dokumen resmi.")
	pdf.Ln(6)
	pdf.Cell(190, 10, "Silakan tunjukkan kode penukaran saat mengambil hadiah.")
	pdf.Ln(6)
	pdf.Cell(190, 10, "© 2025 Cerdikia. All rights reserved.")

	// Create directory for temporary files if it doesn't exist
	tmpDir := filepath.Join("tmp", "receipts")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating directory for PDF",
			"error":   err.Error(),
		})
		return
	}

	// Generate filename
	filename := fmt.Sprintf("receipt_%s_%s.pdf", log.KodePenukaran, time.Now().Format("20060102150405"))
	filePath := filepath.Join(tmpDir, filename)

	// Save PDF to file
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error generating PDF",
			"error":   err.Error(),
		})
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(filePath)

	// Schedule file deletion (after a delay to ensure it's downloaded)
	go func() {
		time.Sleep(5 * time.Minute)
		os.Remove(filePath)
	}()
}

// Helper functions
func getStatusClass(status string) string {
	switch status {
	case "menunggu":
		return "waiting"
	case "selesai":
		return "completed"
	case "dibatalkan":
		return "cancelled"
	default:
		return "waiting"
	}
}

func getStatusLabel(status string) string {
	switch status {
	case "menunggu":
		return "Menunggu"
	case "selesai":
		return "Selesai"
	case "dibatalkan":
		return "Dibatalkan"
	default:
		return "Menunggu"
	}
}
