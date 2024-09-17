package pipeline

import (
	"github.com/Firoz01/go-mongodb-test/mongodb/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildMoviePipeline(query collections.MovieQuery, matchStage bson.M) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "casts"},
					{"localField", "cast_id"},
					{"foreignField", "_id"},
					{"as", "cast_info"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$cast_info"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", matchStage}})
	}

	pipeline = append(pipeline,
		bson.D{{"$project", bson.M{
			"_id":             1,
			"title":           1,
			"year":            1,
			"genres":          1,
			"href":            1,
			"extract":         1,
			"thumbnail":       1,
			"thumbnailwidth":  1,
			"thumbnailheight": 1,
			"cast_id":         1,
			"cast":            "$cast_info.cast",
		}}},
		bson.D{{"$skip", (query.Page - 1) * query.Limit}},
		bson.D{{"$limit", query.Limit}},
	)

	return pipeline
}
