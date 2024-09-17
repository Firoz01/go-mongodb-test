package handlers

import (
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

// DeleteMovie deletes a movie document by its ID in the MongoDB collection.
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db := mongodb.GetDatabase()
	movieCollection := db.Collection("movies")
	castCollection := db.Collection("casts")

	// Get the movie ID from the URL query
	movieID := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Delete the movie from the collection
	result1, err := movieCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || result1.DeletedCount == 0 {
		log.Printf("Error deleting movie: %v", err)
		http.Error(w, "Error deleting movie", http.StatusInternalServerError)
		return
	}

	result2, err := castCollection.DeleteOne(ctx, bson.M{"movie_id": id})
	if err != nil || result2.DeletedCount == 0 {
		log.Printf("Error deleting movie: %v", err)
		http.Error(w, "Error deleting movie", http.StatusInternalServerError)
		return
	}

	logEntry := collections.LogEntry{
		Timestamp: time.Now(),
		Message:   "movie item deleted",
		Level:     "INFO",
		Metadata: map[string]interface{}{
			"service": "order-tracking",
			"status":  "success",
			"data":    movieID,
		},
	}

	_ = mongodb.InsertLogEntry(ctx, logEntry)

	w.WriteHeader(http.StatusOK)
}
