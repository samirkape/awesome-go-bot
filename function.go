package mybot

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// HandleTelegramWebHook parses a POST request from telegram and responds with appropriate actions.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse POST request
	// chatID: is a unique user identifier.
	// msgText: is a message/command user has sent to the bot.
	chatID, msgText, err := requestHandler(w, r)
	if err != nil {
		log.Printf("requesthandler: unable to proceed %v", err)
		return
	}

	// read package list from the databse
	allPackages := ListCategories()

	// handle command given in the msgText
	// e.g /listpackages, /getStats
	executeCommand(msgText, chatID, allPackages)
}

func requestHandler(w http.ResponseWriter, r *http.Request) (int, string, error) {
	var message ReceiveMessage
	chatID := 0
	msgText := ""

	// Parse incoming request
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			log.Println(err)
			return -1, "", err
		}
		r.Body.Close()
	}

	// validate chat id
	if message.Message.Chat.ID > 0 {
		log.Println(message.Message.Chat.ID, message.Message.Text)
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	} else {
		return 0, "", errors.New("invalid user chat id")
	}

	return chatID, msgText, nil
}
