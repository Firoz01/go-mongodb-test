package restaurants

import (
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"time"

	"github.com/Firoz01/go-mongodb-test/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func InsertRestaurant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	db := mongodb.GetDatabase()
	restaurantCollection := db.Collection("restaurants")

	var restaurant collections.Restaurant

	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding restaurant data: %v", err)
		return
	}

	restaurant.ID = primitive.NewObjectID()

	result, err := restaurantCollection.InsertOne(ctx, restaurant)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Restaurant already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to insert restaurant", http.StatusInternalServerError)
		}
		log.Printf("Error inserting restaurant: %v", err)
		return
	}

	logEntry := collections.LogEntry{
		Timestamp: time.Now(),
		Message:   "new restaurant data inserted",
		Level:     "INFO",
		Metadata: map[string]interface{}{
			"service": "go-mongo-test",
			"status":  "success",
			"data":    restaurant,
		},
	}

	_ = mongodb.InsertLogEntry(ctx, logEntry)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"inserted_id": result.InsertedID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
