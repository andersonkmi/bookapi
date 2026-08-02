# Book API toy project in Go

A small REST API for managing books, written in Go with the [Gin](https://github.com/gin-gonic/gin) web framework.

Books are held in an in-memory map guarded by a mutex, so all data is lost when the process exits. There is no database or persistence layer.

## Requirements

- Go 1.26.4 or later

## Running

```bash
go run .
```

The server listens on `:8080`. Trusted proxies are disabled, so client IPs are taken from the connection rather than forwarding headers.

To build a binary instead:

```bash
go build -o bookapi .
./bookapi
```

Gin runs in debug mode by default; set `GIN_MODE=release` to quiet the startup warnings and request logging.

## API

All routes are served under the `/api/v1` prefix.

| Method   | Path              | Description                  | Success response      |
| -------- | ----------------- | ---------------------------- | --------------------- |
| `GET`    | `/api/v1/books`   | List all books               | `200` with JSON array |
| `GET`    | `/api/v1/books/:id` | Fetch a single book by ID  | `200` with the book   |
| `POST`   | `/api/v1/books`   | Create a book                | `201` with the book   |
| `DELETE` | `/api/v1/books/:id` | Delete a book by ID        | `204` with no body    |

### Book

```json
{
  "id": 1,
  "title": "The Go Programming Language",
  "author": "Donovan & Kernighan"
}
```

`title` and `author` are required on create. `id` is assigned by the server from an incrementing counter and any value sent by the client is overwritten.

### Errors

Errors are returned as `{"error": "..."}`:

- `400 Bad Request` — the `:id` path parameter is not an integer, or the request body fails validation.
- `404 Not Found` — no book exists with the given ID.

Note that `GET /api/v1/books` returns books in map iteration order, which is not stable between calls.

## Examples

```bash
# Create a book
curl -X POST http://localhost:8080/api/v1/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"The Go Programming Language","author":"Donovan & Kernighan"}'

# List all books
curl http://localhost:8080/api/v1/books

# Fetch one
curl http://localhost:8080/api/v1/books/1

# Delete one
curl -X DELETE http://localhost:8080/api/v1/books/1
```

## Layout

```
main.go              # server setup, /api/v1 route group, listener on :8080
handlers/handlers.go # Book model, in-memory store, route registration and handlers
```

## License

MIT — see [LICENSE](LICENSE).
