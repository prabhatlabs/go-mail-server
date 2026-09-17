package jobs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prabhatlabs/go-mail-server/internal/lib/database"
)

type Service struct {
	db *database.DB
}

func JobsRoutes(database *database.DB) http.Handler {
	s := &Service{db: database}
	r := chi.NewRouter()

	r.Get("/", s.listJobs)
	r.Post("/", s.enqueue)
	r.Get("/:id", s.getJob)
	r.Get("/:id/status", s.getJobStatus)

	return r
}
