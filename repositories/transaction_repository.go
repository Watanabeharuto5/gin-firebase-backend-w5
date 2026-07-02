package repositories

import (
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/config"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

// Create menyimpan transaksi baru ke database beserta item-itemnya dalam satu batch
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return config.DB.Create(tx).Error
}

// FindAllByUserID mengambil semua riwayat transaksi milik user beserta detail item & produknya
func (r *TransactionRepository) FindAllByUserID(userID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	// Preload Items dan Preload Product di dalam masing-masing item
	result := config.DB.Preload("Items.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&txs)
	return txs, result.Error
}
