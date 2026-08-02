package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andersonkmi/bookapi/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookHandler wires HTTP endpoints to the book repository.
type BookHandler struct {
	repo *repository.BookRepository
}

// NewBookHandler creates a BookHandler backed by the given repository.
func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

// RegisterRoutes wires the book endpoints onto the provided router group.
func (h *BookHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/books", h.listBooks)
	routerGroup.GET("/books/:id", h.getBook)
	routerGroup.POST("/books", h.createBook)
	routerGroup.DELETE("/books/:id", h.deleteBook)
}

func (h *BookHandler) listBooks(ctx *gin.Context) {
	books, err := h.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *BookHandler) getBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.repo.Find(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (h *BookHandler) createBook(ctx *gin.Context) {
	var book repository.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Insert(&book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (h *BookHandler) deleteBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
