package main

import (
	"go-invoice/internal/processor"
	"go-invoice/internal/service"
	"go-invoice/internal/ui"
	"go-invoice/pkg/config"
	"go-invoice/pkg/google"
)

func main() {
	cfg := config.LoadConfig()

	sheetsClient := google.NewSheetsClient(cfg.CredentialsPath)
	geminiClient := google.NewGeminiClient(cfg.GeminiAPIKey)
	driveClient := google.NewDriveClient(cfg.CredentialsPath)

	display := ui.NewDisplay()
	menu := ui.NewMenu()

	parser := processor.NewParser()
	aiSvc := service.NewAIService(geminiClient)

	sheetSvc := service.NewSheetService(sheetsClient, cfg.SheetID, aiSvc)
	driveSvc := service.NewDriveService(driveClient)

	app := service.NewAppService(cfg, parser, sheetSvc, driveSvc, display)

	router := ui.NewRouter(menu, app)
	router.Start()
}
