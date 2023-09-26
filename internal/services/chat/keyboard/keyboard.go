package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/helpers"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/commands"
	"github.com/samirkape/awesome-go-bot/internal/errors"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
)

var index, messageId int
var messages []string
var packages []inmemory.Package
var includeCategoryInMessage bool

type keyboardChat struct {
	chat.Info
	analytics.Service
	*tgbotapi.BotAPI
}

func NewChat(update *tgbotapi.Update) (chat.Info, error) {
	if update.CallbackQuery == nil {
		return nil, errors.NewValidationError("callback query is nil")
	}
	query := update.CallbackQuery
	return &keyboardChat{
		Info: &chat.Chat{
			ChatId:        update.CallbackQuery.Message.Chat.ID,
			CallBackQuery: query,
			CallBack:      true,
		},
	}, nil
}

func NewDefaultKeyboardChat(chat chat.Info, analyticsService analytics.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	return &keyboardChat{
		Info:    chat,
		Service: analyticsService,
		BotAPI:  botService,
	}, nil
}

func NewRegularKeyboardChat(messageInfo chat.Info, analyticsService analytics.Service, botService *tgbotapi.BotAPI) chat.Info {
	return &keyboardChat{
		Info: &chat.Chat{
			ChatId:        messageInfo.GetChatId(),
			CallBackQuery: messageInfo.GetCallBackQuery(),
			Query:         messageInfo.GetQuery(),
			CallBack:      false,
		},
		Service: analyticsService,
		BotAPI:  botService,
	}
}

func (k keyboardChat) HandleQuery() error {
	chatService := k.Info
	analyticsService := k.Service
	botService := k.BotAPI
	query := chatService.GetQuery()
	command := commands.New()

	var err error

	if query != "" {
		index = 0 // reset index for new query
	}

	switch query {
	case command.IsTopN(chatService.GetQuery()):
		includeCategoryInMessage = true
		packages = analyticsService.GetTopPackagesSortedByStars(chatService.GetQuery())
	case command.IsCategoryNumber(chatService.GetQuery()):
		includeCategoryInMessage = false
		packages, err = analyticsService.GetPackagesByCategoryNumber(chatService.GetQuery())
		if err != nil {
			return err
		}
	}

	messages = helpers.BuildStringMessageBatch(packages, includeCategoryInMessage)
	if messages == nil {
		return nil
	}

	if chatService.IsCallBack() {
		chatService.SetMessageId(messageId)
		switch chatService.GetCallBackQuery().Data {
		case "prev":
			if index > 0 {
				index--
			}
		case "next":
			if index < len(messages)-1 {
				index++
			} else {
				index = 0
			}
		default:
			return nil
		}
	}

	messageId, err = gobot.RespondToCallBack(chatService, botService, messages[index], index, len(messages))
	if err != nil {
		return err
	}

	return nil
}
