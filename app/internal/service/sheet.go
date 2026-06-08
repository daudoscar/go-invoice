package service

import (
	"errors"
	"fmt"
	"go-invoice/internal/model"
	"go-invoice/pkg/google"
	"strings"
	"time"
)

const TargetSheet = "Log Transaksi"

type SheetService struct {
	Client  *google.SheetsClient
	SheetID string
	AISvc   interface {
		CleanName(string, string) string
	}
}

func NewSheetService(client *google.SheetsClient, sheetID string, aiSvc interface{ CleanName(string, string) string }) *SheetService {
	return &SheetService{
		Client:  client,
		SheetID: sheetID,
		AISvc:   aiSvc,
	}
}

var validUnits = map[string]bool{
	"Rim": true, "Pack": true, "Pcs": true, "Dus": true,
	"Box": true, "Set": true, "Kg": true, "Galon": true,
	"Unit": true, "Pasang": true, "Btg": true, "Renceng": true,
	"Jerigen": true, "(manual)": true,
}

func (s *SheetService) ProcessItems(
	orderID string,
	seller string,
	items []model.Item,
	bulan, tgl, lokasi string,
) error {
	if len(items) == 0 {
		return errors.New("rincian produk kosong")
	}

	var rawNamesBatch []string
	for _, item := range items {
		cleanBaseName := strings.ReplaceAll(item.BaseName, "\n", " ")
		rawNamesBatch = append(rawNamesBatch, strings.TrimSpace(cleanBaseName))
	}
	allNamesRaw := strings.Join(rawNamesBatch, "\n")

	allOutputRaw := s.AISvc.CleanName(allNamesRaw, "")

	var aiResults []string
	for _, line := range strings.Split(allOutputRaw, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			aiResults = append(aiResults, trimmed)
		}
	}

	useAI := len(aiResults) == len(items)

	for i, item := range items {
		finalName := item.BaseName
		finalSatuan := "(manual)"

		if useAI {
			parts := strings.Split(aiResults[i], "|")
			if len(parts) >= 1 {
				finalName = strings.TrimSpace(parts[0])
			}
			if len(parts) == 2 {
				unitPart := strings.TrimSpace(parts[1])
				if validUnits[unitPart] {
					finalSatuan = unitPart
				}
			}
		}

		row := s.buildRow(bulan, tgl, lokasi, finalName, item, finalSatuan, orderID)

		if err := s.Client.PushRow(s.SheetID, TargetSheet, row); err != nil {
			return fmt.Errorf("gagal push ke sheets untuk item %s: %w", item.BaseName, err)
		}

		if i < len(items)-1 {
			time.Sleep(1200 * time.Millisecond)
		}
	}

	return nil
}

func (s *SheetService) buildRow(bulan, tgl, lokasi, produk string, item model.Item, satuan, orderID string) []interface{} {
	return []interface{}{
		bulan, tgl, lokasi, produk, item.Qty, satuan, item.FinalNominal, "", orderID,
	}
}

func (s *SheetService) FindRowsInMemory(allValues [][]interface{}, invoiceID string) []int {
	searchID := strings.TrimSpace(invoiceID)
	var foundRows []int

	for i, row := range allValues {
		if len(row) > 0 {
			cellValue := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
			if cellValue == searchID {
				foundRows = append(foundRows, i+1)
			}
		}
	}
	return foundRows
}

func (s *SheetService) UpdateHyperlink(row int, col string, link string, label string) error {
	cellAddress := fmt.Sprintf("%s%d", col, row)
	formula := fmt.Sprintf(`=HYPERLINK("%s"; "%s")`, link, label)

	return s.Client.UpdateCell(s.SheetID, TargetSheet, cellAddress, formula)
}
