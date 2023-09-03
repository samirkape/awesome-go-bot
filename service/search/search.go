package search

import "awesome-go-bot-refactored/gopackage"

type Interface interface {
	Search(query string) ([]gopackage.Package, error)
}

type Search struct {
	provider gopackage.DbProvider
}
