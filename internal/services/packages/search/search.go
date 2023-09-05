package search

import (
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	searchtrie "github.com/samirkape/awesome-go-bot/internal/services/packages/search/trie"
)

type Service interface {
	Search(string) []inmemory.Package
}

func NewService(a packages.Service) Service {
	return &searchtrie.Search{Service: a}
}
