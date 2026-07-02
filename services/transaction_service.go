package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Watanabeharuto5/gin-firebase-backend-w5/config"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/repositories"
	"gorm.io/gorm"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	cartRepo        *repositories.CartRepository
	productRepo     *repositories.ProductRepository
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		transactionRepo: repositories.NewTransactionRepository(),
		cartRepo:        repositories.NewCartRepository(),
		productRepo:     repositories.NewProductRepository(),
	}
}

// CheckoutCart memproses pembelian barang di keranjang belanja, mengurangi stok, membuat rekam transaksi, dan membersihkan keranjang
func (s *TransactionService) CheckoutCart(userID uint, paymentMethod string) (*models.Transaction, error) {
	var transaction *models.Transaction

	// Jalankan transaksi dalam scope DB Transaction agar menjamin ACID compliance
	err := config.DB.Transaction(func(db *gorm.DB) error {
		// 1. Ambil seluruh item keranjang milik user
		var cartItems []models.CartItem
		if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return errors.New("keranjang belanja Anda kosong")
		}

		// 2. Buat Invoice Number unik menggunakan Unix Timestamp
		invoiceNumber := fmt.Sprintf("TRX-%d-%d-KP", time.Now().Unix(), userID)

		var totalPrice float64
		var txItems []models.TransactionItem

		// 3. Proses validasi stok produk & kurangi stoknya
		for _, item := range cartItems {
			var product models.Product
			if err := db.First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("produk '%s' tidak ditemukan", item.Product.Name)
			}

			if !product.IsActive {
				return fmt.Errorf("produk '%s' sedang tidak aktif", product.Name)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("stok produk '%s' tidak mencukupi (tersedia: %d)", product.Name, product.Stock)
			}

			// Kurangi stok
			product.Stock -= item.Quantity
			if err := db.Save(&product).Error; err != nil {
				return err
			}

			// Hitung total harga & tambahkan ke list item transaksi
			totalPrice += product.Price * float64(item.Quantity)
			txItems = append(txItems, models.TransactionItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
		}

		// 4. Buat record Transaksi Utama
		transaction = &models.Transaction{
			UserID:        userID,
			InvoiceNumber: invoiceNumber,
			TotalPrice:    totalPrice,
			Status:        "Pending",
			PaymentMethod: paymentMethod,
			Items:         txItems,
		}

		if err := db.Create(transaction).Error; err != nil {
			return err
		}

		// 5. Kosongkan keranjang belanja
		if err := db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionHistory mengambil seluruh riwayat transaksi belanja user
func (s *TransactionService) GetTransactionHistory(userID uint) ([]models.Transaction, error) {
	return s.transactionRepo.FindAllByUserID(userID)
}

// ConfirmPayment mengubah status transaksi menjadi Selesai setelah pembayaran berhasil diverifikasi
func (s *TransactionService) ConfirmPayment(invoiceNumber string) error {
	return config.DB.Model(&models.Transaction{}).
		Where("invoice_number = ?", invoiceNumber).
		Update("status", "Selesai").Error
}
