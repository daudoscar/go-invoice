package ui

import (
	"fmt"
	"go-invoice/pkg/response"
	"strings"
)

type Menu struct{}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) ShowHeader() {
	fmt.Println("\n==============================")
	fmt.Println("      INVOICE MANAGER KING DAUD      ")
	fmt.Println("       PT WIN MEDIKA         ")
	fmt.Println("==============================")
}

func (m *Menu) GetChoice() string {
	fmt.Println("\n[ MENU UTAMA ]")
	fmt.Println("1. Scan & Process PDF Lokal")
	fmt.Println("2. Sync Hyperlink (Drive to Sheets)")
	fmt.Println("0. Keluar")
	fmt.Print("\nPilih menu (1/2/0): ")

	var choice string
	fmt.Scanln(&choice)
	return strings.TrimSpace(choice)
}

func (m *Menu) StartProcess(mode string) {
	fmt.Println()
	response.Info(fmt.Sprintf(">>> MEMULAI PROSES: %s", mode))
}

func (m *Menu) EndProcess() {
	fmt.Println()
	response.Info(">>> PROSES SELESAI")
	fmt.Println("\nTekan [ENTER] untuk kembali ke menu utama...")
	fmt.Scanln()
}

func (m *Menu) ShowExitMessage() {
	fmt.Println("\nTerima kasih, King Daud. Sampai jumpa!")
}
