package search

import (
	"awesome-go-bot/internal/services/packages"
	"awesome-go-bot/internal/services/packages/analytics/inmemory"
	searchtrie "awesome-go-bot/internal/services/packages/search/trie"
)

type Service interface {
	Search(string) []inmemory.Package
}

func NewService(a packages.Service) Service {
	return &searchtrie.Search{Service: a}
}
