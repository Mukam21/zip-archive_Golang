package api

import (
	"net/http"
	"zip-archive_Golang/model"
	"zip-archive_Golang/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg *model.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	th := service.NewTkHandler(cfg)

	r.Post("/task", th.CreateTk)

	r.Post("/task/{id}/url", th.AddFileFlow)

	r.Get("/task/{id}/status", th.GetStatusWk)

	return r
}
