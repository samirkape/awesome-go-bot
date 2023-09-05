package regular

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/helpers"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/gobot/commands"
	"github.com/samirkape/awesome-go-bot/gobot/constant"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/samirkape/awesome-go-bot/internal/services/chat/keyboard"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"strings"
)

type regular struct {
	analytics.Service
	*tgbotapi.BotAPI
	chat.Info
}

func NewRegularChat(update *tgbotapi.Update, analyticsService analytics.Service, botService *tgbotapi.BotAPI) chat.Info {
	query := strings.TrimSpace(update.Message.Text)
	return &regular{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  query,
			Inline: false,
		},
		Service: analyticsService,
		BotAPI:  botService,
	}
}

func (r *regular) HandleQuery() error {
	chatService := r.Info
	analyticsService := r.Service
	botService := r.BotAPI

	command := commands.New()
	switch chatService.GetQuery() {
	case command.GetStart():
		err := gobot.Respond(chatService, botService, constant.Start)
		if err != nil {
			return err
		}
	case command.GetDescription():
		err := gobot.Respond(chatService, botService, constant.Description)
		if err != nil {
			return err
		}
	case command.GetListCategories():
		err := gobot.Respond(chatService, botService, helpers.ListToMessage(analyticsService.GetCategories()))
		if err != nil {
			return err
		}
	case command.IsTopN(chatService.GetQuery()):
		topN := analyticsService.GetTopPackagesSortedByStars(chatService.GetQuery())
		messages := helpers.BuildStringMessageBatch(topN, true)
		err := gobot.RespondToMessages(chatService, botService, messages)
		if err != nil {
			return err
		}
	case command.IsCategoryNumber(chatService.GetQuery()):
		keyboardService := keyboard.NewRegularKeyboardChat(chatService, analyticsService, botService)
		err := keyboardService.HandleQuery()
		if err != nil {
			return err
		}
	default:
		pkg := analyticsService.GetPackageByName(chatService.GetQuery())
		err := gobot.Respond(chatService, botService, helpers.PackageToMsg(pkg, true))
		if err != nil {
			return err
		}
	}
	return nil
}
