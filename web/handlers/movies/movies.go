package movies

import (
	"context"
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"github.com/Firoz01/go-mongodb-test/mongodb/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func countMovies(ctx context.Context, collection *mongo.Collection, matchStage bson.M) (int64, error) {
	return collection.CountDocuments(ctx, matchStage)
}

func parseMovieQuery(r *http.Request) collections.MovieQuery {
	query := r.URL.Query()
	yearStr := query.Get("year")
	year, _ := strconv.Atoi(yearStr) // Ignoring error, zero value will be used if conversion fails
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 100 // Default limit, adjust as needed
	}

	return collections.MovieQuery{
		Title:  query.Get("title"),
		Year:   year,
		Cast:   query.Get("cast"),
		Genres: strings.Split(query.Get("genres"), ","),
		Page:   page,
		Limit:  limit,
	}
}

func Movies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db := mongodb.GetDatabase()
	collection := db.Collection("movies")

	// Parse and build query
	query := parseMovieQuery(r)

	// Build the match stage for counting movies
	matchStage := bson.M{}

	if query.Title != "" {
		matchStage["title"] = bson.M{"$regex": primitive.Regex{Pattern: query.Title, Options: "i"}}
	}

	if query.Year != 0 {
		matchStage["year"] = query.Year
	}

	if len(query.Genres) > 0 && query.Genres[0] != "" {
		matchStage["genres"] = bson.M{"$in": query.Genres}
	}

	if query.Cast != "" {
		matchStage["cast_info.cast"] = bson.M{"$regex": primitive.Regex{Pattern: query.Cast, Options: "i"}}
	}

	// Count total movies matching the query
	totalMovies, err := countMovies(ctx, collection, matchStage)
	if err != nil {
		log.Printf("Error counting movies: %v", err)
		http.Error(w, "Error counting movies", http.StatusInternalServerError)
		return
	}

	pipeline := pipeline.BuildMoviePipeline(query, matchStage)

	// Fetch the matching movies
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode movies
	var movies []collections.MovieWithCasts
	if err = cursor.All(ctx, &movies); err != nil {
		log.Printf("Error decoding query results: %v", err)
		http.Error(w, "Error decoding query results", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalMovies + int64(query.Limit) - 1) / int64(query.Limit)

	pagination := collections.Pagination{
		Total:       totalMovies,
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		PageSize:    query.Limit,
	}

	// Build the response
	response := collections.MovieWithPagination{
		Movie:      movies,
		Pagination: pagination,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
