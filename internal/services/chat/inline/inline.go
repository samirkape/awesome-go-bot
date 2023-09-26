package inline

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/internal/errors"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/search"
	"log"
	"strconv"
	"strings"
)

type inlineChat struct {
	*tgbotapi.BotAPI
	search.Service
	chat.Info
}

func NewChat(update *tgbotapi.Update) (chat.Info, error) {
	if update.InlineQuery == nil {
		return &chat.Chat{}, errors.NewValidationError("inline query is nil")
	}
	query := strings.TrimSpace(update.InlineQuery.Query)
	return inlineChat{Info: &chat.Chat{
		QueryId: update.InlineQuery.ID,
		Query:   query,
		Inline:  true,
	}}, nil
}

func NewInlineChat(chat chat.Info, searchService search.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	return &inlineChat{
		Info:    chat,
		BotAPI:  botService,
		Service: searchService,
	}, nil
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
	for i := range packages {
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
