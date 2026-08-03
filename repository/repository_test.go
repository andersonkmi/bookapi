package repository

import (
	"encoding/json"
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
	require.NoError(t, db.AutoMigrate(&Book{}, &Comment{}))

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

func TestInsertWithComments(t *testing.T) {
	repo := newTestRepository(t)

	book := &Book{
		Title:    "With Comments",
		Author:   "Author",
		Comments: []Comment{{Text: "great read"}, {Text: "loved it"}},
	}
	require.NoError(t, repo.Insert(book))

	found, err := repo.Find(book.ID)
	require.NoError(t, err)
	require.Len(t, found.Comments, 2)
	assert.Equal(t, "great read", found.Comments[0].Text)
	assert.Equal(t, book.ID, found.Comments[0].BookID)
}

func TestBookJSONExposesCommentsAsStrings(t *testing.T) {
	book := Book{
		ID:       1,
		Title:    "The Go Programming Language",
		Author:   "Donovan & Kernighan",
		Comments: []Comment{{ID: 10, BookID: 1, Text: "great read"}},
	}

	data, err := json.Marshal(book)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":1,"title":"The Go Programming Language","author":"Donovan & Kernighan","comments":["great read"]}`, string(data))

	var got Book
	require.NoError(t, json.Unmarshal(data, &got))
	require.Len(t, got.Comments, 1)
	assert.Equal(t, "great read", got.Comments[0].Text)
}
