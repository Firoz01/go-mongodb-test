package handlers

import (
	"context"
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func PatchMovie(w http.ResponseWriter, r *http.Request) {
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

	// Decode the request body into a map for flexible field updates
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare the update document dynamically
	update := bson.M{"$set": bson.M{}}

	// Loop over the provided fields and add them to the update document
	for field, value := range updates {
		// Ensure the fields are valid by checking the movie struct
		switch field {
		case "title", "year", "genres", "href", "extract", "thumbnail", "thumbnailWidth", "thumbnailHeight", "castId":
			update["$set"].(bson.M)[field] = value
		default:
			http.Error(w, "Invalid field in request body", http.StatusBadRequest)
			return
		}
	}

	// Update the movie document only with the fields provided in the request body
	result, err := collection.UpdateByID(ctx, id, update)
	if err != nil || result.MatchedCount == 0 {
		log.Printf("Error updating movie: %v", err)
		http.Error(w, "Error updating movie", http.StatusInternalServerError)
		return
	}

	// Respond with a success message or the updated movie (optional to fetch again)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updates) // This will return the fields that were updated
}
