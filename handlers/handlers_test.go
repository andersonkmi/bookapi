package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andersonkmi/bookapi/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type HandlersSuite struct {
	suite.Suite
	store  *mockBookStore
	router *gin.Engine
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

func (s *HandlersSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *HandlersSuite) SetupTest() {
	s.store = new(mockBookStore)
	handler := NewBookHandler(s.store)

	s.router = gin.New()
	handler.RegisterRoutes(s.router.Group("/api/v1"))
}

func (s *HandlersSuite) request(method, path string, body any) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		s.Require().NoError(err)
		reader = bytes.NewReader(payload)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	return rec
}

func (s *HandlersSuite) TestListBooks() {
	want := []repository.Book{
		{ID: 1, Title: "Book A", Author: "Author A"},
		{ID: 2, Title: "Book B", Author: "Author B"},
	}
	s.store.On("FindAll").Return(want, nil)

	rec := s.request(http.MethodGet, "/api/v1/books", nil)

	s.Equal(http.StatusOK, rec.Code)
	var got []repository.Book
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	s.Equal(want, got)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestGetBookFound() {
	want := &repository.Book{ID: 1, Title: "Book A", Author: "Author A"}
	s.store.On("Find", uint(1)).Return(want, nil)

	rec := s.request(http.MethodGet, "/api/v1/books/1", nil)

	s.Equal(http.StatusOK, rec.Code)
	var got repository.Book
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	s.Equal(*want, got)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestGetBookNotFound() {
	s.store.On("Find", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	rec := s.request(http.MethodGet, "/api/v1/books/99", nil)

	s.Equal(http.StatusNotFound, rec.Code)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestGetBookInvalidID() {
	rec := s.request(http.MethodGet, "/api/v1/books/abc", nil)

	s.Equal(http.StatusBadRequest, rec.Code)
	s.store.AssertNotCalled(s.T(), "Find", mock.Anything)
}

func (s *HandlersSuite) TestCreateBook() {
	s.store.On("Insert", mock.AnythingOfType("*repository.Book")).Return(nil)

	rec := s.request(http.MethodPost, "/api/v1/books", repository.Book{Title: "New", Author: "Author"})

	s.Equal(http.StatusCreated, rec.Code)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestAddComment() {
	want := &repository.Comment{ID: 5, BookID: 1, Text: "great read"}
	s.store.On("AddComment", uint(1), "great read").Return(want, nil)

	rec := s.request(http.MethodPost, "/api/v1/books/1/comments", map[string]string{"text": "great read"})

	s.Equal(http.StatusCreated, rec.Code)
	var got repository.Comment
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	s.Equal(*want, got)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestAddCommentBookNotFound() {
	s.store.On("AddComment", uint(99), "x").Return(nil, gorm.ErrRecordNotFound)

	rec := s.request(http.MethodPost, "/api/v1/books/99/comments", map[string]string{"text": "x"})

	s.Equal(http.StatusNotFound, rec.Code)
	s.store.AssertExpectations(s.T())
}

func (s *HandlersSuite) TestAddCommentMissingText() {
	rec := s.request(http.MethodPost, "/api/v1/books/1/comments", map[string]string{})

	s.Equal(http.StatusBadRequest, rec.Code)
	s.store.AssertNotCalled(s.T(), "AddComment", mock.Anything, mock.Anything)
}

func (s *HandlersSuite) TestDeleteBook() {
	s.store.On("Delete", uint(1)).Return(nil)

	rec := s.request(http.MethodDelete, "/api/v1/books/1", nil)

	s.Equal(http.StatusNoContent, rec.Code)
	s.store.AssertExpectations(s.T())
}
