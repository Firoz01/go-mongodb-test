package handlers

import (
	"context"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

// DeleteMovie deletes a movie document by its ID in the MongoDB collection.
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
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

	// Delete the movie from the collection
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || result.DeletedCount == 0 {
		log.Printf("Error deleting movie: %v", err)
		http.Error(w, "Error deleting movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
