package gopackage

import "go.mongodb.org/mongo-driver/bson/primitive"

type CategoryName string

type Package struct {
	Name     string             `bson:"name" json:"name"`
	URL      string             `bson:"url" json:"url"`
	Info     string             `bson:"info" json:"info"`
	Stars    int                `bson:"stars" json:"stars"`
	Category string             `bson:"title" json:"title"`
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

type DbProviderInterface interface {
	GetAllPackages() (map[CategoryName][]Package, error)
}
