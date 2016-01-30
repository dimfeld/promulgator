// This file contains the models that are passed between different
// servers running within the bot.

package model

import (
	"github.com/nlopes/slack"
)

type Config struct {
	Verbose bool

	// Base URL for Jira
	JiraUrl     string `required:"true"`
	JiraKey     string //`required:"true"`
	JiraAppName string `default:"JiraSlack"`

	// Jira OAuth access data
	JiraAccessToken  string //`required:"true"`
	JiraAccessSecret string //`required:"true"`

	SlackKey            string `required:"true"`
	SlackUser           string `default: "jira"`
	SlackDefaultChannel string `required:"true"`

	WebHookBind string `default:":80"`
	WebHookKey  string `required:"true"`
}

// ChatMessage contains information about an incoming or outgoing message.
type ChatMessage struct {
	// For incoming chat messages, the user who typed the message. This is
	// ignored for outging messages since we can't pretend to be another user.
	FromUser string
	// The user to whom the message is addresesed, if any. For Slack, this is
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
