package internalerrors

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samirkape/awesome-go-bot/gobot"
	"github.com/samirkape/awesome-go-bot/internal/services/chat"
)

type ValidationError struct {
	message     string
	Type        string
	placeHolder interface{}
}

func (v ValidationError) Error() string {
	if v.placeHolder == nil {
		return v.message
	}
	return fmt.Sprintf(v.message+": %+v", v.placeHolder)
}

func NewValidationError(message string, placeHolder ...interface{}) ValidationError {
	return ValidationError{
		message:     message,
		placeHolder: placeHolder,
	}
}

func RespondToError(err error, botService *tgbotapi.BotAPI, chatService chat.Info) error {
	if err == nil {
		return nil
	}
	var validationError ValidationError
	if errors.As(err, &validationError) {
		err := gobot.Respond(chatService, botService, validationError.Error())
		if err != nil {
			return err
		}
	}
	return err
}
