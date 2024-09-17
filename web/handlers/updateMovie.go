package handlers

import (
	"context"
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

// UpdateMovie updates a movie document by its ID in the MongoDB collection.
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db := mongodb.GetDatabase()
	collection := db.Collection("movies")

	// Get the movie ID from the URL query
	movieID := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Decode the request body into a Movie struct
	var movie collections.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the movie in the collection
	update := bson.M{
		"$set": movie,
	}

	result, err := collection.UpdateByID(ctx, id, update)
	if err != nil || result.MatchedCount == 0 {
		log.Printf("Error updating movie: %v", err)
		http.Error(w, "Error updating movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}
