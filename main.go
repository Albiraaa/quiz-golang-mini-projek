package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"projek/config"
	"projek/handlers"
	"projek/middleware"
)

func main() {
	// init DB & migration
	config.InitDB()

	r := gin.Default()

	// Health check simple
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Golang Book API"})
	})

	// Auth routes (tidak pakai JWT middleware)
	auth := r.Group("/api/users")
	{
		auth.POST("/register", handlers.RegisterUser) // opsional
		auth.POST("/login", handlers.Login)           // sesuai brief
	}

	// Protected routes pakai JWT
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// categories
		protected.GET("/kategori", handlers.GetCategories)
		protected.POST("/kategori", handlers.CreateCategory)
		protected.GET("/kategori/:id", handlers.GetCategoryByID)
		protected.PUT("/kategori/:id", handlers.UpdateCategory)
		protected.DELETE("/kategori/:id", handlers.DeleteCategory)
		protected.GET("/kategori/:id/books", handlers.GetBooksByCategory)

		// books
		protected.GET("/buku", handlers.GetBooks)
		protected.POST("/buku", handlers.CreateBook)
		protected.GET("/buku/:id", handlers.GetBookByID)
		protected.PUT("/buku/:id", handlers.UpdateBook)
		protected.DELETE("/buku/:id", handlers.DeleteBook)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
