package google

import (
	"context"
	"fmt"
	"go-invoice/pkg/logger"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsClient struct {
	Service *sheets.Service
}

func NewSheetsClient(credentialsPath string) *SheetsClient {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		logger.Fatal("Gagal inisialisasi Google Sheets Service", err)
	}

	return &SheetsClient{Service: srv}
}

func (c *SheetsClient) FetchColumn(sheetID, tabName, col string) ([][]interface{}, error) {
	readRange := fmt.Sprintf("'%s'!%s:%s", tabName, col, col)
	resp, err := c.Service.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func (c *SheetsClient) PushRow(sheetID, tabName string, row []interface{}) error {
	rb := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}
	writeRange := fmt.Sprintf("'%s'!A1", tabName)

	_, err := c.Service.Spreadsheets.Values.Append(sheetID, writeRange, rb).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	return err
}

func (c *SheetsClient) UpdateCell(sheetID, tabName, cellAddress, value string) error {
	target := fmt.Sprintf("'%s'!%s", tabName, cellAddress)
	rb := &sheets.ValueRange{
		Values: [][]interface{}{{value}},
	}

	_, err := c.Service.Spreadsheets.Values.Update(sheetID, target, rb).
		ValueInputOption("USER_ENTERED").
		Do()
	return err
}
