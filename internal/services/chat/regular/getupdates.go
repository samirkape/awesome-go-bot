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
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics"
	"strings"
)

type regular struct {
	analytics.Service
	*tgbotapi.BotAPI
	chat.Info
}

func NewRegularChat(update *tgbotapi.Update, analyticsService analytics.Service, botService *tgbotapi.BotAPI) (chat.Info, error) {
	if update.Message == nil {
		return nil, fmt.Errorf("regular message is nil: %+v", update)
	}
	query := strings.TrimSpace(update.Message.Text)
	return &regular{
		Info: &chat.Chat{
			ChatId: update.Message.Chat.ID,
			Query:  query,
			Inline: false,
		},
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
		return gobot.Respond(chatService, botService, constant.Start)
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
