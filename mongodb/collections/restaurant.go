package collections

import "go.mongodb.org/mongo-driver/bson/primitive"

type Restaurant struct {
	ID         primitive.ObjectID `bson:"_id"`
	Grades     []int              `bson:"grades"`
	Name       string             `bson:"name"`
	Contact    Contact            `bson:"contact"`
	Stars      int                `bson:"stars"`
	Categories []string           `bson:"categories"`
}

// Contact represents the nested contact details
type Contact struct {
	Phone    string     `bson:"phone"`
	Email    string     `bson:"email"`
	Location [2]float64 `bson:"location"` // [longitude, latitude]
}
