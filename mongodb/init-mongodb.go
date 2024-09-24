package mongodb

import (
	"context"
	"github.com/Firoz01/go-mongodb-test/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"os"
	"time"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

func InitMongoDB() {
	if client == nil || database == nil {
		client, database = initMongoDBInternal()
	}
}

func initMongoDBInternal() (*mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.GetConfig()
	opts := options.Client().ApplyURI(cfg.MongodbURL)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("failed to connect")
		os.Exit(1)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("Failed to ping MongoDB")
		os.Exit(1)
	}

	database := client.Database(cfg.MongodbDatabaseName)
	slog.Info("Connected to MongoDB!")

	return client, database
}

func Disconnect() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			slog.Error("Failed to disconnect from MongoDB")
		}
		slog.Info("Disconnected from MongoDB!")
	}
}
