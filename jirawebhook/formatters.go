package jirawebhook

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"

	"github.com/dimfeld/promulgator/model"
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

	link := fmt.Sprintf("%s/browse/%s", data.JiraURL, data.Issue.Key)

	attachment := slack.Attachment{
		Pretext: fmt.Sprintf("%s commented on %s <%s|%s>",
			data.User.DisplayName, data.Issue.Fields.Type.Name, link, data.Issue.Key),
		Text:       data.Comment.Body,
		MarkdownIn: []string{"text"},
	}

	msg := &model.ChatMessage{
		Attachments: []slack.Attachment{attachment},
	}
	return msg, nil
}
