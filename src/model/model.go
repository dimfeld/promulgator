// This file contains the models that are passed between different
// servers running within the bot.

package model

type Config struct {
	Verbose bool

	// Base URL for Jira
	JiraUrl     string `required:"true"`
	JiraKey     string `required:"true"`
	JiraAppName string `default:"JiraSlack"`

	// Jira OAuth access data
	JiraAccessToken  string `required:"true"`
	JiraAccessSecret string `required:"true"`

	SlackKey string `required:"true"`
	SlackUrl string `default:""`
}

type User struct {
	Id   string
	Name string
}

type Comment struct {
	User     User
	Contents string
}

type Issue struct {
	Id          string
	Assignee    User
	Description string
	Status      string
}

type UpdateType int

const (
	UpdateTypeIssue UpdateType = iota
	UpdateTypeComment
)

type TrackerUpdate struct {
	UpdateType      UpdateType
	UpdatedObjectId string
	Edits           map[string]string
}

type ChatMessage struct {
	User     string
	Title    string
	Contents string
}
