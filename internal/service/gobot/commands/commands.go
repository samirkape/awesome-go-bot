package commands

import (
	"strconv"
	"strings"
)

type Commands struct {
	Start          string
	ListCategories string
	GetPackages    string
	TopN           string
	Description    string
}

func New() *Commands {
	return &Commands{
		Start:          "/start",
		ListCategories: "/list_categories",
		GetPackages:    "/select_category",
		TopN:           "/top",
		Description:    "/description",
	}
}

func (c *Commands) GetStart() string {
	return c.Start
}

func (c *Commands) GetListCategories() string {
	return c.ListCategories
}

func (c *Commands) GetGetPackages() string {
	return c.GetPackages
}

func (c *Commands) IsTopN(query string) string {
	if strings.HasPrefix(query, c.TopN) {
		return query
	}
	return ""
}

func (c *Commands) IsCategoryNumber(query string) string {
	_, err := strconv.Atoi(query)
	if err != nil {
		return ""
	}
	return query
}

func (c *Commands) GetDescription() string {
	return c.Description
}
