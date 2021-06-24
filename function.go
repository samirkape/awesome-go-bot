package mybot

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Telegram asynchronusly does HTTP POST request on the trigger URL we have given,
// (the one we get from Google cloud functions).
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse POST request
	// chatID: is a unique user identifier.
	// msgText: is a message/command user has sent to the bot.
	chatID, msgText, err := responseDecoder(w, r)
	if err != nil {
		log.Printf("requesthandler: something went wrong when parsing the response %v", err)
		return
	}

	// Head package list from the databse
	allPackages := ListCategories()

	// Handle command given in the msgText
	// e.g /listpackages, /getStats
	executeCommand(msgText, chatID, allPackages)
}

// A responseDecoder() parses JSON response from the POST request.
// if all goes right, then user id, message string and error value are returned.
func responseDecoder(w http.ResponseWriter, r *http.Request) (int, string, error) {
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

	// Validate chat id
	if message.Message.Chat.ID > 0 {
		log.Println(message.Message.Chat.ID, message.Message.Text)
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	} else {
		return 0, "", errors.New("responseDecoder: invalid user chat id")
	}

	return chatID, msgText, nil
}
