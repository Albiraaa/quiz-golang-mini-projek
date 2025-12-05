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
		protected.GET("/categories", handlers.GetCategories)
		protected.POST("/categories", handlers.CreateCategory)
		protected.GET("/categories/:id", handlers.GetCategoryByID)
		protected.PUT("/categories/:id", handlers.UpdateCategory)
		protected.DELETE("/categories/:id", handlers.DeleteCategory)
		protected.GET("/categories/:id/books", handlers.GetBooksByCategory)

		// books
		protected.GET("/books", handlers.GetBooks)
		protected.POST("/books", handlers.CreateBook)
		protected.GET("/books/:id", handlers.GetBookByID)
		protected.PUT("/books/:id", handlers.UpdateBook)
		protected.DELETE("/books/:id", handlers.DeleteBook)
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
