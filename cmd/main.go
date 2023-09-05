package main

import (
	"awesome-go-bot/domain/gopackage/mongodb"
	"awesome-go-bot/gobot"
	"awesome-go-bot/gobot/config"
	"awesome-go-bot/internal/services/chat/factory"
	"awesome-go-bot/internal/services/packages"
	"awesome-go-bot/internal/services/packages/analytics"
	"awesome-go-bot/internal/services/packages/search"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func main() {
	URL, found := os.LookupEnv("ATLAS_URI")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}

	token, found := os.LookupEnv("TEST_TOKEN")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}

	botService, err := gobot.New(&config.Config{Token: token})
	if err != nil {
		return
	}

	dbConfig := mongodb.NewConfig(mongodb.TABLENAME, URL)
	client, err := mongodb.New(dbConfig)

	packageService := packages.NewService(client)
	// get all analyticsService from the database
	analyticsService := analytics.NewService(packageService)
	// create new search service
	searchService := search.NewService(packageService)
	// create new chat

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := botService.GetUpdatesChan(u)

	for update := range updates {
		chatService := factory.NewChatService(&update, packageService, analyticsService, searchService, botService)
		err := chatService.HandleQuery()
		if err != nil {
			return
		}
	}
}

func tryKeyboard(updates tgbotapi.UpdatesChannel, messages []string, goBot *tgbotapi.BotAPI) {
	var index, messageId int
	for update := range updates {
		if update.Message != nil {
			// Handle incoming messages here
			// For example, send a message with an inline keyboard
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages[index])
			msg.ReplyMarkup = createInlineKeyboard()

			sentMessage, err := goBot.Send(msg)
			if err != nil {
				log.Panic(err)
			}

			// Store the message ID for future edits
			messageId = sentMessage.MessageID

			// Handle button clicks
		} else if update.CallbackQuery != nil {
			callback := update.CallbackQuery.Data

			switch callback {
			case "prev":
				index--

				// Handle previous button click
				// Implement logic for going back to the previous content
				editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, messageId, messages[index])
				markup := createInlineKeyboard()
				editMsg.ReplyMarkup = &markup
				goBot.Send(editMsg)

			case "next":
				index++

				// Handle next button click
				// Implement logic for showing the next content
				editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, messageId, messages[index])
				markup := createInlineKeyboard()
				editMsg.ReplyMarkup = &markup
				goBot.Send(editMsg)
			}
		}

	}
}

func createInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	// Create an inline keyboard with prev and next buttons
	prevButton := tgbotapi.NewInlineKeyboardButtonData("Previous", "prev")
	nextButton := tgbotapi.NewInlineKeyboardButtonData("Next", "next")
	row := []tgbotapi.InlineKeyboardButton{prevButton, nextButton}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	return keyboard
}
