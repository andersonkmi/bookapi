package handlers

import (
	"context"

	"github.com/andersonkmi/bookapi/repository"
	"github.com/stretchr/testify/mock"
)

// mockBookStore is a testify-based mock implementation of BookStore.
type mockBookStore struct {
	mock.Mock
}

func (m *mockBookStore) Insert(ctx context.Context, book *repository.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *mockBookStore) Find(ctx context.Context, id uint) (*repository.Book, error) {
	args := m.Called(ctx, id)
	book, _ := args.Get(0).(*repository.Book)
	return book, args.Error(1)
}

func (m *mockBookStore) FindAll(ctx context.Context) ([]repository.Book, error) {
	args := m.Called(ctx)
	books, _ := args.Get(0).([]repository.Book)
	return books, args.Error(1)
}

func (m *mockBookStore) AddComment(ctx context.Context, bookID uint, text string) (*repository.Comment, error) {
	args := m.Called(ctx, bookID, text)
	comment, _ := args.Get(0).(*repository.Comment)
	return comment, args.Error(1)
}

func (m *mockBookStore) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
