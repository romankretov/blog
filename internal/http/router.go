package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"local/blog/internal/http/handlers"
	appmw "local/blog/internal/http/middleware"
)

func NewRouter(log zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(appmw.RequestLogger(log))
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(10 * time.Second))

	hh := handlers.NewHealthHandler()
	r.Get("/healthz", hh.Healthz)
	r.Get("/readyz", hh.Readyz)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	return r
}
