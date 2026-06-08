package google

import (
	"context"
	"go-invoice/pkg/logger" // Sesuaikan dengan nama project baru

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveClient struct {
	Service *drive.Service
}

func NewDriveClient(credentialsPath string) *DriveClient {
	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		logger.Fatal("Gagal inisialisasi Google Drive Service", err)
	}

	return &DriveClient{Service: srv}
}
