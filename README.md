# Shopee Invoice Manager

Internal automation tool for **PT Win Medika** that turns Shopee invoice PDFs into clean, reconciled bookkeeping records.

It reads invoice PDFs, extracts the line-item data with **Google Gemini**, applies **pro-rata price allocation** so the totals match bank mutations exactly, writes the results to **Google Sheets**, and archives the source files to **Google Drive** — replacing each Order ID with a clickable Drive link.

---

## ✨ Features

- **PDF parsing** — reads Shopee invoice PDFs from a local folder.
- **AI extraction** — uses Gemini to pull structured line-item data out of each invoice.
- **Pro-rata pricing** — proportionally distributes the final paid amount across items, so the Google Sheets total always reconciles with the bank mutation even with discounts or bundled packages.
- **Google Sheets sync** — pushes the processed data straight into your spreadsheet.
- **Google Drive archiving** — uploads processed files and rewrites the Order ID column into clickable Drive hyperlinks.
- **Interactive CLI** — simple menu-driven workflow.

---

## 📂 Project Structure

```
go-invoice/
├── app/                    # Go application
│   ├── cmd/app/main.go     # Entry point
│   ├── internal/
│   │   ├── model/          # Domain models (invoice)
│   │   ├── processor/      # PDF parsing & helpers
│   │   ├── service/        # App, sheet, drive, AI services
│   │   └── ui/             # Menu, router, display
│   ├── pkg/
│   │   ├── config/         # Config loader (.env, bot.json)
│   │   ├── google/         # Sheets, Drive, Gemini clients
│   │   ├── logger/         # Logging
│   │   ├── response/       # Output helpers
│   │   └── utils/          # Reader utilities
│   ├── go.mod
│   └── go.sum
├── invoices/               # Invoice documents (not committed)
│   ├── input/              # Drop Shopee invoice PDFs here
│   └── output/             # Auto-archived results (by month/location)
└── README.md
```

> `invoices/` and all secrets are excluded from Git via `.gitignore`.

---

## 🔧 Requirements

- [Go 1.25+](https://go.dev/dl/)
- A Google Cloud **service account** with access to your Google Sheet & Drive folder (`bot.json`)
- A **Gemini API key**

---

## ⚙️ Configuration

Create `app/pkg/config/.env`:

```env
# --- Google Cloud ---
GOOGLE_CREDENTIALS_PATH="pkg/config/bot.json"
GOOGLE_SHEET_ID="your_spreadsheet_id"
GOOGLE_DRIVE_FOLDER_ID="your_drive_folder_id"

# --- Local Paths ---
INVOICE_INPUT_PATH="../invoices/input"
INVOICE_OUTPUT_PATH="../invoices/output"

# --- AI ---
GEMINI_API_KEY="your_gemini_api_key"
```

Place your service-account credentials at `app/pkg/config/bot.json`.

> ⚠️ **Never commit `.env` or `bot.json`.** They contain secrets and are already git-ignored.

---

## 🚀 Usage

From the `app/` folder:

```bash
# Build
go build -o invoice-manager.exe ./cmd/app

# Run
./invoice-manager.exe
```

Then pick a menu option:

- **`1` — Scan & Process PDF Lokal**
  Reads PDFs from `invoices/input`, computes pro-rata pricing, pushes data to Google Sheets, and moves the files into `invoices/output`.
- **`2` — Sync Hyperlink (Drive → Sheets)**
  Uploads archived files to Google Drive and updates the Order ID column in Sheets to clickable links.
- **`0` — Keluar**
  Exit.

---

## 📌 Notes

- Run the executable from inside `app/` so the relative `INVOICE_INPUT_PATH` / `INVOICE_OUTPUT_PATH` resolve correctly.
- Pro-rata allocation guarantees per-item totals always sum to the actual amount paid.

---

© 2026 PT Win Medika — Internal Enterprise Tools
