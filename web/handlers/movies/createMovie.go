package movies

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func insertMovie(ctx context.Context, inputMovie collections.InputMovie, moviesCollection, castsCollection *mongo.Collection) (interface{}, error) {

	castID, err := insertCast(ctx, inputMovie.Cast, castsCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cast: %v", err)
	}

	movie := collections.Movie{
		ID:              primitive.NewObjectID(),
		Title:           inputMovie.Title,
		Year:            inputMovie.Year,
		Genres:          inputMovie.Genres,
		Href:            inputMovie.Href,
		Extract:         inputMovie.Extract,
		Thumbnail:       inputMovie.Thumbnail,
		ThumbnailWidth:  inputMovie.ThumbnailWidth,
		ThumbnailHeight: inputMovie.ThumbnailHeight,
		CastID:          castID,
	}

	movieResult, err := moviesCollection.InsertOne(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("failed to insert movie: %v", err)
	}

	if err := updateCastWithMovieID(ctx, castID, movieResult.InsertedID.(primitive.ObjectID), castsCollection); err != nil {
		return nil, fmt.Errorf("failed to update cast with movie ID: %v", err)
	}

	return movieResult.InsertedID.(primitive.ObjectID), nil
}

func insertCast(ctx context.Context, cast []string, castsCollection *mongo.Collection) (primitive.ObjectID, error) {
	castDoc := collections.Cast{
		ID:   primitive.NewObjectID(),
		Cast: cast,
	}

	result, err := castsCollection.InsertOne(ctx, castDoc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func updateCastWithMovieID(ctx context.Context, castID, movieID primitive.ObjectID, castsCollection *mongo.Collection) error {
	_, err := castsCollection.UpdateOne(
		ctx,
		bson.M{"_id": castID},
		bson.M{"$set": bson.M{"movie_id": movieID}},
	)
	return err
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db := mongodb.GetDatabase()
	movieCollection := db.Collection("movies")
	castCollection := db.Collection("casts")

	inputMovie := collections.InputMovie{}

	if err := json.NewDecoder(r.Body).Decode(&inputMovie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := insertMovie(ctx, inputMovie, movieCollection, castCollection)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logEntry := collections.LogEntry{
		Timestamp: time.Now(),
		Message:   "new movie inserted",
		Level:     "INFO",
		Metadata: map[string]interface{}{
			"service": "order-tracking",
			"status":  "success",
			"data":    inputMovie,
		},
	}

	_ = mongodb.InsertLogEntry(ctx, logEntry)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
