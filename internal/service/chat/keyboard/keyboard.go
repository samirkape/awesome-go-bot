package keyboard

import (
	"awesome-go-bot/gopackage"
	"awesome-go-bot/gopackage/helper"
	"awesome-go-bot/internal/service/chat"
	"awesome-go-bot/internal/service/gobot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var index, messageId int
var packages []gopackage.Package

type keyboardChat struct {
	chat.Info
}

func NewKeyboardChat(update *tgbotapi.Update) chat.Info {
	query := update.CallbackQuery
	return &keyboardChat{
		Info: &chat.Chat{
			ChatId:        update.CallbackQuery.Message.Chat.ID,
			CallBackQuery: query,
			CallBack:      true,
		},
	}
}

func ProcessUsingInlineKeyboard(botService *tgbotapi.BotAPI, packageService gopackage.AllPackages, chat chat.Info) error {
	var err error
	if chat.GetQuery() != "" {
		packages = packageService.GetPackagesByCategoryNumber(chat.GetQuery())
	}
	messages := helper.BuildStringMessageBatch(packages, false)
	if chat.IsCallBack() {
		chat.SetMessageId(messageId)
		switch chat.GetCallBackQuery().Data {
		case "prev":
			if index > 0 {
				index--
			}
		case "next":
			if index < len(messages)-1 {
				index++
			}
		}
	}
	messageId, err = gobot.Respond(chat, botService, messages[index], true)
	if err != nil {
		return err
	}
	return nil
}
