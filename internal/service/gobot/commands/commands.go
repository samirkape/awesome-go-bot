package commands

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
		ListCategories: "/list-categories",
		GetPackages:    "/select-category",
		TopN:           "/top-n",
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

func (c *Commands) GetTopN() string {
	return c.TopN
}

func (c *Commands) GetDescription() string {
	return c.Description
}
