package service

import (
	"fmt"
	"go-invoice/internal/model"
	"go-invoice/pkg/google"
	"strings"
	"time"
)

type DriveService struct {
	Client *google.DriveClient
}

func NewDriveService(client *google.DriveClient) *DriveService {
	return &DriveService{Client: client}
}

// GetFolderIDByName mencari ID folder berdasarkan namanya di Google Drive
func (s *DriveService) GetFolderIDByName(folderName string) (string, error) {
	// Memberikan jeda agar API tidak throttling
	time.Sleep(300 * time.Millisecond)

	query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and trashed = false", folderName)

	res, err := s.Client.Service.Files.List().
		Q(query).
		Fields("files(id, name)").
		PageSize(1).
		Do()

	if err != nil {
		return "", fmt.Errorf("error mencari folder: %w", err)
	}

	if len(res.Files) == 0 {
		return "", fmt.Errorf("folder dengan nama '%s' tidak ditemukan di Drive", folderName)
	}

	return res.Files[0].Id, nil
}

// GetAllFileLinksByFolderName adalah entry point utama jika ingin mencari berdasarkan nama folder
func (s *DriveService) GetAllFileLinksByFolderName(targetFolderName string) ([]model.FileInfo, error) {
	// 1. Cari ID foldernya terlebih dahulu
	folderID, err := s.GetFolderIDByName(targetFolderName)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DRIVE] Fokus pencarian di folder: %s (ID: %s)\n", targetFolderName, folderID)

	// 2. Jalankan pencarian rekursif menggunakan ID yang ditemukan
	return s.GetAllFileLinksInFolder(folderID)
}

// GetAllFileLinksInFolder mengumpulkan semua link PDF dalam satu folder (rekursif)
func (s *DriveService) GetAllFileLinksInFolder(folderID string) ([]model.FileInfo, error) {
	var allFiles []model.FileInfo

	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	call := s.Client.Service.Files.List().
		Q(query).
		Fields("nextPageToken, files(id, name, webViewLink, mimeType)").
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true)

	for {
		// Strategi Free Tier: Delay untuk menghindari 429 Too Many Requests
		time.Sleep(600 * time.Millisecond)

		res, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("error listing files in folder %s: %w", folderID, err)
		}

		for _, f := range res.Files {
			switch f.MimeType {
			case "application/vnd.google-apps.folder":
				// Rekursi ke sub-folder
				time.Sleep(200 * time.Millisecond)
				subFiles, err := s.GetAllFileLinksInFolder(f.Id)
				if err != nil {
					fmt.Printf("Warning: failed to scan sub-folder %s: %v\n", f.Name, err)
					continue
				}
				allFiles = append(allFiles, subFiles...)

			default:
				// Filter hanya file PDF
				if strings.HasSuffix(strings.ToLower(f.Name), ".pdf") {
					allFiles = append(allFiles, model.FileInfo{
						Name: f.Name,
						Link: f.WebViewLink,
					})
				}
			}
		}

		if res.NextPageToken == "" {
			break
		}
		call.PageToken(res.NextPageToken)
	}

	return allFiles, nil
}

// GetFileLink mencari link webView untuk file spesifik (berdasarkan nama file tunggal)
func (s *DriveService) GetFileLink(fileName string) (string, error) {
	time.Sleep(500 * time.Millisecond)

	safeFileName := strings.ReplaceAll(fileName, "'", "\\'")
	query := fmt.Sprintf("name contains '%s' and mimeType = 'application/pdf' and trashed = false", safeFileName)

	res, err := s.Client.Service.Files.List().
		Q(query).
		Fields("files(id, name, webViewLink)").
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		PageSize(1).
		Do()

	if err != nil {
		return "", fmt.Errorf("error listing files: %w", err)
	}

	if len(res.Files) > 0 {
		return res.Files[0].WebViewLink, nil
	}

	return "", nil
}
