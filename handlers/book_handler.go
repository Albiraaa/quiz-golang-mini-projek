package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"projek/config"
	"projek/models"
)

type BookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	TotalPage   int    `json:"total_page" binding:"required"`
	CategoryID  int    `json:"category_id" binding:"required"`
}

// logic thickness
func calcThickness(totalPage int) string {
	if totalPage > 100 {
		return "tebal"
	}
	return "tipis" // total_page <= 100
}

// GET /api/books
func GetBooks(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url, release_year,
		       price, total_page, thickness, category_id,
		       created_at, created_by, modified_at, modified_by
		FROM books
		ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query books"})
		return
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.ImageURL,
			&b.ReleaseYear,
			&b.Price,
			&b.TotalPage,
			&b.Thickness,
			&b.CategoryID,
			&b.CreatedAt,
			&b.CreatedBy,
			&b.ModifiedAt,
			&b.ModifiedBy,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan book"})
			return
		}
		books = append(books, b)
	}

	c.JSON(http.StatusOK, books)
}

// POST /api/books
func CreateBook(c *gin.Context) {
	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// validasi release_year 1980–2024
	if req.ReleaseYear < 1980 || req.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "release_year must be between 1980 and 2024",
		})
		return
	}

	thickness := calcThickness(req.TotalPage)
	username, _ := c.Get("username")
	now := time.Now()

	err := config.DB.QueryRow(
		`INSERT INTO books (
			title, description, image_url, release_year,
			price, total_page, thickness, category_id,
			created_at, created_by, modified_at, modified_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id`,
		req.Title,
		req.Description,
		req.ImageURL,
		req.ReleaseYear,
		req.Price,
		req.TotalPage,
		thickness,
		req.CategoryID,
		now,
		username,
		now,
		username,
	).Scan(new(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book created"})
}

// GET /api/books/:id
func GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var b models.Book
	err := config.DB.QueryRow(`
		SELECT id, title, description, image_url, release_year,
		       price, total_page, thickness, category_id,
		       created_at, created_by, modified_at, modified_by
		FROM books WHERE id = $1`, id).
		Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.ImageURL,
			&b.ReleaseYear,
			&b.Price,
			&b.TotalPage,
			&b.Thickness,
			&b.CategoryID,
			&b.CreatedAt,
			&b.CreatedBy,
			&b.ModifiedAt,
			&b.ModifiedBy,
		)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query book"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// PUT /api/books/:id (validasi update)
func UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.ReleaseYear < 1980 || req.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "release_year must be between 1980 and 2024",
		})
		return
	}

	// cek book ada
	var exists bool
	err := config.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)`, id).
		Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check book"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found for update"})
		return
	}

	thickness := calcThickness(req.TotalPage)
	username, _ := c.Get("username")
	now := time.Now()

	_, err = config.DB.Exec(`
		UPDATE books
		SET title = $1,
		    description = $2,
		    image_url = $3,
		    release_year = $4,
		    price = $5,
		    total_page = $6,
		    thickness = $7,
		    category_id = $8,
		    modified_at = $9,
		    modified_by = $10
		WHERE id = $11`,
		req.Title,
		req.Description,
		req.ImageURL,
		req.ReleaseYear,
		req.Price,
		req.TotalPage,
		thickness,
		req.CategoryID,
		now,
		username,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}

// DELETE /api/books/:id
func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res, err := config.DB.Exec(`DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found for delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}
