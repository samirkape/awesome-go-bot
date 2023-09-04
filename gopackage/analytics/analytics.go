package analytics

import "awesome-go-bot/gopackage"

type Service interface {
	GetCategories() []gopackage.CategoryName
	GetPackagesByCategory(gopackage.CategoryName) []gopackage.Package
	GetPackagesByCategoryNumber(string) []gopackage.Package
	GetTopPackagesSortedByStars(string) []gopackage.Package
	GetPackageByName(string) gopackage.Package
}
