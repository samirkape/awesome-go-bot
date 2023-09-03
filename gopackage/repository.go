package gopackage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
)

type CategoryName string
type AllPackages map[CategoryName][]Package

type Package struct {
	Name     string             `bson:"name" json:"name"`
	URL      string             `bson:"url" json:"url"`
	Info     string             `bson:"info" json:"info"`
	Stars    int                `bson:"stars" json:"stars"`
	Category string             `bson:"title" json:"title"`
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

type GetterInterface interface {
	GetAllPackages() (AllPackages, error)
}

func (a AllPackages) GetCategories() []CategoryName {
	var categories []CategoryName
	for k := range a {
		categories = append(categories, k)
	}
	return categories
}

func (a AllPackages) GetPackagesByCategory(category CategoryName) []Package {
	return a[category]
}

func (a AllPackages) GetPackagesByCategoryNumber(categoryNumber int) []Package {
	return a[a.GetCategories()[categoryNumber]]
}

func (a AllPackages) GetTopPackagesSortedByStars(n int) []Package {
	var packages []Package
	for _, v := range a {
		packages = append(packages, v...)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Stars > packages[j].Stars
	})
	return packages[:n]
}
