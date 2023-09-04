package regular

import (
	"awesome-go-bot-refactored/gopackage"
	"awesome-go-bot-refactored/gopackage/helper"
	"awesome-go-bot-refactored/internal/service/chat"
	"awesome-go-bot-refactored/internal/service/gobot"
	"awesome-go-bot-refactored/internal/service/gobot/commands"
	"awesome-go-bot-refactored/internal/service/gobot/constant"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type regular struct {
	chat.Info
}

func NewRegularChat(update *tgbotapi.Update) chat.Info {
	return &regular{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  update.Message.Text,
			Inline: false,
		},
	}
}

func HandleQuery(botService *tgbotapi.BotAPI, packageService gopackage.AllPackages, chat chat.Info) error {
	command := commands.New()
	switch chat.GetQuery() {
	case command.GetStart():
		err := gobot.Respond(chat, botService, constant.Start)
		if err != nil {
			return err
		}
	case command.GetDescription():
		err := gobot.Respond(chat, botService, constant.Description)
		if err != nil {
			return err
		}
	case command.GetListCategories():
		err := gobot.Respond(chat, botService, helper.ListToMessage(packageService.GetCategories()))
		if err != nil {
			return err
		}
	}
	return nil
}
