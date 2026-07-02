package services

import (
	"errors"

	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/repositories"
	"gorm.io/gorm"
)

type CartService struct {
	cartRepo    *repositories.CartRepository
	productRepo *repositories.ProductRepository
}

func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repositories.NewCartRepository(),
		productRepo: repositories.NewProductRepository(),
	}
}

// GetCart mendapatkan semua item keranjang milik user beserta relasi produknya
func (s *CartService) GetCart(userID uint) ([]models.CartItem, error) {
	return s.cartRepo.FindAllByUserID(userID)
}

// AddToCart menambahkan produk ke keranjang, atau menambah kuantitasnya jika produk sudah ada di keranjang
func (s *CartService) AddToCart(userID, productID uint, quantity int) (*models.CartItem, error) {
	// 1. Verifikasi apakah produk valid dan aktif
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	if !product.IsActive {
		return nil, errors.New("produk sedang tidak aktif")
	}

	// 2. Cek ketersediaan stok
	if product.Stock < quantity {
		return nil, errors.New("stok produk tidak mencukupi")
	}

	// 3. Cek apakah produk sudah ada di keranjang user
	item, err := s.cartRepo.FindByUserAndProduct(userID, productID)
	if err == nil {
		// Produk sudah ada, update kuantitasnya
		newQty := item.Quantity + quantity
		if product.Stock < newQty {
			return nil, errors.New("total kuantitas melebihi stok produk")
		}
		item.Quantity = newQty
		if err := s.cartRepo.Update(item); err != nil {
			return nil, err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Produk belum ada, buat baru
		item = &models.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := s.cartRepo.Create(item); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Load detail produk untuk response
	itemWithProduct, err := s.cartRepo.FindByUserAndProduct(userID, productID)
	if err == nil {
		itemWithProduct.Product = *product
		return itemWithProduct, nil
	}

	return item, nil
}

// UpdateCartItem memperbarui kuantitas produk dalam keranjang belanja
func (s *CartService) UpdateCartItem(userID, itemID uint, quantity int) (*models.CartItem, error) {
	// 1. Dapatkan item keranjang
	item, err := s.cartRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item keranjang tidak ditemukan")
		}
		return nil, err
	}

	// 2. Pastikan kepemilikan user cocok
	if item.UserID != userID {
		return nil, errors.New("akses ditolak: bukan pemilik keranjang")
	}

	// 3. Verifikasi stok produk
	product, err := s.productRepo.FindByID(item.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Stock < quantity {
		return nil, errors.New("stok produk tidak mencukupi")
	}

	// 4. Update kuantitas
	item.Quantity = quantity
	if err := s.cartRepo.Update(item); err != nil {
		return nil, err
	}

	item.Product = *product
	return item, nil
}

// RemoveFromCart menghapus item dari keranjang belanja
func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	// 1. Verifikasi keberadaan dan kepemilikan item
	item, err := s.cartRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item keranjang tidak ditemukan")
		}
		return err
	}

	if item.UserID != userID {
		return errors.New("akses ditolak: bukan pemilik keranjang")
	}

	// 2. Hapus
	return s.cartRepo.Delete(itemID, userID)
}

// ClearCart menghapus seluruh item keranjang milik user
func (s *CartService) ClearCart(userID uint) error {
	return s.cartRepo.Clear(userID)
}
