package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Firoz01/go-mongodb-test/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `json:"title"`
	Year            int                `json:"year"`
	Genres          []string           `json:"genres"`
	Href            string             `json:"href"`
	Extract         string             `json:"extract"`
	Thumbnail       string             `json:"thumbnail"`
	ThumbnailWidth  int                `json:"thumbnail_width"`
	ThumbnailHeight int                `json:"thumbnail_height"`
	CastID          primitive.ObjectID `bson:"cast_id,omitempty"`
}

type Cast struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	Cast    []string           `json:"cast"`
}

type InputMovie struct {
	Title           string   `json:"title"`
	Year            int      `json:"year"`
	Cast            []string `json:"cast"`
	Genres          []string `json:"genres"`
	Href            string   `json:"href"`
	Extract         string   `json:"extract"`
	Thumbnail       string   `json:"thumbnail"`
	ThumbnailWidth  int      `json:"thumbnail_width"`
	ThumbnailHeight int      `json:"thumbnail_height"`
}

func main() {

	config.LoadConfig()
	cfg := config.GetConfig()

	client, err := connectToMongoDB(cfg.MongodbURL)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
		return
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(cfg.MongodbDatabaseName)
	moviesCollection := db.Collection("movies")
	castsCollection := db.Collection("casts")

	if err := processJSONFiles("./cmd/seedmongodb/json-files", moviesCollection, castsCollection); err != nil {
		log.Fatalf("Error processing JSON files: %v", err)
		return
	}

	fmt.Println("All files have been processed successfully!")
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func processJSONFiles(dir string, moviesCollection, castsCollection *mongo.Collection) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		if err := processSingleJSONFile(filePath, moviesCollection, castsCollection); err != nil {
			log.Printf("Error processing file %s: %v", filePath, err)
		}
	}
	return nil
}

func processSingleJSONFile(filePath string, moviesCollection, castsCollection *mongo.Collection) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	var inputMovies []InputMovie
	if err := json.Unmarshal(data, &inputMovies); err != nil {
		return fmt.Errorf("error parsing JSON data: %v", err)
	}

	for _, inputMovie := range inputMovies {
		if err := insertMovie(inputMovie, moviesCollection, castsCollection); err != nil {
			log.Printf("Failed to insert movie: %s. Error: %v", inputMovie.Title, err)
		} else {
			fmt.Printf("Successfully inserted movie: %s\n", inputMovie.Title)
		}
	}
	return nil
}

func insertMovie(inputMovie InputMovie, moviesCollection, castsCollection *mongo.Collection) error {

	castID, err := insertCast(inputMovie.Cast, castsCollection)
	if err != nil {
		return fmt.Errorf("failed to insert cast: %v", err)
	}

	movie := Movie{
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

	movieResult, err := moviesCollection.InsertOne(context.TODO(), movie)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %v", err)
	}

	if err := updateCastWithMovieID(castID, movieResult.InsertedID.(primitive.ObjectID), castsCollection); err != nil {
		return fmt.Errorf("failed to update cast with movie ID: %v", err)
	}

	return nil
}

func insertCast(cast []string, castsCollection *mongo.Collection) (primitive.ObjectID, error) {
	castDoc := Cast{
		ID:   primitive.NewObjectID(),
		Cast: cast,
	}

	result, err := castsCollection.InsertOne(context.TODO(), castDoc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func updateCastWithMovieID(castID, movieID primitive.ObjectID, castsCollection *mongo.Collection) error {
	_, err := castsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": castID},
		bson.M{"$set": bson.M{"movie_id": movieID}},
	)
	return err
}
