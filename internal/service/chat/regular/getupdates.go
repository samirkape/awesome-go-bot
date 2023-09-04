package regular

import (
	"awesome-go-bot/gopackage"
	"awesome-go-bot/gopackage/helper"
	"awesome-go-bot/internal/service/chat"
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
	case command.IsTopN(chat.GetQuery()):
		topN := packageService.GetTopPackagesSortedByStars(chat.GetQuery())
		messages := helper.BuildStringMessageBatch(topN, true)
		err := gobot.RespondToMessages(chat, botService, messages)
		if err != nil {
			return err
		}
	case command.IsCategoryNumber(chat.GetQuery()):
		packages := packageService.GetPackagesByCategoryNumber(chat.GetQuery())
		messages := helper.BuildStringMessageBatch(packages, false)
		err := gobot.RespondToMessages(chat, botService, messages)
		if err != nil {
			return err
		}
	default:
		pkg := packageService.GetPackageByName(chat.GetQuery())
		err := gobot.Respond(chat, botService, helper.PackageToMsg(pkg, true))
		if err != nil {
			return err
		}
	}
	return nil
}
