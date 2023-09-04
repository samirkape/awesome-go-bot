package regular

import (
	"awesome-go-bot/gopackage"
	"awesome-go-bot/gopackage/helper"
	"awesome-go-bot/internal/service/chat"
	"awesome-go-bot/internal/service/chat/keyboard"
	"awesome-go-bot/internal/service/gobot"
	"awesome-go-bot/internal/service/gobot/commands"
	"awesome-go-bot/internal/service/gobot/constant"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type regular struct {
	chat.Info
}

func NewRegularChat(update *tgbotapi.Update) chat.Info {
	query := strings.TrimSpace(update.Message.Text)
	return &regular{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  query,
			Inline: false,
		},
	}
}

func HandleQuery(botService *tgbotapi.BotAPI, packageService gopackage.AllPackages, chat chat.Info) error {
	command := commands.New()
	switch chat.GetQuery() {
	case command.GetStart():
		_, err := gobot.Respond(chat, botService, constant.Start, false)
		if err != nil {
			return err
		}
	case command.GetDescription():
		_, err := gobot.Respond(chat, botService, constant.Description, false)
		if err != nil {
			return err
		}
	case command.GetListCategories():
		_, err := gobot.Respond(chat, botService, helper.ListToMessage(packageService.GetCategories()), false)
		if err != nil {
			return err
		}
	case command.IsTopN(chat.GetQuery()):
		topN := packageService.GetTopPackagesSortedByStars(chat.GetQuery())
		messages := helper.BuildStringMessageBatch(topN, true)
		err := gobot.RespondToMessages(chat, botService, messages)
		if err != nil {
			return err
		}
	case command.IsCategoryNumber(chat.GetQuery()):
		err := keyboard.ProcessUsingInlineKeyboard(botService, packageService, chat)
		if err != nil {
			return err
		}
	default:
		pkg := packageService.GetPackageByName(chat.GetQuery())
		_, err := gobot.Respond(chat, botService, helper.PackageToMsg(pkg, true), false)
		if err != nil {
			return err
		}
	}
	return nil
}
