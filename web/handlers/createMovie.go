package handlers

import (
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db := mongodb.GetDatabase()
	collection := db.Collection("movies")

	// Decode the request body into a Movie struct
	var movie collections.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set a new ObjectID if one isn't provided
	if movie.ID.IsZero() {
		movie.ID = primitive.NewObjectID()
	}

	// Insert the movie into the collection
	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		log.Printf("Error inserting movie: %v", err)
		http.Error(w, "Error creating movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
