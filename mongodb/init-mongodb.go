package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
	database   *mongo.Database
)

func GetClient(ctx context.Context, uri string, dbName string) (*mongo.Client, error) {
	clientOnce.Do(func() {
		opts := options.Client().ApplyURI(uri)
		opts.SetMaxPoolSize(100)
		opts.SetMinPoolSize(5)
		opts.SetMaxConnIdleTime(30 * time.Minute)

		var err error
		client, err = mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Ping the database to verify connection
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Println("Connected to MongoDB!")
		database = client.Database(dbName)

		// Create the time series logs collection if it doesn't exist
		if err := ensureLogsCollectionExists(ctx, database); err != nil {
			log.Fatalf("Failed to create logs collection: %v", err)
		}
	})

	return client, nil
}

func GetDatabase() *mongo.Database {
	return database
}

func Disconnect(ctx context.Context) {
	if client != nil {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}
}

// ensureLogsCollectionExists checks if the "logs" collection exists, and if not, creates it as a time series collection
func ensureLogsCollectionExists(ctx context.Context, db *mongo.Database) error {
	collectionName := "logs"

	// Check if the collection already exists
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	// If the collection exists, return early
	if len(collections) > 0 {
		log.Printf("Collection %s already exists", collectionName)
		return nil
	}

	// Create the time series collection if it doesn't exist
	metadataField := "metadata" // Metadata field should be a pointer to string
	opts := options.CreateCollection().SetTimeSeriesOptions(&options.TimeSeriesOptions{
		TimeField: "timestamp",    // Field for the time series timestamp
		MetaField: &metadataField, // Pointer to the metadata field name
	})

	err = db.CreateCollection(ctx, collectionName, opts)
	if err != nil {
		return fmt.Errorf("failed to create time series collection: %v", err)
	}

	log.Printf("Time series collection %s created.", collectionName)
	return nil
}
