package util

import (
	"errors"
	"fmt"
	"strings"
)

// ListToMsg Converts slice of strings into a single string.
func ListToMsg(list []string) string {
	msg := strings.Builder{}
	for i, pkg := range list {
		msg.WriteString(fmt.Sprint(i) + ". " + pkg + string("\n")) // 3 = remove ## from start
	}
	return msg.String()
}

// Check for unhandled command and invalid index number
func ValidateMessage(msgText string) error {
	// Input validation: Check if it is a unhandled command
	if strings.HasPrefix(msgText, "/") {
		return errors.New("invalid response, try numeric input")
	}
	return nil
}
