package packages

import (
	"github.com/samirkape/awesome-go-bot/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Packages []Package

var SortedByStars []Package

type Package struct {
	Name  string             `bson:"name" json:"name"`
	URL   string             `bson:"url" json:"url"`
	Info  string             `bson:"info" json:"info"`
	Title string             `bson:"title" json:"title"`
	Stars int                `bson:"stars" json:"stars"`
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

type AllData struct {
	AllPackages map[string][]Package
	Packages    []Package
	repository.CategoryList
}
