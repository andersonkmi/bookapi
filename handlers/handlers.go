package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/andersonkmi/bookapi/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookStore describes the data-access operations required by the handlers.
//
// Implementations are expected to report a missing book with
// gorm.ErrRecordNotFound so the handlers can answer 404 instead of 500.
// [github.com/andersonkmi/bookapi/repository.BookRepository] satisfies this
// interface.
type BookStore interface {
	Insert(ctx context.Context, book *repository.Book) error
	Find(ctx context.Context, id uint) (*repository.Book, error)
	FindAll(ctx context.Context) ([]repository.Book, error)
	AddComment(ctx context.Context, bookID uint, text string) (*repository.Comment, error)
	Delete(ctx context.Context, id uint) error
}

// internalError logs the underlying error and responds with a generic 500 so
// implementation details are never leaked to the client.
func internalError(ctx *gin.Context, err error) {
	log.Printf("%s %s: %v", ctx.Request.Method, ctx.Request.URL.Path, err)
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

// BookHandler wires HTTP endpoints to the book repository.
//
// It holds no mutable state, so a single instance can serve concurrent
// requests. Use [NewBookHandler] to build one and
// [BookHandler.RegisterRoutes] to attach it to a router.
type BookHandler struct {
	repo BookStore
}

// NewBookHandler creates a BookHandler backed by the given repository.
//
// The repo must be non-nil; every endpoint delegates to it.
func NewBookHandler(repo BookStore) *BookHandler {
	return &BookHandler{repo: repo}
}

// RegisterRoutes wires the book endpoints onto the provided router group.
//
// The paths are relative to the group's prefix, so a group of "/api/v1"
// produces:
//
//	GET    /books              list every book
//	GET    /books/:id          fetch one book
//	POST   /books              create a book
//	POST   /books/:id/comments add a comment to a book
//	DELETE /books/:id          delete a book
func (handler *BookHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/books", handler.listBooks)
	routerGroup.GET("/books/:id", handler.getBook)
	routerGroup.POST("/books", handler.createBook)
	routerGroup.POST("/books/:id/comments", handler.addComment)
	routerGroup.DELETE("/books/:id", handler.deleteBook)
}

// listBooks handles GET /books and responds with 200 and the JSON array of all
// books, or 500 if the store fails.
func (handler *BookHandler) listBooks(ctx *gin.Context) {
	books, err := handler.repo.FindAll(ctx.Request.Context())
	if err != nil {
		internalError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// getBook handles GET /books/:id and responds with 200 and the book, 400 for a
// non-numeric id, 404 when it does not exist, or 500 if the store fails.
func (handler *BookHandler) getBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := handler.repo.Find(ctx.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		internalError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, book)
}

// createBook handles POST /books. The request body is a book document; title
// and author are required and any supplied comments are stored with it. It
// responds with 201 and the created book, 400 for a malformed or incomplete
// payload, or 500 if the store fails.
func (handler *BookHandler) createBook(ctx *gin.Context) {
	var book repository.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(book.Title) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if strings.TrimSpace(book.Author) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "author is required"})
		return
	}

	if err := handler.repo.Insert(ctx.Request.Context(), &book); err != nil {
		internalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

// addComment handles POST /books/:id/comments. The request body is
// {"text": "..."} and the text must not be blank. It responds with 201 and the
// created comment, 400 for a non-numeric id or invalid body, 404 when the book
// does not exist, or 500 if the store fails.
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

	comment, err := handler.repo.AddComment(ctx.Request.Context(), uint(id), body.Text)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		internalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

// deleteBook handles DELETE /books/:id and responds with 204 on success, 400
// for a non-numeric id, or 500 if the store fails. Deleting an unknown id is
// treated as success.
func (handler *BookHandler) deleteBook(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := handler.repo.Delete(ctx.Request.Context(), uint(id)); err != nil {
		internalError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
