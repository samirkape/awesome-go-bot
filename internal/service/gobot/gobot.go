// gobot will be responsible for creation of client
// and router that will handle the incoming requests

package gobot

import (
	"awesome-go-bot/internal/service/chat"
	"awesome-go-bot/internal/service/gobot/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func New(config *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.GetToken())
	if err != nil {
		return nil, err
	}
	return bot, nil
}

// Respond will send msg string to user with userid
func Respond(chat chat.Info, bot *tgbotapi.BotAPI, messageText string) error {
	messageConfig := defaultMessageConfig(chat.GetChatId(), messageText)
	_, err := bot.Send(messageConfig)
	if err != nil {
		return fmt.Errorf("message sending failed: %v", err)
	}
	return nil
}

// RespondToMessages will send msg string to user with userid
func RespondToMessages(chat chat.Info, bot *tgbotapi.BotAPI, messages []string) error {
	for _, msg := range messages {
		err := Respond(chat, bot, msg)
		if err != nil {
			logrus.Error("message sending failed: %v", err)
			continue
		}
	}
	return nil
}

func defaultMessageConfig(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		ParseMode:             "markdown",
		Text:                  text,
		DisableWebPagePreview: true,
	}
}

func defaultEditMessageConfig(chatID int64, messageId int, text string) tgbotapi.EditMessageTextConfig {
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageId,
		},
		ParseMode:             "markdown",
		Text:                  text,
		DisableWebPagePreview: true,
	}
}

func trash() tgbotapi.InlineKeyboardMarkup {
	// Create an inline keyboard with prev and next buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData("Previous", "prev")
	nextButton := tgbotapi.NewInlineKeyboardButtonData("Next", "next")
	row := []tgbotapi.InlineKeyboardButton{prevButton, nextButton}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	return keyboard
}

func RespondToCallBack(chat chat.Info, bot *tgbotapi.BotAPI, messageText string, currentPage, totalPages int) (int, error) {
	var sentMessage tgbotapi.Message
	var err error
	if chat.GetMessageId() != 0 {
		messageConfig := defaultEditMessageConfig(chat.GetChatId(), chat.GetMessageId(), messageText)
		markup := createInlineKeyboard(currentPage, totalPages)
		messageConfig.ReplyMarkup = &markup
		sentMessage, err = bot.Send(messageConfig)
		if err != nil {
			return 0, fmt.Errorf("message sending failed: %v", err)
		}
	} else {
		messageConfig := defaultMessageConfig(chat.GetChatId(), messageText)
		markup := createInlineKeyboard(currentPage, totalPages)
		messageConfig.ReplyMarkup = &markup
		sentMessage, err = bot.Send(messageConfig)
		if err != nil {
			return 0, fmt.Errorf("message sending failed: %v", err)
		}
	}
	messageId := sentMessage.MessageID
	return messageId, nil
}

// createInlineKeyboard creates an inline keyboard markup with "Previous," "Total Pages," and "Next" buttons.
func createInlineKeyboard(currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	// Calculate the page navigation buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData("Previous", "prev")
	nextButtonText := fmt.Sprintf("Next (%d/%d)", currentPage+1, totalPages)
	nextButton := tgbotapi.NewInlineKeyboardButtonData(nextButtonText, "next")
	totalPagesButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Total Pages: %d", totalPages), "total_pages")

	// Create a row for navigation buttons and total pages
	navigationRow := []tgbotapi.InlineKeyboardButton{prevButton, totalPagesButton, nextButton}

	// Combine navigation row
	keyboard := tgbotapi.NewInlineKeyboardMarkup(navigationRow)

	return keyboard
}
