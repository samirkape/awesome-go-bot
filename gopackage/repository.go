package gopackage

import (
	"github.com/shivamMg/trie"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"strings"
	"sync"
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

func (a AllPackages) Search(query string) []Package {
	index := buildPackageIndex(a)
	query = strings.ToLower(query)
	results := searchTrie(index, query, false)
	return a.buildPackagesFromSearchResults(results)
}

func buildPackageIndex(a AllPackages) *trie.Trie {
	var index *trie.Trie
	var once sync.Once
	once.Do(func() {
		index = initTrie(trie.New(), a)
	})
	return index
}

func initTrie(t *trie.Trie, allPackages AllPackages) *trie.Trie {
	categoryIndex := 0
	for categoryName, packages := range allPackages {
		buildCategoryIndex(t, categoryIndex, categoryName)
		for _, pkg := range packages {
			pkgInfoStrings := getPkgInfoStrings(pkg)
			insertPackageInfo(t, pkgInfoStrings, pkg)
		}
	}
	return t
}

func buildCategoryIndex(t *trie.Trie, categoryIndex int, categoryName CategoryName) {
	categoryParts := strings.Split(string(categoryName), " ")
	for _, part := range categoryParts {
		t.Put([]string{part}, categoryIndex)
	}
}

// getPkgInfoStrings extracts and preprocesses package information into a string slice
func getPkgInfoStrings(pkg Package) []string {
	pkgInfoStrings := []string{}

	// Add package name to the slice
	pkgInfoStrings = append(pkgInfoStrings, pkg.Name)

	// Split the package category and add its parts to the slice
	categoryParts := strings.Split(pkg.Category, " ")
	pkgInfoStrings = append(pkgInfoStrings, categoryParts...)

	// Extract the GitHub repository name from the URL and add it to the slice
	urlParts := strings.Split(strings.TrimPrefix(pkg.URL, "https://github.com/"), "/")
	pkgInfoStrings = append(pkgInfoStrings, urlParts...)

	// Split the package info and add its parts to the slice
	infoParts := strings.Split(pkg.Info, " ")
	pkgInfoStrings = append(pkgInfoStrings, infoParts...)

	// convert all strings to lowercase
	for i := 0; i < len(pkgInfoStrings); i++ {
		pkgInfoStrings[i] = strings.ToLower(pkgInfoStrings[i])
	}
	return pkgInfoStrings
}

// InsertPackageInfo inserts package information strings into the trie
func insertPackageInfo(t *trie.Trie, pkgInfoStrings []string, pkg Package) {
	for index := range pkgInfoStrings {
		indexed := pkgInfoStrings[index:]
		t.Put(indexed, pkg)
		for _, word := range indexed {
			for i := len(word); i > 0; i-- {
				indexedChars := word[:i]
				t.Put([]string{indexedChars}, pkg)
			}
		}
	}
}

func searchTrie(index *trie.Trie, query string, exact bool) *trie.SearchResults {
	var results *trie.SearchResults
	querySlice := strings.Split(query, " ")
	if exact {
		results = index.Search(querySlice, trie.WithMaxResults(1))
	} else {
		results = index.Search(querySlice, trie.WithMaxResults(50))
	}
	return results
}

func (a AllPackages) buildPackagesFromSearchResults(results *trie.SearchResults) []Package {
	var packages []Package
	for _, result := range results.Results {
		categoryNumber, isCategory := result.Value.(int)
		if isCategory {
			packages = append(packages, a.GetPackagesByCategoryNumber(categoryNumber)...)
			return packages
		}
		packages = append(packages, result.Value.(Package))
	}
	return packages
}
