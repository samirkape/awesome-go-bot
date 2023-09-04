package search

import (
	"awesome-go-bot-refactored/gopackage"
	searchtrie "awesome-go-bot-refactored/gopackage/search/trie"
)

type Service interface {
	Search(string) []gopackage.Package
}

func NewService(a gopackage.AllPackages) Service {
	return searchtrie.Search{AllPackages: a}
}
