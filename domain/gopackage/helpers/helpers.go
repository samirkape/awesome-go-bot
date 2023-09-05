package helpers

import (
	"awesome-go-bot/internal/services/packages/analytics/inmemory"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// ListToMessage Converts slice of strings into a single string.
func ListToMessage(list []inmemory.CategoryName) string {
	var msg strings.Builder
	for i, pkg := range list {
		index := i + 1
		markdown := fmt.Sprintf("[%d %s](/%v)\n", index, pkg, index)
		msg.WriteString(markdown)
	}
	return msg.String()
}

func BuildStringMessageBatch(packages []inmemory.Package, forTop bool) []string {
	const batchSize = 5
	var batch []string
	for start := 0; start < len(packages); start += batchSize {
		end := start + batchSize
		if end > len(packages) {
			end = len(packages)
		}
		mergedMsg := packagesToMsg(packages[start:end], forTop)
		batch = append(batch, mergedMsg)
	}
	return batch
}

func packagesToMsg(packages []inmemory.Package, forTop bool) string {
	var msg strings.Builder
	for _, pkg := range packages {
		msg.WriteString(PackageToMsg(pkg, forTop))
		msg.WriteString("\n\n")
	}
	return msg.String()
}

func PackageToMsg(pkg inmemory.Package, forTopN bool) string {
	var category string
	name := cases.Title(language.AmericanEnglish).String(pkg.Name)
	stars := fmt.Sprintf("Stars: %d\n", pkg.Stars)
	url := fmt.Sprintf("[%s](%s)\n", name, pkg.URL)
	info := pkg.Info

	if forTopN {
		category = fmt.Sprintf("Category: %s\n", pkg.Category)
	}

	return fmt.Sprintf("%s%s%s%s", url, stars, category, info)
}