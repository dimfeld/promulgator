package jirawebhook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"

	"github.com/dimfeld/promulgator/model"
)

var replacer *strings.Replacer

func IssueCreatedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}

func IssueDeletedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}

func IssueUpdatedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	if data.Comment == nil || data.Comment.Created != data.Comment.Updated {
		// Right now we only care about new comments.
		return nil, nil
	}

	link := fmt.Sprintf("%sbrowse/%s", data.JiraURL, data.Issue.Key)

	attachment := slack.Attachment{
		Fallback: fmt.Sprintf("%s commented on %s: %s",
			data.User.DisplayName, data.Issue.Key, data.Comment.Body),
		Pretext: fmt.Sprintf("%s commented on %s <%s|%s>",
			data.User.DisplayName, data.Issue.Fields.Type.Name, link, data.Issue.Key),
		Text:       replacer.Replace(data.Comment.Body),
		MarkdownIn: []string{"text"},
	}

	msg := &model.ChatMessage{
		Attachments: []slack.Attachment{attachment},
	}
	return msg, nil
}
