package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestRepository(t *testing.T) *BookRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&Book{}))

	return NewBookRepository(db)
}

func TestInsertAndFind(t *testing.T) {
	repo := newTestRepository(t)

	book := &Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"}
	require.NoError(t, repo.Insert(book))
	require.NotZero(t, book.ID)

	found, err := repo.Find(book.ID)
	require.NoError(t, err)
	assert.Equal(t, book.Title, found.Title)
	assert.Equal(t, book.Author, found.Author)
}

func TestFindAll(t *testing.T) {
	repo := newTestRepository(t)

	require.NoError(t, repo.Insert(&Book{Title: "Book A", Author: "Author A"}))
	require.NoError(t, repo.Insert(&Book{Title: "Book B", Author: "Author B"}))

	books, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, books, 2)
}

func TestDelete(t *testing.T) {
	repo := newTestRepository(t)

	book := &Book{Title: "Ephemeral", Author: "Author"}
	require.NoError(t, repo.Insert(book))

	require.NoError(t, repo.Delete(book.ID))

	_, err := repo.Find(book.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
