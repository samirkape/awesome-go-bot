package regular

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/helpers"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/commands"
	"github.com/samirkape/awesome-go-bot/gobot/constant"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/keyboard"
	"github.com/samirkape/awesome-go-bot/internal/services/internalerrors"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"strconv"
	"strings"
)

type regular struct {
	analytics.Service
	*tgbotapi.BotAPI
	chat.Info
}

var emptyUpdateError = "regular message is nil"
var emptyQueryError = "regular message query is empty"
var invalidQueryError = "query should start with /, try with"
var nonNumericQueryError = "query cannot be parsed, try"

func NewValidatedChat(update *tgbotapi.Update) (chat.Info, error) {
	if update.Message == nil {
		return nil, internalerrors.NewValidationError(emptyUpdateError)
	}

	query := strings.TrimSpace(update.Message.Text)
	if query == "" {
		return nil, internalerrors.NewValidationError(emptyQueryError)
	}

	chatId := update.Message.Chat.ID
	if !strings.HasPrefix(query, constant.CommandPrefix) {
		startCommand := commands.New().Start
		_, err := strconv.Atoi(query)
		if err == nil {
			updatedQuery := fmt.Sprintf("%s%s", constant.CommandPrefix, query)
			return newRegular(chatId, query), internalerrors.NewValidationError(invalidQueryError, updatedQuery)
		} else {
			return newRegular(chatId, query), internalerrors.NewValidationError(nonNumericQueryError, startCommand)
		}
	}

	return newRegular(chatId, query), nil
}

func NewRegularChat(chat chat.Info, analyticsService analytics.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	return &regular{
		Info:    chat,
		Service: analyticsService,
		BotAPI:  botService,
	}, nil
}

func (r *regular) HandleQuery() error {
	chatService := r.Info
	analyticsService := r.Service
	botService := r.BotAPI

	command := commands.New()
	switch chatService.GetQuery() {
	case command.GetStart():
		return gobot.Respond(chatService, botService, constant.SupportedCommands, gobot.WithCustomParsing(tgbotapi.ModeHTML))
	case command.GetDescription():
		return gobot.Respond(chatService, botService, constant.Description)
	case command.GetListCategories():
		return gobot.Respond(chatService, botService, helpers.ListToMessage(analyticsService.GetCategories()))
	case command.IsTopN(chatService.GetQuery()), command.IsCategoryNumber(chatService.GetQuery()):
		keyboardService := keyboard.NewRegularKeyboardChat(chatService, analyticsService, botService)
		return keyboardService.HandleQuery()
	default:
		pkg := analyticsService.GetPackageByName(chatService.GetQuery())
		return gobot.Respond(chatService, botService, helpers.PackageToMsg(pkg, true))
	}
}

func newRegular(chatId int64, query string) *regular {
	return &regular{
		Info: &chat.Chat{
			ChatId: chatId,
			Query:  query,
			Inline: false,
		},
	}
}
