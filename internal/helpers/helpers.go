package helpers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/config"
	"github.com/samirkape/awesome-go-bot/gobot/constant"
	"github.com/samirkape/awesome-go-bot/internal/errors"
	"github.com/samirkape/awesome-go-bot/internal/logger"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/factory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/mongodb"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
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

func ExecuteCommand(incomingRequest *tgbotapi.Update) error {
	// create new bot
	botService, err := gobot.New(config.WithDefaultConfig())
	if err != nil {
		return err
	}
	// create new chat
	chatInfo, err := factory.New(incomingRequest)
	if err != nil {
		errors.RespondToError(err, botService, chatInfo)
		return err
	}
	// create new mongodb client
	client, err := mongodb.New(mongodb.WithDefaultConfig())
	if err != nil {
		return err
	}
	// create package service from mongodb client
	packageService := packages.NewService(client)
	// create analytics interface from package service
	analyticsService := analytics.NewService(packageService)
	if err != nil {
		return err
	}
	// create new search service based on package service
	searchService := search.NewService(packageService)
	// create new chat service
	chatService, err := factory.NewService(chatInfo, analyticsService, searchService, botService)
	if chatService == nil {
		return err
	}
	logger.FieldLogger("query: ", chatService.GetQuery()).Info("query received")
	return chatService.HandleQuery()
}
