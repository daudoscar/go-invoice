# Shopee Invoice Manager (v1) - PT Win Medika

Aplikasi berbasis Go untuk otomatisasi pembacaan invoice Shopee (PDF), ekstraksi data menggunakan AI (Gemini), dan sinkronisasi ke Google Sheets serta pengarsipan ke Google Drive.

## 📂 Struktur Direktori
.
├── Dokumen_Invoice/
│   ├── input/          # Letakkan file PDF invoice di sini
│   └── output/         # Hasil arsip otomatis (Terorganisir per Bulan/Lokasi)
└── go-invoice/
    ├── invoice-reader.exe  # File eksekusi utama
    ├── .env                # Konfigurasi API & Path
    └── bot.json            # Google Service Account Credentials

## 🚀 Cara Penggunaan
1.  Pastikan folder `input` berisi file PDF invoice Shopee.
2.  Buka terminal atau Command Prompt di dalam folder `go-invoice`.
3.  Jalankan aplikasi:
    ```bash
    ./run.exe
    ```
4.  Pilih Menu pada layar:
    - **[1] Scan & Process PDF Lokal**: Membaca PDF, hitung harga (pro-rata), kirim ke GSheets, dan pindahkan file ke folder `output`.
    - **[2] Sync Drive Hyperlink**: Mengunggah file dari `output` ke Google Drive dan memperbarui kolom Order ID di GSheets menjadi link yang bisa diklik.

## 🛠️ Konfigurasi (.env)
Pastikan file `.env` di dalam folder `go-invoice` memiliki pengaturan yang benar:
- `INVOICE_INPUT_PATH="../Dokumen_Invoice/input"`
- `INVOICE_OUTPUT_PATH="../Dokumen_Invoice/output"`
- `GOOGLE_CREDENTIALS_PATH="bot.json"`
- `GOOGLE_SHEET_ID="ID_SPREADSHEET_ANDA"`
- `GEMINI_API_KEY="API_KEY_GEMINI_ANDA"`

## 📌 Catatan Teknis
- **Logika Pro-rata**: Aplikasi secara otomatis membagi total pembayaran akhir ke setiap item secara proporsional. Hal ini menjamin total di Google Sheets akan tepat sama dengan mutasi bank meskipun ada diskon atau bundling paket.
- **Keamanan**: Jangan membagikan file `bot.json` atau `GEMINI_API_KEY` kepada pihak luar.

---
© 2026 PT Win Medika - Internal Enterprise Tools