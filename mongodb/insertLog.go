package mongodb

import (
	"context"
	"fmt"
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
)

func InsertLogEntry(ctx context.Context, logEntry collections.LogEntry) error {
	db := GetDatabase()
	collection := db.Collection("logs")

	// Insert the log entry into the time series collection
	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		return err
	}
	fmt.Println("Log entry inserted.")
	return nil
}
