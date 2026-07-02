package models

import "gorm.io/gorm"

// Transaction mewakili rekaman transaksi checkout
type Transaction struct {
	gorm.Model

	UserID        uint              `gorm:"not null;index" json:"user_id"`
	InvoiceNumber string            `gorm:"size:100;uniqueIndex;not null" json:"invoice_number"`
	TotalPrice    float64           `gorm:"not null" json:"total_price"`
	Status        string            `gorm:"size:50;default:Selesai;index" json:"status"` // Selesai, Diproses, Batal
	PaymentMethod string            `gorm:"size:100;not null" json:"payment_method"`
	
	// Relasi ke TransactionItems (Preload-able)
	Items []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}

// TransactionItem mewakili detail barang yang dibeli dalam suatu transaksi
type TransactionItem struct {
	gorm.Model

	TransactionID uint    `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	Quantity      int     `gorm:"not null;default:1" json:"quantity"`
	Price         float64 `gorm:"not null" json:"price"` // Menyimpan harga produk saat dibeli

	// Relasi ke Product (Preload-able)
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// DTO untuk Checkout Request
type CheckoutRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}
