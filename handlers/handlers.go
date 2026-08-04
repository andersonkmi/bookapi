package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/andersonkmi/bookapi/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookStore describes the data-access operations required by the handlers.
type BookStore interface {
	Insert(book *repository.Book) error
	Find(id uint) (*repository.Book, error)
	FindAll() ([]repository.Book, error)
	AddComment(bookID uint, text string) (*repository.Comment, error)
	Delete(id uint) error
}

// BookHandler wires HTTP endpoints to the book repository.
type BookHandler struct {
	repo BookStore
}

// NewBookHandler creates a BookHandler backed by the given repository.
func NewBookHandler(repo BookStore) *BookHandler {
	return &BookHandler{repo: repo}
}

// RegisterRoutes wires the book endpoints onto the provided router group.
func (handler *BookHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/books", handler.listBooks)
	routerGroup.GET("/books/:id", handler.getBook)
	routerGroup.POST("/books", handler.createBook)
	routerGroup.POST("/books/:id/comments", handler.addComment)
	routerGroup.DELETE("/books/:id", handler.deleteBook)
}

func (handler *BookHandler) listBooks(ctx *gin.Context) {
	books, err := handler.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (handler *BookHandler) getBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := handler.repo.Find(uint(id))
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

func (handler *BookHandler) createBook(ctx *gin.Context) {
	var book repository.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if book.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "author is required"})
		return
	}

	if err := handler.repo.Insert(&book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

func (handler *BookHandler) addComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Text string `json:"text" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(body.Text) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}

	comment, err := handler.repo.AddComment(uint(id), body.Text)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (handler *BookHandler) deleteBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := handler.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
