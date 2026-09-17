package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prabhatlabs/go-mail-server/internal/lib/database"
	"github.com/prabhatlabs/go-mail-server/internal/lib/env"
	"github.com/prabhatlabs/go-mail-server/internal/lib/middlewares"
	"github.com/prabhatlabs/go-mail-server/internal/services/jobs"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatalln("Env load error: ", err)
	}

	ctx := context.Background()
	database, err := database.ConnectDB(ctx, env.Vars.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middlewares.AccessCodeMiddleware)

	r.Mount("/jobs", jobs.JobsRoutes(database))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
