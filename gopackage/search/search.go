package search

import "awesome-go-bot-refactored/gopackage"

type Service interface {
	Search(string) []gopackage.Package
}
