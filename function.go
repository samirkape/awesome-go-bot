package awesome_go_bot

import (
	"awesome-go-bot-refactored/gopackage"
	"awesome-go-bot-refactored/gopackage/mongodb"
	"awesome-go-bot-refactored/gopackage/search"
	"awesome-go-bot-refactored/internal/logger"
	"awesome-go-bot-refactored/service/chat"
	"awesome-go-bot-refactored/service/chat/inline"
	"awesome-go-bot-refactored/service/chat/regular"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

var queryError = errors.New("unable to handle query")

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := parseRequest(r.Body)
	if err != nil {
		log.Fatal("parsing error")
	}
	chat := getChat(ctx, request)
	err = ExecuteCommand(ctx, chat)
	if err != nil {
		logger.FieldLogger("failed for chatId: ", request.Message.Chat.ID).Error(err)
	}
}

func ExecuteCommand(ctx context.Context, chat chat.Info) error {
	var botService *tgbotapi.BotAPI // TODO create botService
	client := mongodb.GetClient()
	packages, err := client.GetAllPackages()
	if err != nil {
		return err
	}
	searchService := search.NewSearchService(packages)
	if chat.IsInline() {
		handleInlineQuery(botService, searchService, chat)
		fmt.Println("inline query")
	} else {
		fmt.Println("regular query")
	}
	return nil
}

func handleInlineQuery(botService *tgbotapi.BotAPI, service search.Service, chat chat.Info) {
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

func parseRequest(body io.ReadCloser) (*tgbotapi.Update, error) {
	var update *tgbotapi.Update
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&update)
	if err != nil {
		return nil, err
	}
	return update, nil
}

func getChat(ctx context.Context, request *tgbotapi.Update) chat.Info {
	if request.InlineQuery != nil {
		return inline.NewInlineChat(request)
	} else {
		return regular.NewRegularChat(request)
	}
}
