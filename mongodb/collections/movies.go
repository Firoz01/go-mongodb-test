package collections

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Title           string             `json:"title" bson:"title" validate:"required"`
	Year            int                `json:"year" bson:"year"`
	Genres          []string           `json:"genres" bson:"genres"`
	Href            string             `json:"href" bson:"href"`
	Extract         string             `json:"extract" bson:"extract"`
	Thumbnail       string             `json:"thumbnail" bson:"thumbnail"`
	ThumbnailWidth  int                `json:"thumbnailWidth" bson:"thumbnailwidth"`
	ThumbnailHeight int                `json:"thumbnailHeight" bson:"thumbnailheight"`
	CastID          primitive.ObjectID `json:"castId" bson:"cast_id"`
	Cast            []string           `json:"cast" bson:"cast,omitempty"`
}

type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"currentPage"`
	TotalPages  int64 `json:"totalPages"`
	PageSize    int   `json:"pageSize"`
}
type MovieWithPagination struct {
	Movie      []Movie    `json:"movie"`
	Pagination Pagination `json:"pagination"`
}

type MovieStatsResponse struct {
	TotalMovies int64              `json:"totalMovies"`
	YearStats   map[string]float64 `json:"yearStats"`
}

type MovieQuery struct {
	Title  string
	Year   int
	Cast   string
	Genres []string
	Page   int
	Limit  int
}
