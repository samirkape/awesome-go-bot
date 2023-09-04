package search

import (
	"awesome-go-bot/gopackage"
	searchtrie "awesome-go-bot/gopackage/search/trie"
)

type Service interface {
	Search(string) []gopackage.Package
}

func NewService(a gopackage.AllPackages) Service {
	return searchtrie.Search{AllPackages: a}
}
