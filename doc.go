// Command bookapi starts the Book REST API server.
//
// The service exposes a small CRUD API for books and their comments, backed by
// PostgreSQL through GORM and served with the Gin HTTP framework. Routes are
// registered under the /api/v1 prefix by [handlers.BookHandler.RegisterRoutes]:
//
//	GET    /api/v1/books             list all books
//	GET    /api/v1/books/:id         fetch a single book
//	POST   /api/v1/books             create a book
//	POST   /api/v1/books/:id/comments add a comment to a book
//	DELETE /api/v1/books/:id         delete a book
//
// # Configuration
//
// The database connection string is read from the DATABASE_URL environment
// variable. When it is unset, a local development DSN is used
// (host=localhost user=pguser password=pgpwd dbname=bookapi port=5432
// sslmode=disable). The server always listens on port 8080.
//
// # Layout
//
//   - [github.com/andersonkmi/bookapi/repository] holds the persistence models
//     and data-access code.
//   - [github.com/andersonkmi/bookapi/handlers] translates HTTP requests into
//     repository calls and renders JSON responses.
//
// # Usage
//
//	DATABASE_URL="host=db user=pguser password=pgpwd dbname=bookapi port=5432 sslmode=disable" bookapi
package main
