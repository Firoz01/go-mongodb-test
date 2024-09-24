package mongodb

import "go.mongodb.org/mongo-driver/mongo"

func GetDatabase() *mongo.Database {
	return database
}
