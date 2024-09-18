package web

import (
	"github.com/Firoz01/go-mongodb-test/web/handlers"
	"github.com/Firoz01/go-mongodb-test/web/handlers/movies"
	"github.com/Firoz01/go-mongodb-test/web/handlers/restaurants"
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
			http.HandlerFunc(movies.Movies),
		),
	)
	mux.Handle(
		"GET /movies_statistics",
		manager.With(
			http.HandlerFunc(movies.MovieStatistics),
		),
	)

	mux.Handle(
		"POST /movies",
		manager.With(
			http.HandlerFunc(movies.CreateMovie),
		),
	)
	mux.Handle(
		"PATCH /movies",
		manager.With(
			http.HandlerFunc(movies.PatchMovie),
		),
	)
	mux.Handle(
		"DELETE /movies",
		manager.With(
			http.HandlerFunc(movies.DeleteMovie),
		),
	)

	mux.Handle(
		"GET /restaurants",
		manager.With(
			http.HandlerFunc(restaurants.FindRestaurants),
		),
	)

	mux.Handle(
		"POST /restaurants",
		manager.With(
			http.HandlerFunc(restaurants.InsertRestaurant),
		),
	)

}
