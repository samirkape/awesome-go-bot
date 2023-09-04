package main

import (
	"awesome-go-bot/gopackage/mongodb"
	"awesome-go-bot/gopackage/search"
	"awesome-go-bot/internal/service/chat/factory"
	"awesome-go-bot/internal/service/chat/inline"
	"awesome-go-bot/internal/service/chat/keyboard"
	"awesome-go-bot/internal/service/chat/regular"
	"awesome-go-bot/internal/service/gobot"
	"awesome-go-bot/internal/service/gobot/config"
	"context"
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
	if err != nil {
		return
	}

	packageService, err := client.GetAllPackages()
	searchService := search.NewService(packageService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600
	updates := botService.GetUpdatesChan(u)

	for update := range updates {
		chat := factory.New(context.Background(), &update)
		if chat.IsInline() {
			inline.HandleQuery(botService, searchService, chat)
		} else if chat.IsCallBack() {
			keyboard.ProcessUsingInlineKeyboard(botService, packageService, chat)

		} else {
			regular.HandleQuery(botService, packageService, chat)
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
