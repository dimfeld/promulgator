package jirawebhook

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"model"
	// "path/filepath"
)

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

	attachment := slack.Attachment{
		Title:     fmt.Sprintf("%s commented on %s", data.User.Name, data.Issue.Key),
		TitleLink: data.Comment.Self,
		Text:      data.Comment.Body,
	}

	msg := &model.ChatMessage{
		Attachments: []slack.Attachment{attachment},
	}
	return msg, nil
}
