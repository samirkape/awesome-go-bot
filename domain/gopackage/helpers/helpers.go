package helpers

import (
	"fmt"
	"github.com/samirkape/awesome-go-bot/gobot/constant"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// CategoriesToMessage Converts slice of strings into a single string.
func CategoriesToMessage(list []inmemory.CategoryName) string {
	var msg strings.Builder
	for i, pkg := range list {
		index := i + 1
		markdown := fmt.Sprintf("%d %s\n", index, pkg)
		msg.WriteString(markdown)
	}
	msg.WriteString("\n")
	msg.WriteString(constant.CategoryHelper)
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

func PackageToMsg(pkg inmemory.Package, forTopN bool) string {
	var category string
	name := cases.Title(language.AmericanEnglish).String(pkg.Name)
	stars := fmt.Sprintf("★ %d\n", pkg.Stars)
	url := fmt.Sprintf("[%s](%s)\n", name, pkg.URL)
	info := pkg.Info

	if forTopN {
		category = fmt.Sprintf("Category: %s\n", pkg.Category)
	}

	return fmt.Sprintf("%s%s%s%s", url, stars, category, info)
}

func packagesToMsg(packages []inmemory.Package, forTop bool) string {
	var msg strings.Builder
	for _, pkg := range packages {
		msg.WriteString(PackageToMsg(pkg, forTop))
		msg.WriteString("\n\n")
	}
	return msg.String()
}
