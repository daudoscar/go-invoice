package ui

import (
	"fmt"
	"go-invoice/pkg/response"
	"strings"
)

type Display struct{}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) ShowScanSummary(totalFiles int) {
	fmt.Println()
	response.Info(fmt.Sprintf("Terdeteksi %d file PDF di folder input.", totalFiles))
	fmt.Println(strings.Repeat("-", 35))
}

func (d *Display) ShowProcessResult(stats map[string]int) {
	fmt.Println("\n" + strings.Repeat("=", 35))
	fmt.Println("       RINGKASAN PROSES INVOICE      ")
	fmt.Println(strings.Repeat("-", 35))

	total := 0
	for cabang, jumlah := range stats {
		fmt.Printf(" %-15s : %d invoice\n", strings.ToUpper(cabang), jumlah)
		total += jumlah
	}

	fmt.Println(strings.Repeat("-", 35))
	fmt.Printf(" TOTAL BERHASIL    : %d invoice\n", total)
	fmt.Println(strings.Repeat("=", 35))
}

func (d *Display) ShowSyncResult(success int, totalFiles int) {
	fmt.Println("\n" + strings.Repeat("=", 35))
	fmt.Println("      RINGKASAN SINKRONISASI      ")
	fmt.Println(strings.Repeat("-", 35))
	fmt.Printf(" Berhasil Update   : %d baris\n", success)
	fmt.Printf(" Total Link PDF    : %d file\n", totalFiles)
	fmt.Println(strings.Repeat("=", 35))
}
