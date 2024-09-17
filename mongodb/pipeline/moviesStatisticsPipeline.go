package pipeline

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildYearRangePipeline() mongo.Pipeline {
	yearRanges := []struct {
		Start int
		End   int
	}{
		{1900, 1910}, {1911, 1920}, {1921, 1930}, {1931, 1940}, {1941, 1950}, {1951, 1960},
		{1961, 1970}, {1971, 1980}, {1981, 1990}, {1991, 2000}, {2001, 2010}, {2011, 2020}, {2021, 2030},
	}

	var branches []bson.M
	for _, yr := range yearRanges {
		branch := bson.M{
			"case": bson.M{
				"$and": []bson.M{
					{"$gte": []interface{}{"$year", yr.Start}},
					{"$lte": []interface{}{"$year", yr.End}},
				},
			},
			"then": bson.M{
				"$concat": []interface{}{
					bson.M{"$toString": yr.Start}, "-", bson.M{"$toString": yr.End},
				},
			},
		}
		branches = append(branches, branch)
	}

	return mongo.Pipeline{
		{{"$project", bson.M{
			"yearRange": bson.M{
				"$switch": bson.M{
					"branches": branches,
					"default":  "Other",
				},
			},
		}}},
		{{"$group", bson.M{
			"_id":   "$yearRange",
			"count": bson.M{"$sum": 1},
		}}},
	}
}
