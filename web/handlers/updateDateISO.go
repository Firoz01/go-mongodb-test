package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Firoz01/go-mongodb-test/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func UpdateDateISO(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 6000*time.Second)
	defer cancel()

	db := mongodb.GetDatabase()
	tweetsCollection := db.Collection("tweets")

	// Define the layout for parsing the date
	layout := "Mon Jan 02 15:04:05 -0700 2006"

	// Example filter to select documents where created_at is a string (optional)
	filter := bson.M{
		"created_at": bson.M{"$type": "string"},
	}

	// Find the documents matching the filter
	cursor, err := tweetsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterate over the documents and update the created_at field
	for cursor.Next(ctx) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}

		// Parse the existing created_at string to time.Time
		if createdAtStr, ok := document["created_at"].(string); ok {
			parsedTime, err := time.Parse(layout, createdAtStr)
			if err != nil {
				log.Printf("Failed to parse created_at for document %v: %v", document["_id"], err)
				continue
			}

			// Update the document with the ISODate
			update := bson.M{
				"$set": bson.M{"created_at": parsedTime},
			}

			_, err = tweetsCollection.UpdateOne(ctx, bson.M{"_id": document["_id"]}, update)
			if err != nil {
				log.Printf("Failed to update document %v: %v", document["_id"], err)
			} else {
				fmt.Printf("Successfully updated document %v\n", document["_id"])
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update process completed.")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("success"); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
