package search

import "awesome-go-bot-refactored/gopackage"

type Service interface {
	Search(query string) []gopackage.Package
}

func Packages(service Service, query string) []gopackage.Package {
	return service.Search(query)
}
