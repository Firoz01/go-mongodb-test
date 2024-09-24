package restaurants

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Firoz01/go-mongodb-test/logger"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"

	"github.com/Firoz01/go-mongodb-test/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindRestaurants(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	db := mongodb.GetDatabase()
	restaurantCollection := db.Collection("restaurant")

	queryParams := r.URL.Query()

	filter := bson.M{}

	if starsParam := queryParams.Get("stars"); starsParam != "" {
		stars, err := strconv.Atoi(starsParam)
		if err == nil {
			filter["stars"] = bson.M{"$gt": stars}
		}
	}

	if nameParam := queryParams.Get("name"); nameParam != "" {
		filter["name"] = bson.M{"$regex": nameParam, "$options": "i"}
	}

	if categoriesParam := queryParams.Get("categories"); categoriesParam != "" {
		filter["categories"] = bson.M{"$in": []string{categoriesParam}}
	}

	findOptions := options.Find()
	findOptions.SetLimit(10)

	cursor, err := restaurantCollection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
		slog.Error("Failed to fetch restaurants", logger.Extra(map[string]any{
			"error": err.Error(),
		}))
		return
	}
	defer cursor.Close(ctx)

	var restaurants []collections.Restaurant

	if err := cursor.All(ctx, &restaurants); err != nil {
		http.Error(w, "Failed to decode restaurants", http.StatusInternalServerError)
				slog.Error("Error decoding restaurants", logger.Extra(map[string]any{
			"error": err.Error(),
		}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(restaurants); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
