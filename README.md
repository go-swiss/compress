# Compression Middleware

This package provides a middleware for the Go's [net/http](https://pkg.go.dev/net/http) (or other compatible routers) that compresses the response using [brotli](https://en.wikipedia.org/wiki/Brotli) or [gzip](https://en.wikipedia.org/wiki/Gzip) compression algorithms.

It depends only on <https://github.com/andybalholm/brotli>.

## Usage

1. Import the package

    ```go
    import "github.com/go-swiss/compress"
    ```

2. Add `compress.Middleware` to your router.

### With `net/http`

```go
import (
    "net/http"

    "github.com/go-swiss/compress"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", HomeHandler)

    http.ListenAndServe(":8080", compress.Middleware(mux))
}
```

### With `gorilla/mux`

```go
import (
    "net/http"

    "github.com/go-swiss/compress"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.Use(compress.Middleware)
    r.HandleFunc("/", HomeHandler)

    http.ListenAndServe(":8080", r)
}
```

### With `go-chi/chi`

```go
import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-swiss/compress"
)

func main() {
    r := chi.NewRouter()
    r.Use(compress.Middleware)
    r.Get("/", HomeHandler)

    http.ListenAndServe(":8080", r)
}
```
