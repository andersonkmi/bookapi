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
func RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/books", listBooks)
	routerGroup.GET("/books/:id", getBook)
	routerGroup.POST("/books", createBook)
	routerGroup.DELETE("/books/:id", deleteBook)
}

func listBooks(ctx *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	out := make([]Book, 0, len(books))
	for _, book := range books {
		out = append(out, book)
	}
	ctx.JSON(http.StatusOK, out)
}

func getBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	book, ok := books[id]
	mu.Unlock()

	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func createBook(ctx *gin.Context) {
	var book Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	book.ID = nextID
	nextID++
	books[book.ID] = book
	mu.Unlock()

	ctx.JSON(http.StatusCreated, book)
}

func deleteBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	_, ok := books[id]
	delete(books, id)
	mu.Unlock()

	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
