package analytics

import (
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
)

type Service interface {
	GetCategories() []inmemory.CategoryName
	GetPackagesByCategory(inmemory.CategoryName) []inmemory.Package
	GetPackagesByCategoryNumber(string) []inmemory.Package
	GetTopPackagesSortedByStars(string) []inmemory.Package
	GetPackageByName(string) inmemory.Package
}

func NewService(getter packages.Service) Service {
	return getter.GetAllPackages()
}
