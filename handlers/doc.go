// Package handlers implements the HTTP layer of the Book API.
//
// [BookHandler] binds Gin routes to a [BookStore], validates incoming requests
// and renders JSON responses. The concrete store used in production is
// [github.com/andersonkmi/bookapi/repository.BookRepository]; the [BookStore]
// interface exists so tests can substitute a mock.
//
// Errors are reported with a JSON body of the form {"error": "message"} and one
// of the following status codes:
//
//	400 Bad Request           malformed path parameter or invalid/missing payload field
//	404 Not Found             the referenced book does not exist
//	500 Internal Server Error the store returned an unexpected error
package handlers
