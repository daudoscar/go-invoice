package config

import (
	"go-invoice/pkg/logger"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	SheetID         string
	CredentialsPath string
	InputPath       string // Folder PDF (Luar)
	OutputPath      string // Folder Hasil (Luar)
	DriveFolderID   string
	GeminiAPIKey    string
}

func LoadConfig() *Config {
	envPath := filepath.Join("pkg", "config", ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		logger.Warn("[CONFIG] File .env tidak ditemukan, menggunakan env system")
	}

	return &Config{
		SheetID:         os.Getenv("GOOGLE_SHEET_ID"),
		CredentialsPath: os.Getenv("GOOGLE_CREDENTIALS_PATH"),
		InputPath:       os.Getenv("INVOICE_INPUT_PATH"),
		OutputPath:      os.Getenv("INVOICE_OUTPUT_PATH"),
		DriveFolderID:   os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		GeminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
	}
}
