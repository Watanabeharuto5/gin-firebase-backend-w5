package repositories

import (
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/config"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

// FindAllByUserID mengambil semua item keranjang milik user beserta detail produknya
func (r *CartRepository) FindAllByUserID(userID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	result := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&items)
	return items, result.Error
}

// FindByUserAndProduct mencari item keranjang spesifik berdasarkan user_id dan product_id
func (r *CartRepository) FindByUserAndProduct(userID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	result := config.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// FindByID mengambil item keranjang berdasarkan ID miliknya
func (r *CartRepository) FindByID(id uint) (*models.CartItem, error) {
	var item models.CartItem
	result := config.DB.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// Create menyimpan item keranjang baru
func (r *CartRepository) Create(item *models.CartItem) error {
	return config.DB.Create(item).Error
}

// Update menyimpan perubahan kuantitas item keranjang
func (r *CartRepository) Update(item *models.CartItem) error {
	return config.DB.Save(item).Error
}

// Delete menghapus item dari keranjang berdasarkan ID dan kepemilikan user
func (r *CartRepository) Delete(id, userID uint) error {
	return config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{}).Error
}

// Clear menghapus seluruh isi keranjang belanja user (misal setelah checkout)
func (r *CartRepository) Clear(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}
