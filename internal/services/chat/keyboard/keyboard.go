package keyboard

import (
	"awesome-go-bot/domain/gopackage/helpers"
	"awesome-go-bot/gobot"
	"awesome-go-bot/internal/services/chat"
	"awesome-go-bot/internal/services/packages/analytics"
	"awesome-go-bot/internal/services/packages/analytics/inmemory"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var index, messageId int
var categoryPackages []inmemory.Package

type keyboardChat struct {
	chat.Info
	analytics.Service
	*tgbotapi.BotAPI
}

func NewDefaultKeyboardChat(update *tgbotapi.Update, analyticsService analytics.Service, botService *tgbotapi.BotAPI) chat.Info {
	query := update.CallbackQuery
	return &keyboardChat{
		Info: &chat.Chat{
			ChatId:        update.CallbackQuery.Message.Chat.ID,
			CallBackQuery: query,
			CallBack:      true,
		},
		Service: analyticsService,
		BotAPI:  botService,
	}
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
	err := error(nil)
	chatService := k.Info
	analyticsService := k.Service
	botService := k.BotAPI

	if chatService.GetQuery() != "" {
		index = 0
		categoryPackages = analyticsService.GetPackagesByCategoryNumber(chatService.GetQuery())
	}

	messages := helpers.BuildStringMessageBatch(categoryPackages, false)

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
		}
	}

	messageId, err = gobot.RespondToCallBack(chatService, botService, messages[index], index, len(messages))
	if err != nil {
		return err
	}

	return nil
}
