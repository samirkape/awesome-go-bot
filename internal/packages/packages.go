package packages

import (
	"fmt"
	"strings"
)

// PackageToMsg merges single Package struct elements into a message string.
func (input Package) PackageToMsg(forTop bool) string {
	msgString := strings.Builder{}
	name := strings.Title(strings.ToLower(input.Name))
	if forTop {
		msgString.WriteString(fmt.Sprintf("[%s](%s)\nStars: %d\nCategory: %s%s\n", name, input.URL, input.Stars, input.Title, input.Info))
	} else {
		msgString.WriteString(fmt.Sprintf("[%s](%s)\nStars: %d\n%s\n", name, input.URL, input.Stars, input.Info))
	}
	return msgString.String()
}

// PackagesToMsg The packagesToList method works on len(receiver)
// and merge them together
func (input Packages) PackagesToMsg(forTop bool) string {
	msg := strings.Builder{}
	for _, pkg := range input {
		msg.WriteString(pkg.PackageToMsg(forTop))
		msg.WriteString("\n")
	}
	return msg.String()
}
