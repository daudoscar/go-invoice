package service

import (
	"bufio"
	"fmt"
	"go-invoice/internal/model"
	"go-invoice/internal/processor"
	"go-invoice/internal/ui"
	"go-invoice/pkg/config"
	"go-invoice/pkg/response"
	"go-invoice/pkg/utils"
	"os"
	"strings"
	"time"
)

type AppService struct {
	Cfg     *config.Config
	Parser  *processor.Parser
	Sheet   *SheetService
	Drive   *DriveService
	Display *ui.Display
}

func NewAppService(cfg *config.Config, p *processor.Parser, s *SheetService, d *DriveService, disp *ui.Display) *AppService {
	return &AppService{
		Cfg:     cfg,
		Parser:  p,
		Sheet:   s,
		Drive:   d,
		Display: disp,
	}

}

func (a *AppService) RunLocalProcess() {
	fmt.Println("[DEBUG] Mencari PDF di:", a.Cfg.InputPath)
	pdfFiles := utils.GetPDFFiles(a.Cfg.InputPath)

	a.Display.ShowScanSummary(len(pdfFiles))

	if len(pdfFiles) == 0 {
		return
	}

	stats := make(map[string]int)

	for _, path := range pdfFiles {
		data, err := a.Parser.Parse(path)
		if err != nil {
			response.Error("PARSER", "Gagal proses file: "+path, err)
			continue
		}

		response.Info("Processing Order: " + data.OrderID)

		// Hitung total harga
		var totalHarga int64 = 0
		for _, item := range data.Items {
			totalHarga += int64(item.Subtotal)
		}

		err = a.Sheet.ProcessItems(data.OrderID, data.SellerName, data.Items, data.MonthYear, data.FullDate, data.Location)
		if err != nil {
			response.Error("SHEETS", "Gagal input ke Sheets: "+data.OrderID, err)
			continue
		}

		err = a.Parser.ArchiveToResult(path, a.Cfg.OutputPath, data.OrderID, data.MonthYear, data.Location)
		if err != nil {
			response.Error("ARCHIVE", "Gagal arsip: "+data.OrderID, err)
		} else {
			response.Success(data.OrderID, "Data terinput dan file diarsipkan", int(totalHarga))
			stats[data.Location]++
		}
	}

	a.Display.ShowProcessResult(stats)
}

func (a *AppService) RunSyncHyperlink() {
	fmt.Print("\nMasukkan Nama Folder di Drive (kosongkan untuk Scan Global): ")

	scanner := bufio.NewScanner(os.Stdin)
	var targetFolder string
	if scanner.Scan() {
		targetFolder = strings.TrimSpace(scanner.Text())
	}

	response.Info("Memulai sinkronisasi hyperlink...")

	var files []model.FileInfo
	var err error

	if targetFolder != "" {
		// Gunakan fitur baru: Pencarian terfokus berdasarkan nama folder
		files, err = a.Drive.GetAllFileLinksByFolderName(targetFolder)
	} else {
		// Fallback ke Scan Global menggunakan DriveFolderID dari .env
		files, err = a.Drive.GetAllFileLinksInFolder(a.Cfg.DriveFolderID)
	}

	if err != nil {
		response.Error("DRIVE", "Gagal mengambil data file dari Drive", err)
		return
	}

	// Fetch data kolom I (Order ID) dari Google Sheets untuk pencocokan memori
	allSheetValues, err := a.Sheet.Client.FetchColumn(a.Sheet.SheetID, TargetSheet, "I")
	if err != nil {
		response.Error("SHEETS", "Gagal fetch kolom I", err)
		return
	}

	successCount := 0
	for _, file := range files {
		orderID := strings.TrimSuffix(file.Name, ".pdf")
		rowNumbers := a.Sheet.FindRowsInMemory(allSheetValues, orderID)

		if len(rowNumbers) == 0 {
			continue
		}

		fmt.Printf("[PROCESSING] ID: %s (%d baris ditemukan)\n", orderID, len(rowNumbers))

		for _, rowNum := range rowNumbers {
			// Update kolom I dengan HYPERLINK, teks label menggunakan orderID (Program Lama)
			err = a.Sheet.UpdateHyperlink(rowNum, "I", file.Link, orderID)

			if err != nil {
				response.Error("SHEETS", fmt.Sprintf("Gagal update baris %d", rowNum), err)
			} else {
				successCount++
			}

			// Delay untuk menghindari limit rate Google Sheets API
			time.Sleep(1200 * time.Millisecond)
		}
	}

	a.Display.ShowSyncResult(successCount, len(files))
}
