package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

var (
	mu     sync.Mutex
	books  = map[int]Book{}
	nextID = 1
)

// RegisterRoutes wires the book endpoints onto the provided router group.
func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/books", listBooks)
	rg.GET("/books/:id", getBook)
	rg.POST("/books", createBook)
	rg.DELETE("/books/:id", deleteBook)
}

func listBooks(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	out := make([]Book, 0, len(books))
	for _, b := range books {
		out = append(out, b)
	}
	c.JSON(http.StatusOK, out)
}

func getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	b, ok := books[id]
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func createBook(c *gin.Context) {
	var b Book
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	b.ID = nextID
	nextID++
	books[b.ID] = b
	mu.Unlock()

	c.JSON(http.StatusCreated, b)
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	_, ok := books[id]
	delete(books, id)
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
