package handlers

import (
	"encoding/json"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"github.com/Firoz01/go-mongodb-test/mongodb/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"net/http"
	"strings"
)

func MovieStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db := mongodb.GetDatabase()
	moviesCollection := db.Collection("movies")

	query := r.URL.Query()
	genresParam := query.Get("genres")
	genres := strings.Split(genresParam, ",")

	// Case-insensitive genre matching
	matchStage := bson.M{}
	if len(genres) > 0 && genres[0] != "" {
		var genreConditions []bson.M
		for _, genre := range genres {
			genreConditions = append(genreConditions, bson.M{"genres": bson.M{"$regex": genre, "$options": "i"}})
		}
		matchStage["$or"] = genreConditions
	}

	yearPipeline := pipeline.BuildYearRangePipeline()

	// Count total movies
	totalMovies, err := moviesCollection.CountDocuments(ctx, matchStage)
	if err != nil {
		log.Printf("Error counting movies: %v", err)
		http.Error(w, "Error counting movies", http.StatusInternalServerError)
		return
	}

	// Year statistics aggregation
	yearCursor, err := moviesCollection.Aggregate(ctx, yearPipeline)
	if err != nil {
		log.Printf("Error executing year aggregation: %v", err)
		http.Error(w, "Error executing year aggregation", http.StatusInternalServerError)
		return
	}
	defer yearCursor.Close(ctx)

	yearStats := make(map[string]float64)
	for yearCursor.Next(ctx) {
		var result struct {
			Range string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := yearCursor.Decode(&result); err != nil {
			log.Printf("Error decoding year result: %v", err)
			continue
		}
		if totalMovies > 0 {
			percentage := (float64(result.Count) / float64(totalMovies)) * 100
			// Round to 2 decimal places for readability
			percentage = math.Round(percentage*100) / 100
			if percentage > 0 {
				yearStats[result.Range] = percentage
			}
		} else {
			yearStats[result.Range] = 0
		}
	}

	// Build the response
	response := collections.MovieStatsResponse{
		TotalMovies: totalMovies,
		YearStats:   yearStats,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
