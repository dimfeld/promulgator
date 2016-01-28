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

	WebHookBind string `default:":80"`
	WebHookKey  string `required:"true"`

	TemplateDir string `default:"./templates"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type ChatAttachment struct {
	Fallback   string            `json:"fallback,omitempty"`
	Color      string            `json:"color,omitempty"`
	Pretext    string            `json:"pretext,omitempty"`
	AuthorName string            `json:"author_name,omitempty"`
	AuthorLink string            `json:"author_link,omitempty"`
	AuthorIcon string            `json:"author_icon,omitempty"`
	Title      string            `json:"title,omitempty"`
	TitleLink  string            `json:"title_link,omitempty"`
	Text       string            `json:"text,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	ImageUrl   string            `json:"image_url,omitempty"`
	ThumbUrl   string            `json:"thumb_url,omitempty"`
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
	// Attachment, if any
	Attachment *ChatAttachment
}
