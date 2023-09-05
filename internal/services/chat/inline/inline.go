package inline

import (
	"awesome-go-bot/internal/services/chat"
	"awesome-go-bot/internal/services/packages/analytics/inmemory"
	"awesome-go-bot/internal/services/packages/search"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

type inlineChat struct {
	*tgbotapi.BotAPI
	search.Service
	chat.Info
}

func NewInlineChat(update *tgbotapi.Update, searchService search.Service, botService *tgbotapi.BotAPI) chat.Info {
	query := strings.TrimSpace(update.InlineQuery.Query)
	return &inlineChat{
		Info: &chat.Chat{
			QueryId: update.InlineQuery.ID,
			Query:   query,
			Inline:  true,
		},
		BotAPI:  botService,
		Service: searchService,
	}
}

func (i inlineChat) HandleQuery() error {
	var results []interface{}
	chatService := i.Info
	searchService := i.Service
	botService := i.BotAPI

	packages := searchService.Search(chatService.GetQuery())

	results = createInlineQueryArticle(packages, results)

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: chatService.GetQueryId(),
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	if _, err := botService.Request(inlineConf); err != nil {
		log.Println(err)
	}
	return nil
}

func createInlineQueryArticle(packages []inmemory.Package, results []interface{}) []interface{} {
	for i, _ := range packages {
		article := tgbotapi.NewInlineQueryResultArticle(
			strconv.Itoa(i),
			packages[i].Name,
			packages[i].Name,
		)
		article.Description = "Stars: " + strconv.Itoa(packages[i].Stars) + "\n" + packages[i].Info
		results = append(results, article)
	}
	return results
}