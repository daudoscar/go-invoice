package service

import (
	"context"
	"fmt"
	"go-invoice/pkg/google"
	"strings"

	"google.golang.org/genai"
)

type AIService struct {
	gemini *google.GeminiClient
}

func NewAIService(gemini *google.GeminiClient) *AIService {
	return &AIService{gemini: gemini}
}

func (s *AIService) CleanName(rawNames string, databaseContext string) string {
	ctx := context.Background()

	modelName := "gemini-2.5-flash"

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`Anda adalah pakar pemetaan data inventaris.
Tugas Anda: Menyelaraskan nama produk dan menentukan satuannya.

ATURAN KETAT:
1. FORMAT OUTPUT: Nama Produk|Satuan (Gunakan tanda pipa | sebagai pemisah).
2. DAFTAR SATUAN RESMI: [Rim, Pack, Pcs, Dus, Box, Set, Kg, Galon, Unit, Pasang, Btg, Renceng, Jerigen].
3. Jika satuan tidak ada di daftar resmi atau Anda ragu, WAJIB gunakan: (manual).
4. PRIORITAS DATABASE: Jika mirip dengan DAFTAR DATABASE, gunakan Nama dan Satuan dari database secara eksak.
5. PEMBERSIHAN: Hapus merk toko (Nona.id), promo, dan deskripsi panjang.
6. JUMLAH BARIS OUTPUT HARUS SAMA PERSIS DENGAN INPUT.`, genai.RoleUser),
	}

	fullPrompt := fmt.Sprintf("DAFTAR DATABASE KAMI:\n%s\n\nINPUT PDF:\n%s", databaseContext, rawNames)

	result, err := s.gemini.Client.Models.GenerateContent(ctx, modelName, genai.Text(fullPrompt), config)
	if err != nil {
		return rawNames
	}

	return strings.TrimSpace(result.Text())
}
