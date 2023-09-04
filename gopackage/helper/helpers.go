package helper

import (
	"awesome-go-bot-refactored/gopackage"
	"fmt"
	"strings"
)

// ListToMessage Converts slice of strings into a single string.
func ListToMessage(list []gopackage.CategoryName) string {
	var msg strings.Builder
	for i, pkg := range list {
		msg.WriteString(fmt.Sprintf("%d. %s\n", i, pkg))
	}
	return msg.String()
}
