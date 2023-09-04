package inline

import (
	"awesome-go-bot-refactored/gopackage"
	"awesome-go-bot-refactored/gopackage/search"
	"awesome-go-bot-refactored/internal/service/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type inlineChat struct {
	chat.Info
}

func NewInlineChat(update *tgbotapi.Update) chat.Info {
	return &inlineChat{
		Info: &chat.Chat{
			QueryId: update.InlineQuery.ID,
			Query:   update.InlineQuery.Query,
			Inline:  true,
		},
	}
}

func HandleQuery(botService *tgbotapi.BotAPI, service search.Service, chat chat.Info) error {
	var results []interface{}
	packages := service.Search(chat.GetQuery())

	results = createInlineQueryArticle(packages, results)

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: chat.GetQueryId(),
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	if _, err := botService.Request(inlineConf); err != nil {
		log.Println(err)
	}
	return nil
}

func createInlineQueryArticle(packages []gopackage.Package, results []interface{}) []interface{} {
	for i, _ := range packages {
		article := tgbotapi.NewInlineQueryResultArticle(
			strconv.Itoa(i),
			packages[i].Name,
			packages[i].Name,
		)
		article.Description = "stars: " + strconv.Itoa(packages[i].Stars) + "\n" + packages[i].Info
		results = append(results, article)
	}
	return results
}
