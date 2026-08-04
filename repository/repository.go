package repository

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

// Book is the GORM model persisted in the database.
//
// ID is assigned by the database on insert. Title and Author are mandatory.
// Comments is the one-to-many association to [Comment]; it is only populated
// when the query preloads it, and deleting a book cascades to its comments.
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
//
// Only the text of each comment is emitted; comment IDs and book references are
// omitted. A book without comments marshals "comments" as null.
//
// It implements [encoding/json.Marshaler] and returns any error reported by the
// underlying encoder.
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
//
// Each string becomes a [Comment] with only Text set; IDs and BookID are left
// zero and are assigned when the book is persisted. Any previously decoded
// comments on b are discarded.
//
// It implements [encoding/json.Unmarshaler] and returns a decoding error when
// data is not a valid book document.
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
//
// BookID is the foreign key to the owning [Book] and Text holds the comment
// body; both are mandatory.
type Comment struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	BookID uint   `gorm:"not null" json:"book_id"`
	Text   string `gorm:"not null" json:"text"`
}

// BookRepository provides data-access methods for the Book entity.
//
// A BookRepository is safe for concurrent use as long as the underlying
// *gorm.DB is, which is the case for a handle obtained from gorm.Open.
type BookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a BookRepository backed by the given GORM DB.
//
// The db handle must be non-nil and already connected; the repository does not
// run migrations, so the books and comments tables are expected to exist.
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Insert persists a new book and returns the stored record.
//
// The book is updated in place with the generated ID; any comments attached to
// it are inserted in the same operation. It returns the database error when the
// insert fails.
func (r *BookRepository) Insert(ctx context.Context, book *Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

// Find retrieves a book by its ID, including its comments.
//
// It returns gorm.ErrRecordNotFound when no book matches id, or any other
// database error encountered while querying.
func (r *BookRepository) Find(ctx context.Context, id uint) (*Book, error) {
	var book Book
	if err := r.db.WithContext(ctx).Preload("Comments").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// FindAll retrieves all books, including their comments.
//
// An empty table yields an empty slice and a nil error; the result is not
// paginated.
func (r *BookRepository) FindAll(ctx context.Context) ([]Book, error) {
	var books []Book
	if err := r.db.WithContext(ctx).Preload("Comments").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// AddComment appends a comment to an existing book and returns the stored comment.
//
// The book identified by bookID must exist; otherwise gorm.ErrRecordNotFound is
// returned and nothing is written. The returned comment carries its generated
// ID. The text is stored verbatim and is not validated here.
func (r *BookRepository) AddComment(ctx context.Context, bookID uint, text string) (*Comment, error) {
	if err := r.db.WithContext(ctx).First(&Book{}, bookID).Error; err != nil {
		return nil, err
	}

	comment := Comment{BookID: bookID, Text: text}
	if err := r.db.WithContext(ctx).Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete removes a book by its ID.
//
// Its comments are removed as well through the cascading foreign key. Deleting
// an ID that does not exist is a no-op and reports no error.
func (r *BookRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Book{}, id).Error
}
