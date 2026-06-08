package processor

import (
	"fmt"
	"go-invoice/internal/model"
	"go-invoice/pkg/utils"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(path string) (*model.InvoiceData, error) {
	content, err := utils.ReadPDF(path)
	if err != nil {
		return nil, err
	}

	// fmt.Println("\n================ DEBUG PDF CONTENT START ================")
	// fmt.Println(content)
	// fmt.Print("================= DEBUG PDF CONTENT END =================\n\n")

	bulanIndo, tglLengkap := ParseDate(content)
	lokasi := GetCleanLocation(content)

	return &model.InvoiceData{
		OrderID:    p.getOrderID(content),
		SellerName: p.getSellerName(content),
		Location:   lokasi,
		MonthYear:  bulanIndo,
		FullDate:   tglLengkap,
		Items:      p.extractItems(content),
	}, nil
}

// --- CORE REGEX LOGIC (PRIVATE) ---

func (p *Parser) getOrderID(content string) string {
	reOrder := regexp.MustCompile(`\d{6}[A-Z0-9]{8}`)
	match := reOrder.FindString(content)
	if match != "" {
		return match
	}
	return "UNKNOWN"
}

func (p *Parser) getSellerName(content string) string {
	re := regexp.MustCompile(`Nama Penjual:\s*(.*?)No\. Pesanan`)
	m := re.FindStringSubmatch(content)
	if len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// extractItems: Logic utama dengan sistem fallback dan rekonsiliasi pembayaran akhir
func (p *Parser) extractItems(text string) []model.Item {
	// 1. Coba metode original (Linear Regex)
	items := p.extractItemsOriginal(text)

	// 2. Jika gagal (tidak ada item tertangkap), gunakan metode Special Case (Structural/Prorate)
	if len(items) == 0 {
		items = p.extractItemsSpecialCase(text)
	}

	// 3. Logic Rekonsiliasi (Akurasi 100% terhadap mutasi bank/QRIS)
	if len(items) > 0 {
		totalPembayaran := p.ExtractMoney(text, `Total Pembayaran\s*Rp([\d\.]+)`)
		currentTotalBruto := 0
		for _, item := range items {
			currentTotalBruto += item.Subtotal
		}

		if currentTotalBruto > 0 {
			runningTotal := 0
			for i := range items {
				// Hitung rasio kontribusi item terhadap total bruto
				ratio := float64(items[i].Subtotal) / float64(currentTotalBruto)

				if i == len(items)-1 {
					// Penyesuaian terakhir (rounding) agar pas dengan Total Pembayaran di PDF
					items[i].FinalNominal = totalPembayaran - runningTotal
				} else {
					// Distribusikan total pembayaran berdasarkan rasio bruto
					val := math.Round(float64(totalPembayaran) * ratio)
					items[i].FinalNominal = int(val)
					runningTotal += items[i].FinalNominal
				}
			}
		}
	}

	return items
}

func (p *Parser) extractItemsOriginal(text string) []model.Item {
	start := strings.Index(text, "Rincian Pesanan")
	if start == -1 {
		return nil
	}
	section := text[start:]
	if end := strings.Index(section, "Nota Pesanan"); end != -1 {
		section = section[:end]
	}

	// Pattern standar Shopee: No, Produk/Variasi, Harga, Qty, Subtotal
	reItem := regexp.MustCompile(`(\d+)\s+([\s\S]*?)Rp([\d\.]+)\s+(\d+)\s+Rp([\d\.]+)`)
	matches := reItem.FindAllStringSubmatch(section, -1)

	var items []model.Item
	for _, m := range matches {
		name := regexp.MustCompile(`\s+`).ReplaceAllString(m[2], " ")
		name = strings.TrimSpace(name)
		name = strings.TrimPrefix(name, "Produk Variasi Harga Produk Kuantitas Subtotal ")

		qty, _ := strconv.Atoi(m[4])
		subItem := p.CleanAmount(m[5])

		if subItem > 0 {
			items = append(items, model.Item{
				BaseName: name,
				Qty:      qty,
				Subtotal: subItem,
			})
		}
	}
	return items
}

func (p *Parser) extractItemsSpecialCase(text string) []model.Item {
	start := strings.Index(text, "Rincian Pesanan")
	if start == -1 {
		return nil
	}
	section := text[start:]
	if end := strings.Index(section, "Nota Pesanan"); end != -1 {
		section = section[:end]
	}

	reRow := regexp.MustCompile(`(?m)^(\d+)\s*$`)
	indices := reRow.FindAllStringIndex(section, -1)

	var items []model.Item
	for i := 0; i < len(indices); i++ {
		currStart := indices[i][1]
		var currEnd int
		if i < len(indices)-1 {
			currEnd = indices[i+1][0]
		} else {
			currEnd = len(section)
		}

		itemBlock := section[currStart:currEnd]

		// --- 1. DETEKSI QTY (ULTRA AGGRESSIVE SCAN) ---
		qty := 1
		lines := strings.Split(strings.TrimSpace(itemBlock), "\n")

		// Scan dari bawah ke atas untuk mencari angka kuantitas
		for j := len(lines) - 1; j >= 0; j-- {
			line := strings.TrimSpace(lines[j])

			// Abaikan baris yang jelas-jelas nominal uang atau subtotal besar
			if strings.Contains(line, "Rp") || strings.Contains(line, ".") {
				continue
			}

			// Bersihkan baris dari tanda hubung dan spasi agar sisa angka saja
			// Contoh: "- 5 -" menjadi "5"
			cleanLine := strings.NewReplacer("-", "", " ", "").Replace(line)

			if val, err := strconv.Atoi(cleanLine); err == nil {
				// Validasi: Qty retail Shopee biasanya 1-500
				if val > 0 && val < 500 {
					qty = val
					break
				}
			}
		}

		// --- 2. PEMBERSIHAN NAMA (WHITELIST) ---
		rawName := itemBlock
		noise := []string{"Produk", "Variasi", "Harga Produk", "Kuantitas", "Subtotal", "多件多折", "-", "\r", "\n", "PAKET DISKON", "Rp"}

		for _, n := range noise {
			rawName = strings.ReplaceAll(rawName, n, " ")
		}

		// Hapus angka nominal uang (dengan titik ribuan)
		rawName = regexp.MustCompile(`[\d\.]+`).ReplaceAllString(rawName, " ")

		// Hanya sisakan karakter Latin, Angka, dan simbol dasar untuk database [cite: 32, 53]
		reClean := regexp.MustCompile(`[^a-zA-Z0-9\s\-\.\/\(\)]+`)
		name := reClean.ReplaceAllString(rawName, "")
		name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
		name = strings.TrimSpace(name)

		if name != "" {
			items = append(items, model.Item{
				BaseName: name,
				Qty:      qty,
				Subtotal: 0,
			})
		}
	}

	// --- 3. PRO-RATA HARGA (SINKRONISASI TOTAL) ---
	totalBruto := p.ExtractMoney(text, `Subtotal Pesanan\s*Rp([\d\.]+)`)
	if totalBruto == 0 {
		totalBruto = p.ExtractMoney(text, `Subtotal\s*Rp([\d\.]+)`)
	}

	if len(items) > 0 && totalBruto > 0 {
		var totalQty int
		for _, it := range items {
			totalQty += it.Qty
		}

		if totalQty > 0 {
			// Membagi rata total bayar ke semua item berdasarkan jumlah [cite: 56, 63]
			avgPrice := float64(totalBruto) / float64(totalQty)
			for i := range items {
				items[i].Subtotal = int(math.Round(avgPrice * float64(items[i].Qty)))
			}
		}
	}

	return items
}

// --- ARCHIVE LOGIC ---

func (p *Parser) ArchiveToResult(oldPath, baseOutputPath, orderID, bulanPanjang, lokasi string) error {
	parts := strings.Split(bulanPanjang, " ")
	namaBulanCaps := "UNKNOWN"
	tahunDuaDigit := "00"

	if len(parts) >= 2 {
		namaBulanCaps = strings.ToUpper(parts[0])
		tahunPenuh := parts[1]
		if len(tahunPenuh) >= 4 {
			tahunDuaDigit = tahunPenuh[2:]
		}
	}

	folderBulanTahun := fmt.Sprintf("%s %s", namaBulanCaps, tahunDuaDigit)
	folderBulanLokasi := fmt.Sprintf("%s %s", namaBulanCaps, lokasi)
	targetFolder := filepath.Join(baseOutputPath, folderBulanTahun, folderBulanLokasi)

	err := os.MkdirAll(targetFolder, 0755)
	if err != nil {
		return fmt.Errorf("gagal membuat folder arsip: %w", err)
	}

	newPath := filepath.Join(targetFolder, orderID+".pdf")

	if oldPath != newPath {
		if _, err := os.Stat(newPath); err == nil {
			return os.Remove(oldPath)
		} else {
			return os.Rename(oldPath, newPath)
		}
	}
	return nil
}

// --- HELPERS ---

func (p *Parser) CleanAmount(input string) int {
	cleanStr := strings.NewReplacer("Rp", "", ".", "", " ", "", ",", "").Replace(input)
	val, _ := strconv.Atoi(cleanStr)
	return val
}

func (p *Parser) ExtractMoney(text string, pattern string) int {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return p.CleanAmount(match[1])
	}
	return 0
}
