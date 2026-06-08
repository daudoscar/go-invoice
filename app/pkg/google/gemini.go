package google

import (
	"context"
	"go-invoice/pkg/logger"

	"google.golang.org/genai"
)

type GeminiClient struct {
	Client *genai.Client
}

func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()

	// Menggunakan genai.ClientConfig sesuai SDK terbaru
	cfg := &genai.ClientConfig{
		APIKey: apiKey,
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		logger.Fatal("Gagal inisialisasi Google GenAI Service", err)
	}

	return &GeminiClient{Client: client}
}
