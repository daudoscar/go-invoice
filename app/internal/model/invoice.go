package model

type Item struct {
	BaseName     string
	Qty          int
	Subtotal     int
	FinalNominal int
}

type InvoiceData struct {
	OrderID    string
	SellerName string
	Location   string
	MonthYear  string
	FullDate   string
	Items      []Item
}

type FileInfo struct {
	Name string
	Link string
}
