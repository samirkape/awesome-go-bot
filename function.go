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
	response, err := responseDecoder(w, r)
	if err != nil {
		log.Printf("requesthandler: something went wrong when parsing the response %v", err)
		return
	}

	// Head package list from the databse
	AllData := GetAllData()

	// Handle command given in the msgText
	// e.g /listpackages, /getStats
	ExecuteCommand(response, AllData)
}

// A responseDecoder() parses JSON response from the POST request.
// if all goes right, user id, message string and error value are returned.
func responseDecoder(w http.ResponseWriter, r *http.Request) (*BotResponse, error) {
	var message ReceiveMessage
	var response BotResponse

	// Parse incoming request
	if r.Method == http.MethodPost {
		log.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		r.Body.Close()
	}

	// Validate chat id
	if message.Message.Chat.ID > 0 {
		log.Println(message.Message.Chat.ID, message.Message.Text)
		response.ChatID = message.Message.Chat.ID
		response.MsgText = message.Message.Text
	} else {
		return nil, errors.New("responseDecoder: invalid user chat id")
	}

	return &response, nil
}
