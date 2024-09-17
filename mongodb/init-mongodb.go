package mongodb

import (
	"context"
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
