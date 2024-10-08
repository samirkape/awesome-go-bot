package commands

import (
	"strconv"
	"strings"
)

type Commands struct {
	Start          string
	Commands       string
	ListCategories string
	GetPackages    string
	TopN           string
	Description    string
}

func New() *Commands {
	return &Commands{
		Start:          "/start",
		Commands:       "/commands",
		ListCategories: "/listcategories",
		GetPackages:    "/selectcategory",
		TopN:           "/top",
		Description:    "/description",
	}
}

func (c *Commands) GetStart() string {
	return c.Start
}

func (c *Commands) GetSupportedCommands() string {
	return c.Commands
}

func (c *Commands) GetListCategories() string {
	return c.ListCategories
}

func (c *Commands) GetPackages() string {
	return c.GetPackages
}

func (c *Commands) IsTopN(query string) string {
	if strings.HasPrefix(query, c.TopN) {
		return query
	}
	return "nope"
}

func (c *Commands) IsTop(query string) string {
	if strings.HasPrefix(query, c.Top) { // TODO
		return query
	}
	return "nope"
}

func (c *Commands) IsCategoryNumber(query string) string {
	var newQuery string
	if strings.HasPrefix(query, "/") {
		newQuery = strings.TrimPrefix(query, "/")
	}
	_, err := strconv.Atoi(newQuery)
	if err != nil {
		return "nope"
	}
	return query
}

func (c *Commands) GetDescription() string {
	return c.Description
}
