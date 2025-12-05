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

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// GET /api/categories
func GetCategories(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, created_at, created_by, modified_at, modified_by
		FROM categories
		ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query categories"})
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.CreatedAt,
			&cat.CreatedBy,
			&cat.ModifiedAt,
			&cat.ModifiedBy,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan category"})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// POST /api/categories
func CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	username, _ := c.Get("username")
	now := time.Now()

	err := config.DB.QueryRow(
		`INSERT INTO categories (name, created_at, created_by, modified_at, modified_by)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		req.Name,
		now, username,
		now, username,
	).Scan(new(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "category created"})
}

// GET /api/categories/:id
func GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var cat models.Category
	err := config.DB.QueryRow(
		`SELECT id, name, created_at, created_by, modified_at, modified_by
		 FROM categories WHERE id = $1`,
		id,
	).Scan(
		&cat.ID,
		&cat.Name,
		&cat.CreatedAt,
		&cat.CreatedBy,
		&cat.ModifiedAt,
		&cat.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query category"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// PUT /api/categories/:id (biar bisa validasi update)
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// cek ada atau nggak
	var exists bool
	err := config.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`, id).
		Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check category"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found for update"})
		return
	}

	username, _ := c.Get("username")
	now := time.Now()

	_, err = config.DB.Exec(
		`UPDATE categories
		 SET name = $1, modified_at = $2, modified_by = $3
		 WHERE id = $4`,
		req.Name, now, username, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated"})
}

// DELETE /api/categories/:id
func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res, err := config.DB.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found for delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

// GET /api/categories/:id/books
func GetBooksByCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url, release_year, price,
		       total_page, thickness, category_id,
		       created_at, created_by, modified_at, modified_by
		FROM books
		WHERE category_id = $1
		ORDER BY id`, id)
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
