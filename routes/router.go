package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/handlers"
	"github.com/Watanabeharuto5/gin-firebase-backend-w5/middleware"
)

func SetupRouter() *gin.Engine {
	// gin.Default() sudah include Logger & Recovery middleware
	r := gin.Default()

	// ─── CORS Middleware ───────────────────────────────────
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ─── Init handlers ─────────────────────────────────────
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	cartHandler := handlers.NewCartHandler()

	// ─── API v1 group ──────────────────────────────────────
	v1 := r.Group("/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "gin-firebase-backend",
			})
		})

		// ── Auth routes (public) ─────────────────────────
		auth := v1.Group("/auth")
		{
			auth.POST("/verify-token", authHandler.VerifyToken)
		}

		// ── Protected routes (butuh JWT) ────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Products
			products := protected.Group("/products")
			{
				// GET /v1/products
				products.GET("", productHandler.GetAll)

				// GET /v1/products/:id
				products.GET("/:id", productHandler.GetByID)

				// ── Admin only ─────────────────────────
				adminProducts := products.Group("")
				adminProducts.Use(middleware.AdminOnly())
				{
					// POST /v1/products
					adminProducts.POST("", productHandler.Create)

					// PUT /v1/products/:id
					adminProducts.PUT("/:id", productHandler.Update)

					// DELETE /v1/products/:id
					adminProducts.DELETE("/:id", productHandler.Delete)
				}
			}

			// Cart
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)
				cart.POST("", cartHandler.AddToCart)
				cart.POST("/checkout", cartHandler.Checkout)
				cart.PUT("/:id", cartHandler.UpdateCartItem)
				cart.DELETE("/:id", cartHandler.RemoveFromCart)
			}

			// Transactions
			protected.GET("/transactions", cartHandler.GetHistory)
			protected.POST("/transactions/confirm", cartHandler.ConfirmPayment)
			protected.POST("/transactions/cancel", cartHandler.CancelPayment)
		}
	}

	return r
}