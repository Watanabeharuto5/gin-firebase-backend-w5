package models

import "gorm.io/gorm"

// CartItem mewakili produk yang dimasukkan ke dalam keranjang belanja oleh pengguna
type CartItem struct {
	gorm.Model

	UserID    uint    `gorm:"not null;index" json:"user_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`

	// Relasi ke Product (Preload-able)
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Request DTOs untuk Cart

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}
