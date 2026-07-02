package handlers

import (
	"net/http"
	"strconv"

	"github.com/Watanabeharuto5/gin-firebase-backend-w5/models"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/services"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService        *services.CartService
	transactionService *services.TransactionService
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService:        services.NewCartService(),
		transactionService: services.NewTransactionService(),
	}
}

// helper untuk mengambil user_id dari Gin Context secara aman (di-cast dari float64)
func getContextUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return uint(v), true
	case uint:
		return v, true
	case int:
		return uint(v), true
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return uint(id), true
		}
	}
	return 0, false
}

// GetCart godoc
// GET /cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	items, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data keranjang",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang berhasil diambil",
		"data":    items,
	})
}

// AddToCart godoc
// POST /cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid: product_id dan quantity (min 1) wajib diisi",
		})
		return
	}

	item, err := h.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk berhasil ditambahkan ke keranjang",
		"data":    item,
	})
}

// UpdateCartItem godoc
// PUT /cart/:id
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	idParam := c.Param("id")
	itemID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item keranjang tidak valid",
		})
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Kuantitas wajib diisi dengan angka lebih besar dari 0",
		})
		return
	}

	item, err := h.cartService.UpdateCartItem(userID, uint(itemID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kuantitas keranjang berhasil diperbarui",
		"data":    item,
	})
}

// RemoveFromCart godoc
// DELETE /cart/:id
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	idParam := c.Param("id")
	itemID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID item keranjang tidak valid",
		})
		return
	}

	err = h.cartService.RemoveFromCart(userID, uint(itemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk berhasil dihapus dari keranjang",
	})
}

// Checkout godoc
// POST /cart/checkout
func (h *CartHandler) Checkout(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	// Parse request body
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "payment_method wajib diisi",
		})
		return
	}

	// Proses checkout (ACID compliant database transaction)
	transaction, err := h.transactionService.CheckoutCart(userID, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Checkout berhasil diproses",
		"data":    transaction,
	})
}

// GetHistory godoc
// GET /transactions
func (h *CartHandler) GetHistory(c *gin.Context) {
	userID, ok := getContextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Autentikasi gagal: user_id tidak ditemukan",
		})
		return
	}

	txs, err := h.transactionService.GetTransactionHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat transaksi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Riwayat transaksi berhasil diambil",
		"data":    txs,
	})
}

// ConfirmPayment godoc
// POST /transactions/confirm
func (h *CartHandler) ConfirmPayment(c *gin.Context) {
	var req struct {
		InvoiceNumber string `json:"invoice_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invoice_number wajib diisi",
		})
		return
	}

	err := h.transactionService.ConfirmPayment(req.InvoiceNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengonfirmasi pembayaran",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pembayaran berhasil dikonfirmasi",
	})
}
