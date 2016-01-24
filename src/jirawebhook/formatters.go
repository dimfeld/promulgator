package jirawebhook

import (
	"errors"
	"model"
)

func IssueCreatedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}

func IssueDeletedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}

func IssueUpdatedFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}

func CommentFormatter(data *JiraWebhook) (*model.ChatMessage, error) {
	return nil, errors.New("Not implemented")
}
