package utils

import (
	"encoding/csv"
	"go-invoice/pkg/logger"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ReadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	return reader.ReadAll()
}

func ReadPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var str strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		str.WriteString(content)
	}

	return strings.ReplaceAll(str.String(), "\r", ""), nil
}

func GetPDFFiles(dir string) []string {
	var files []string

	// Gunakan signature yang tepat: (string, fs.DirEntry, error)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		// 1. Cek jika ada error saat mengakses path
		if err != nil {
			return err
		}

		// 2. Jika ini adalah direktori, kita lewati (lanjut ke isinya)
		if d.IsDir() {
			return nil
		}

		// 3. Cek ekstensi file secara case-insensitive
		if strings.ToLower(filepath.Ext(path)) == ".pdf" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		logger.Error("Gagal menyisir direktori invoice", err)
	}

	return files
}
