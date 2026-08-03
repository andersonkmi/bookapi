package repository

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Book is the GORM model persisted in the database.
type Book struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title    string    `gorm:"not null" json:"title"`
	Author   string    `gorm:"not null" json:"author"`
	Comments []Comment `gorm:"constraint:OnDelete:CASCADE" json:"comments"`
}

// bookJSON is the wire representation of a Book, exposing comments as plain strings.
type bookJSON struct {
	ID       uint     `json:"id"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Comments []string `json:"comments"`
}

// MarshalJSON renders a Book with its comments as a flat array of strings.
func (b Book) MarshalJSON() ([]byte, error) {
	var comments []string
	for _, comment := range b.Comments {
		comments = append(comments, comment.Text)
	}
	return json.Marshal(bookJSON{
		ID:       b.ID,
		Title:    b.Title,
		Author:   b.Author,
		Comments: comments,
	})
}

// UnmarshalJSON parses a Book whose comments are a flat array of strings.
func (b *Book) UnmarshalJSON(data []byte) error {
	var aux bookJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	b.ID = aux.ID
	b.Title = aux.Title
	b.Author = aux.Author
	b.Comments = nil
	for _, text := range aux.Comments {
		b.Comments = append(b.Comments, Comment{Text: text})
	}
	return nil
}

// Comment is a single free-text comment belonging to a Book.
type Comment struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	BookID uint   `gorm:"not null" json:"book_id"`
	Text   string `gorm:"not null" json:"text"`
}

// BookRepository provides data-access methods for the Book entity.
type BookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a BookRepository backed by the given GORM DB.
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Insert persists a new book and returns the stored record.
func (r *BookRepository) Insert(book *Book) error {
	return r.db.Create(book).Error
}

// Find retrieves a book by its ID, including its comments.
func (r *BookRepository) Find(id uint) (*Book, error) {
	var book Book
	if err := r.db.Preload("Comments").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// FindAll retrieves all books, including their comments.
func (r *BookRepository) FindAll() ([]Book, error) {
	var books []Book
	if err := r.db.Preload("Comments").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// AddComment appends a comment to an existing book and returns the stored comment.
func (r *BookRepository) AddComment(bookID uint, text string) (*Comment, error) {
	if err := r.db.First(&Book{}, bookID).Error; err != nil {
		return nil, err
	}

	comment := Comment{BookID: bookID, Text: text}
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete removes a book by its ID.
func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&Book{}, id).Error
}
