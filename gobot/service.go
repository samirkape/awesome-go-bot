package gobot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
	"github.com/sirupsen/logrus"
)

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

func RespondToCallBack(chat chat.Info, bot *tgbotapi.BotAPI, messageText string, currentPage, totalPages int) (int, error) {
	var sentMessage tgbotapi.Message
	var err error
	messageConfig := getMessageConfig(chat, currentPage, totalPages, messageText)
	sentMessage, err = bot.Send(messageConfig)
	if err != nil {
		return 0, fmt.Errorf("message sending failed: %v", err)
	}
	messageId := sentMessage.MessageID
	return messageId, nil
}

func getMessageConfig(chat chat.Info, currentPage, totalPages int, messageText string) tgbotapi.Chattable {
	if chat.GetMessageId() != 0 {
		messageConfig := defaultEditMessageConfig(chat.GetChatId(), chat.GetMessageId(), messageText)
		markup := createInlineKeyboard(currentPage, totalPages)
		messageConfig.ReplyMarkup = &markup
		return messageConfig
	} else {
		messageConfig := defaultMessageConfig(chat.GetChatId(), messageText)
		markup := createInlineKeyboard(currentPage, totalPages)
		messageConfig.ReplyMarkup = &markup
		return messageConfig
	}
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
