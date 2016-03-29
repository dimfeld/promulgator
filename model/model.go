// This file contains the models that are passed between different
// servers running within the bot.

package model

import (
	"github.com/nlopes/slack"
)

type Config struct {
	Verbose bool

	// Base URL for Jira. This should include basic auth information.
	JiraUrl string `envconfig:"JIRA_URL" required:"true"`
	// JiraApiKey string `envconfig:"JIRA_API_KEY"` //`required:"true"`
	// JiraAppName    string `envconfig:"JIRA_APPNAME" default:"Promulgator"`
	JiraWebHookKey string `envconfig:"JIRA_WEBHOOK_KEY" required:"true"`

	// Jira OAuth access data
	//JiraAccessToken  string `envconfig:"JIRA_ACCESS_TOKEN"`  //`required:"true"`
	//JiraAccessSecret string `envconfig:"JIRA_ACCESS_SECRET"` //`required:"true"`

	// The token for the bot configuration.
	SlackKey string `envconfig:"SLACK_KEY" required:"true"`
	// The token for the Slack slash command. If empty, this gets the same value
	// as SlackKey, which is the desired behavior when using a full "app" as opposed
	// to team-specific bots and slash commands.
	SlackSlashCommandKey string `envconfig:"SLACK_SLASH_COMMAND_KEY"`
	// The name the bot should post as.
	SlackUser string `envconfig:"SLACK_USER" default:"jira"`
	// The Icon to use
	SlackIcon string `envconfig:"SLACK_ICON" default:"https://slack.global.ssl.fastly.net/12d4/img/services/jira_48.png"`
	// The channel to post messages to, when not invoked via chatbot DM.
	SlackDefaultChannel string `envconfig:"SLACK_DEFAULT_CHANNEL" required:"true"`

	// Timeout for HTTP requests, in milliseconds.
	RequestTimeout int `envconfig:"REQUEST_TIMEOUT" default:"30000"`

	// Listen on this IP/Port for webhooks.
	WebHookBind string `envconfig:"WEBHOOK_BIND" default:":80"`
}

// ChatMessage contains information about an incoming or outgoing message.
type ChatMessage struct {
	// For incoming chat messages, the user who typed the message. This is
	// ignored for outging messages since we can't pretend to be another user.
	FromUser string
	// The user to whom the message is addressed, if any. For Slack, this is
	// the @user at the beginning of the message.
	ToUser string
	// The name of a channel to send to.
	Channel string
	Text    string
	// For incoming messages, true if the message is addressed to the bot.
	// This can either be through a DM channel or by starting with @botuser.
	ToBot bool
	// Attachments, if any
	Attachments []slack.Attachment
}
