package repository

import (
	"gorm.io/gorm"
)

// Book is the GORM model persisted in the database.
type Book struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"not null" json:"title"`
	Author string `gorm:"not null" json:"author"`
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

// Find retrieves a book by its ID.
func (r *BookRepository) Find(id uint) (*Book, error) {
	var book Book
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// FindAll retrieves all books.
func (r *BookRepository) FindAll() ([]Book, error) {
	var books []Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// Delete removes a book by its ID.
func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&Book{}, id).Error
}
