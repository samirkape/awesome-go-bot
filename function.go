package mybot

import (
	"encoding/json"
	"errors"
	"github.com/samirkape/awesome-go-bot/internal/bot"
	"github.com/samirkape/awesome-go-bot/internal/commands"
	"github.com/samirkape/awesome-go-bot/internal/repository"
	"log"
	"net/http"
)

// HandleTelegramWebHook is a Google cloud function that handles all the requests made by the user.
// Telegram asynchronously does HTTP POST request on this function.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parsing POST request made by telegram on behalf of user
	response, err := responseDecoder(w, r)
	if err != nil {
		log.Printf("requesthandler: something went wrong when parsing the response %v", err)
		return
	}

	AllData := repository.GetAllPackages()
	commands.Execute(response, AllData)
}

// A responseDecoder() parses JSON response from the POST request.
// if all goes right, user id, message string and error value are returned.
func responseDecoder(w http.ResponseWriter, r *http.Request) (*bot.Request, error) {
	var request bot.RequestData
	var inputQuery bot.Request

	if r.Method == http.MethodPost {
		log.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = r.Body.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("responseDecoder: invalid request method")
	}

	// validate chat
	if request.Message.Chat.ID > 0 {
		inputQuery.ChatID = request.Message.Chat.ID
		inputQuery.Command = request.Message.Text
	} else {
		return nil, errors.New("responseDecoder: invalid user chat-id")
	}

	return &inputQuery, nil
}
