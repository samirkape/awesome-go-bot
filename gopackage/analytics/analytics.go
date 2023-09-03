package analytics

import "awesome-go-bot-refactored/gopackage"

type Service interface {
	GetCategories() []gopackage.CategoryName
	GetPackagesByCategory(gopackage.CategoryName) []gopackage.Package
	GetPackagesByCategoryNumber(int) []gopackage.Package
	GetTopPackagesSortedByStars(int) []gopackage.Package
}
