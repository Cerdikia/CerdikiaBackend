package controllers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ModuleKey struct {
	Mapel     string
	Kelas     string
	Judul     string
	Deskripsi string
}

// Fungsi untuk mendapatkan rich text dari sel Excel
func getRichTextFromCell(f *excelize.File, sheet, cell string) (string, error) {
	// Coba dapatkan rich text dari sel
	richText, err := f.GetCellRichText(sheet, cell)
	if err != nil {
		// Jika error, mungkin bukan rich text, coba dapatkan nilai biasa
		plainText, err := f.GetCellValue(sheet, cell)
		if err != nil {
			// Jika tidak bisa mendapatkan nilai, kembalikan string kosong tanpa error
			// untuk menghindari sel kosong dianggap error
			return "", nil
		}
		return plainText, nil
	}

	// Jika tidak ada rich text, kembalikan string kosong
	if len(richText) == 0 {
		plainText, _ := f.GetCellValue(sheet, cell)
		return plainText, nil
	}

	// Jika ada rich text, konversi ke HTML
	var htmlText string
	for _, run := range richText {
		text := run.Text

		// Escape HTML characters
		text = strings.ReplaceAll(text, "&", "&amp;")
		text = strings.ReplaceAll(text, "<", "&lt;")
		text = strings.ReplaceAll(text, ">", "&gt;")

		// Cek formatting dari font
		if run.Font != nil {
			// Bold
			if run.Font.Bold {
				text = "<strong>" + text + "</strong>"
			}
			// Italic
			if run.Font.Italic {
				text = "<em>" + text + "</em>"
			}
			// Underline
			if run.Font.Underline != "" {
				text = "<u>" + text + "</u>"
			}
		}

		// Tambahkan text dengan formatting ke hasil
		htmlText += text
	}

	// Jika tidak ada teks yang dihasilkan, coba dapatkan nilai biasa
	if htmlText == "" {
		plainText, _ := f.GetCellValue(sheet, cell)
		return plainText, nil
	}

	return htmlText, nil
}

// Fungsi untuk mendapatkan style alignment dari sel
func getCellAlignment(f *excelize.File, sheet, cell string) (string, error) {
	// Dapatkan style ID dari sel
	styleID, err := f.GetCellStyle(sheet, cell)
	if err != nil {
		return "", err
	}

	// Dapatkan style dari ID
	style, err := f.GetStyle(styleID)
	if err != nil {
		return "", err
	}

	// Cek horizontal alignment
	if style.Alignment != nil {
		switch style.Alignment.Horizontal {
		case "center":
			return "text-align: center;", nil
		case "right":
			return "text-align: right;", nil
		case "justify":
			return "text-align: justify;", nil
		}
	}

	return "", nil
}

// Fungsi untuk memeriksa apakah sel berisi gambar dan/atau teks
func processCellContent(f *excelize.File, sheet, cell string, rowText string) (string, bool, error) {
	baseurl := os.Getenv("BASEURL")
	// Mendapatkan semua gambar dari sheet
	pics, err := f.GetPictures(sheet, cell)
	if err != nil {
		return "", false, err
	}

	// Cek apakah ini adalah sel E7 atau E8 yang mungkin berisi teks italic
	isSpecialCell := cell == "E7" || cell == "E8"

	// Dapatkan rich text dari sel
	formattedText, err := getRichTextFromCell(f, sheet, cell)
	if err != nil {
		// Jika error, gunakan teks biasa dari row
		formattedText = rowText
	}

	// Jika ini adalah sel khusus dan teks masih kosong, coba perlakuan khusus
	if isSpecialCell && formattedText == "" {
		// Untuk sel E7 dan E8 yang berisi teks italic, kita coba dapatkan nilai langsung
		// dan bungkus dalam tag italic
		plainText, _ := f.GetCellValue(sheet, cell)
		if plainText != "" {
			formattedText = "<em>" + plainText + "</em>"
		}
	}

	// Jika tidak ada teks dari rich text atau dari row, coba dapatkan nilai sel langsung
	if formattedText == "" && rowText == "" {
		formattedText, _ = f.GetCellValue(sheet, cell)
	}

	// Dapatkan alignment
	alignment, _ := getCellAlignment(f, sheet, cell)

	// Jika tidak ada gambar, kembalikan teks saja dengan formatting
	if len(pics) == 0 {
		return formattedText, false, nil
	}

	// Buat nama file unik untuk gambar
	timestamp := time.Now().UnixNano()
	// Gunakan extension default .png karena excelize.Picture tidak memiliki field Name
	fileExt := ".png"
	imgFileName := fmt.Sprintf("%d%s", timestamp, fileExt)
	imgPath := filepath.Join("tmp", "images", imgFileName)
	imgURL := baseurl + "/tmp/images/" + url.PathEscape(imgFileName)

	// Simpan gambar ke file
	err = os.WriteFile(imgPath, pics[0].File, 0644)
	if err != nil {
		return "", true, err
	}

	// Jika ada teks di sel yang sama dengan gambar, gabungkan keduanya
	if formattedText != "" {
		// Kembalikan URL gambar dan teks sebagai array
		textStyle := ""
		if alignment != "" {
			textStyle = fmt.Sprintf(" style='%s'", alignment)
		}
		return imgURL + "|TEXT|" + formattedText + "|STYLE|" + textStyle, true, nil
	}

	// Jika hanya ada gambar, kembalikan URL gambar saja
	return imgURL, true, nil
}

// Fungsi untuk membersihkan HTML dari tag kosong
func cleanHTML(html string) string {
	// Jika HTML kosong, kembalikan string kosong
	if html == "" {
		return ""
	}

	// Hapus tag paragraf yang benar-benar kosong (<p></p>)
	emptyParagraphRegex := regexp.MustCompile(`<p[^>]*></p>`)
	html = emptyParagraphRegex.ReplaceAllString(html, "")

	// Hapus tag paragraf yang hanya berisi whitespace
	emptyParagraphWithSpaceRegex := regexp.MustCompile(`<p[^>]*>[\s\n\r]+</p>`)
	html = emptyParagraphWithSpaceRegex.ReplaceAllString(html, "")

	// Hapus whitespace berlebih di awal dan akhir
	html = strings.TrimSpace(html)

	return html
}

// Fungsi untuk membuat HTML content dari teks atau gambar
func createHTMLContent(content string, isImage bool) string {
	if isImage {
		// Cek apakah content berisi kombinasi gambar dan teks
		parts := strings.Split(content, "|TEXT|")
		if len(parts) > 1 {
			// Ada gambar dan teks
			imgURL := parts[0]

			// Pisahkan teks dan style
			textParts := strings.Split(parts[1], "|STYLE|")
			text := textParts[0]
			style := ""
			if len(textParts) > 1 {
				style = textParts[1]
			}

			// Buat HTML dengan gambar dan teks
			return fmt.Sprintf("<p><img src=\"%s\" style=\"width: 50%%; display: block; margin: auto;\" alt=\"Gambar Soal\" loading='lazy'></p><p%s>%s</p>",
				imgURL, style, text)
		}

		// Jika hanya gambar, buat tag img
		return fmt.Sprintf("<p><img src=\"%s\" style=\"width: 50%%; display: block; margin: auto;\" alt=\"Gambar Soal\" loading='lazy'></p>", content)
	} else {
		// Jika konten kosong, kembalikan string kosong
		if content == "" {
			return ""
		}

		// Cek apakah konten sudah berisi tag HTML
		hasHTML := regexp.MustCompile("<[a-z]+[^>]*>").MatchString(content)
		if hasHTML {
			// Jika sudah berisi HTML, kembalikan apa adanya
			return content
		} else {
			// Jika belum berisi HTML, bungkus dengan tag paragraf
			return "<p>" + content + "</p>"
		}
	}
}

// Fungsi untuk membuat HTML content dari opsi dengan label
func createOptionHTMLContent(optionLabel, content string, isImage bool) string {
	// Gunakan format yang sama dengan createHTMLContent untuk konsistensi
	if isImage {
		// Cek apakah content berisi kombinasi gambar dan teks
		parts := strings.Split(content, "|TEXT|")
		if len(parts) > 1 {
			// Ada gambar dan teks
			imgURL := parts[0]

			// Pisahkan teks dan style
			textParts := strings.Split(parts[1], "|STYLE|")
			text := textParts[0]
			style := ""
			if len(textParts) > 1 {
				style = textParts[1]
			}

			// Buat HTML dengan gambar dan teks
			return fmt.Sprintf("<p><img src=\"%s\" style=\"width: 50%%; display: block; margin: auto;\" alt=\"Opsi %s\" loading='lazy'></p><p%s>%s</p>",
				imgURL, optionLabel, style, text)
		}

		// Jika hanya gambar, buat tag img
		return fmt.Sprintf("<p><img src=\"%s\" style=\"width: 50%%; display: block; margin: auto;\" alt=\"Opsi %s\" loading='lazy'></p>", content, optionLabel)
	} else {
		// Jika konten kosong, kembalikan string kosong
		if content == "" {
			return ""
		}

		// Cek apakah konten sudah berisi tag HTML
		hasHTML := regexp.MustCompile("<[a-z]+[^>]*>").MatchString(content)
		if hasHTML {
			// Jika sudah berisi HTML, kembalikan apa adanya
			return content
		} else {
			// Jika belum berisi HTML, bungkus dengan tag paragraf
			return "<p>" + content + "</p>"
		}
	}
}

func ImportSoalHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File tidak ditemukan"})
		return
	}

	os.MkdirAll("./tmp", 0755)
	os.MkdirAll("./tmp/images", 0755)

	tempPath := "./tmp/" + file.Filename
	c.SaveUploadedFile(file, tempPath)
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal membuka file"})
		return
	}

	sheetName := "Sheet1"
	rows, _ := f.GetRows(sheetName)
	if len(rows) < 2 {
		c.JSON(400, gin.H{"error": "Data kosong"})
		return
	}

	type Soal struct {
		Soal    string `json:"soal"`
		Jenis   string `json:"jenis"`
		OpsiA   string `json:"opsi_a"`
		OpsiB   string `json:"opsi_b"`
		OpsiC   string `json:"opsi_c"`
		OpsiD   string `json:"opsi_d"`
		Jawaban string `json:"jawaban"`
	}
	type Module struct {
		JudulModule     string `json:"judul_module"`
		DeskripsiModule string `json:"deskripsi_module"`
		Soal            []Soal `json:"soal"`
	}
	type Grouped struct {
		Mapel  string   `json:"mapel"`
		Kelas  string   `json:"kelas"`
		Module []Module `json:"module"`
	}

	// map[mapel|kelas]map[judul_module|deskripsi_module]*Module
	groupMap := make(map[string]map[string]*Module)
	resultMap := make(map[string]*Grouped)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		for len(row) < 11 {
			row = append(row, "")
		}
		mapel := strings.TrimSpace(row[0])
		kelas := strings.TrimSpace(row[1])
		judul := strings.TrimSpace(row[2])
		deskripsi := strings.TrimSpace(row[3])
		// Proses soal
		soalCell := fmt.Sprintf("E%d", i+1)
		soalContent, soalHasPic, _ := processCellContent(f, sheetName, soalCell, row[4])
		soalStr := ""
		if soalHasPic {
			soalStr = createHTMLContent(soalContent, true)
		} else {
			soalStr = createHTMLContent(soalContent, false)
		}

		// Proses opsi A
		opsiACell := fmt.Sprintf("G%d", i+1)
		opsiAContent, opsiAHasPic, _ := processCellContent(f, sheetName, opsiACell, row[6])
		opsiAStr := ""
		if opsiAHasPic {
			opsiAStr = createOptionHTMLContent("A", opsiAContent, true)
		} else {
			opsiAStr = createOptionHTMLContent("A", opsiAContent, false)
		}

		// Proses opsi B
		opsiBCell := fmt.Sprintf("H%d", i+1)
		opsiBContent, opsiBHasPic, _ := processCellContent(f, sheetName, opsiBCell, row[7])
		opsiBStr := ""
		if opsiBHasPic {
			opsiBStr = createOptionHTMLContent("B", opsiBContent, true)
		} else {
			opsiBStr = createOptionHTMLContent("B", opsiBContent, false)
		}

		// Proses opsi C
		opsiCCell := fmt.Sprintf("I%d", i+1)
		opsiCContent, opsiCHasPic, _ := processCellContent(f, sheetName, opsiCCell, row[8])
		opsiCStr := ""
		if opsiCHasPic {
			opsiCStr = createOptionHTMLContent("C", opsiCContent, true)
		} else {
			opsiCStr = createOptionHTMLContent("C", opsiCContent, false)
		}

		// Proses opsi D
		opsiDCell := fmt.Sprintf("J%d", i+1)
		opsiDContent, opsiDHasPic, _ := processCellContent(f, sheetName, opsiDCell, row[9])
		opsiDStr := ""
		if opsiDHasPic {
			opsiDStr = createOptionHTMLContent("D", opsiDContent, true)
		} else {
			opsiDStr = createOptionHTMLContent("D", opsiDContent, false)
		}

		soal := Soal{
			Soal:    soalStr,
			Jenis:   strings.TrimSpace(row[5]),
			OpsiA:   opsiAStr,
			OpsiB:   opsiBStr,
			OpsiC:   opsiCStr,
			OpsiD:   opsiDStr,
			Jawaban: strings.ToLower(strings.TrimSpace(row[10])),
		}
		groupKey := mapel + "|" + kelas
		moduleKey := judul + "|" + deskripsi
		if resultMap[groupKey] == nil {
			resultMap[groupKey] = &Grouped{
				Mapel:  mapel,
				Kelas:  kelas,
				Module: []Module{},
			}
		}
		if groupMap[groupKey] == nil {
			groupMap[groupKey] = make(map[string]*Module)
		}
		if groupMap[groupKey][moduleKey] == nil {
			m := &Module{
				JudulModule:     judul,
				DeskripsiModule: deskripsi,
				Soal:            []Soal{},
			}
			groupMap[groupKey][moduleKey] = m
			resultMap[groupKey].Module = append(resultMap[groupKey].Module, *m)
		}
		// Tambahkan soal ke module
		groupMap[groupKey][moduleKey].Soal = append(groupMap[groupKey][moduleKey].Soal, soal)
	}

	// Convert resultMap to slice
	var result []Grouped
	for _, v := range resultMap {
		// Update module slice with actual soal (karena pointer)
		for i, mod := range v.Module {
			key := mod.JudulModule + "|" + mod.DeskripsiModule
			if mPtr, ok := groupMap[v.Mapel+"|"+v.Kelas][key]; ok {
				v.Module[i].Soal = mPtr.Soal
			}
		}
		result = append(result, *v)
	}

	c.JSON(200, result)
}
