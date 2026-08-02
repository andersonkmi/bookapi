package handlers

import (
	"github.com/andersonkmi/bookapi/repository"
	"github.com/stretchr/testify/mock"
)

// mockBookStore is a testify-based mock implementation of BookStore.
type mockBookStore struct {
	mock.Mock
}

func (m *mockBookStore) Insert(book *repository.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookStore) Find(id uint) (*repository.Book, error) {
	args := m.Called(id)
	book, _ := args.Get(0).(*repository.Book)
	return book, args.Error(1)
}

func (m *mockBookStore) FindAll() ([]repository.Book, error) {
	args := m.Called()
	books, _ := args.Get(0).([]repository.Book)
	return books, args.Error(1)
}

func (m *mockBookStore) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
