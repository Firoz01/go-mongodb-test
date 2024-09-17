package web

import (
	"github.com/Firoz01/go-mongodb-test/web/handlers"
	"github.com/Firoz01/go-mongodb-test/web/middlewares"
	"net/http"
)

func InitRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle(
		"GET /hello",
		manager.With(
			http.HandlerFunc(handlers.Hello),
		),
	)
	mux.Handle(
		"GET /movies",
		manager.With(
			http.HandlerFunc(handlers.Movies),
		),
	)
	mux.Handle(
		"GET /movies_statistics",
		manager.With(
			http.HandlerFunc(handlers.MovieStatistics),
		),
	)

	mux.Handle(
		"POST /movies",
		manager.With(
			http.HandlerFunc(handlers.CreateMovie),
		),
	)
	mux.Handle(
		"PATCH /movies",
		manager.With(
			http.HandlerFunc(handlers.PatchMovie),
		),
	)
	mux.Handle(
		"DELETE /movies",
		manager.With(
			http.HandlerFunc(handlers.DeleteMovie),
		),
	)

}
