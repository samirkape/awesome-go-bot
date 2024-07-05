package searchtrie

import (
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	"github.com/shivamMg/trie"
	"slices"
	"sort"
	"strings"
	"sync"
)

type Search struct {
	index *trie.Trie
	once  sync.Once
	packages.Service
}

func (s *Search) Search(query string) []inmemory.Package {
	if query == "" {
		return nil
	}
	s.once.Do(func() {
		s.index = buildPackageIndex(s.GetAllPackages())
	})
	query = strings.ToLower(query)
	results := searchTrie(s.index, query, false)
	return buildPackagesFromSearchResults(s.GetAllPackages(), results)
}

func buildPackageIndex(a inmemory.AllPackages) *trie.Trie {
	var index *trie.Trie
	index = initTrie(trie.New(), a)
	return index
}

func initTrie(t *trie.Trie, allPackages inmemory.AllPackages) *trie.Trie {
	//categoryIndex := 0
	for _, pkgs := range allPackages {
		//buildCategoryIndex(t, categoryIndex, categoryName)
		for _, pkg := range pkgs {
			pkgInfoStrings := getPkgInfoStrings(pkg)
			insertPackageInfo(t, pkgInfoStrings, pkg)
		}
	}
	return t
}

func buildCategoryIndex(t *trie.Trie, categoryIndex int, categoryName inmemory.CategoryName) {
	categoryParts := strings.Split(string(categoryName), " ")
	for _, part := range categoryParts {
		t.Put([]string{part}, categoryIndex)
	}
}

// getPkgInfoStrings extracts and preprocesses package information into a string slice
func getPkgInfoStrings(pkg inmemory.Package) []string {
	var pkgInfoStrings []string

	// Add package name to the slice
	pkgInfoStrings = append(pkgInfoStrings, splitString(pkg.Name, "")...)

	// Split the package category and add its parts to the slice
	pkgInfoStrings = append(pkgInfoStrings, splitString(pkg.Category, "")...)

	// Extract the GitHub repository name from the URL and add it to the slice
	urlParts := strings.Split(strings.TrimPrefix(pkg.URL, "https://github.com/"), "/")
	for _, urlPart := range urlParts {
		pkgInfoStrings = append(pkgInfoStrings, splitString(urlPart, "")...)
		break // no need to add the rest of the parts of the URL since it is already added in the package name
	}

	// Split the package info and add its parts to the slice
	pkgInfoStrings = append(pkgInfoStrings, splitString(pkg.Info, " ")...)

	// convert all strings to lowercase
	for i := 0; i < len(pkgInfoStrings); i++ {
		pkgInfoStrings[i] = strings.ToLower(pkgInfoStrings[i])
	}
	return pkgInfoStrings
}

func splitString(pkg string, ch string) []string {
	var stringToChars []string
	chars := strings.Split(pkg, ch)
	stringToChars = append(stringToChars, chars...)
	return stringToChars
}

// InsertPackageInfo inserts package information strings into the trie
func insertPackageInfo(t *trie.Trie, pkgInfoStrings []string, pkg inmemory.Package) {
	for index := range pkgInfoStrings {
		indexed := pkgInfoStrings[index:]
		t.Put(indexed, pkg)
		for _, word := range indexed {
			for i := 1; i <= len(word); i++ {
				indexedChars := word[:i]
				t.Put([]string{indexedChars}, pkg)
			}
		}
	}
}

func searchTrie(index *trie.Trie, query string, exact bool) *trie.SearchResults {
	var results *trie.SearchResults
	query = strings.TrimSuffix(query, " ")
	querySlice := strings.Split(query, "")
	if exact {
		results = index.Search(querySlice, trie.WithMaxResults(1))
	} else {
		results = index.Search(querySlice, trie.WithMaxResults(50))
	}
	return results
}

func buildPackagesFromSearchResults(a inmemory.AllPackages, results *trie.SearchResults) []inmemory.Package {
	var packages []inmemory.Package
	for _, result := range results.Results {
		pkg := result.Value.(inmemory.Package)
		pkg.Name = strings.ToLower(pkg.Name)
		packages = append(packages, pkg)
	}
	packages = sortByStarsAndCleanup(packages)
	return packages
}

func sortByStarsAndCleanup(packages []inmemory.Package) []inmemory.Package {
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Stars > packages[j].Stars
	})
	packages = slices.Compact(packages)
	return packages
}
