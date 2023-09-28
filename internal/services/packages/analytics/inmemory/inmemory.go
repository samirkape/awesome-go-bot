package inmemory

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func (a AllPackages) GetCategories() []CategoryName {
	var categories []CategoryName
	for k := range a {
		categories = append(categories, k)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i] < categories[j]
	})
	return categories
}

func (a AllPackages) GetPackagesByCategory(category CategoryName) []Package {
	return a[category]
}

// TODO optimize this
func (a AllPackages) GetPackageByName(name string) Package {
	for _, v := range a {
		for _, pkg := range v {
			if strings.ToLower(pkg.Name) == name {
				return pkg
			}
		}
	}
	return Package{}
}

func (a AllPackages) GetPackagesByCategoryNumber(query string) ([]Package, error) {
	if strings.HasPrefix(query, "/") {
		query = strings.TrimPrefix(query, "/")
	}
	categoryNumber, err := strconv.Atoi(query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	// sort categories
	if categoryNumber > len(a.GetCategories()) {
		return nil, errors.New("category number is out of range")
	}
	packages := a[a.GetCategories()[categoryNumber-1]]
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Stars > packages[j].Stars
	})
	return packages, nil
}

// TODO optimize this
func (a AllPackages) GetTopPackagesSortedByStars(query string) []Package {
	var packages []Package
	for _, v := range a {
		packages = append(packages, v...)
	}
	n, err := getNumberOutOfQuery(query)
	if err != nil {
		n = len(packages)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Stars > packages[j].Stars
	})
	return packages[:n]
}

func getNumberOutOfQuery(query string) (int, error) {
	re := regexp.MustCompile(`(\d+)$`)
	match := re.FindStringSubmatch(query)
	if len(match) == 2 {
		// Convert the matched string to an integer
		number, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		return number, nil
	}
	return 0, fmt.Errorf("no number found in the input string")
}
