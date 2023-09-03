package search

import (
	"awesome-go-bot-refactored/gopackage"
	searchtrie "awesome-go-bot-refactored/gopackage/search/trie"
)

type Service interface {
	Search(string) []gopackage.Package
}

func NewSearchService(a gopackage.AllPackages) Service {
	return searchtrie.Search{a}
}
