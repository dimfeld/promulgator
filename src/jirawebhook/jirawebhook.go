package jirawebhook

import (
	"encoding/json"
	"github.com/Carevoyance/go-jira"
	"model"
	"net/http"
	"sync"
)

type JiraChangeItem struct {
	ToString   string `json:"toString"`
	To         string `json:"to"`
	FromString string `json:"fromString"`
	From       string `json:"from"`
	FieldType  string `json:"fieldtype"`
	Field      string `json:"field"`
}

type JiraChangelog struct {
	Items []JiraChangeItem `json:"items"`
	Id    int              `json:"id"`
}

type JiraWebhook struct {
	Timestamp uint64         `json:"timestamp"`
	Event     string         `json:"event"`
	User      *jira.Assignee `json:"user"`
	Issue     *jira.Issue    `json:"user"`
	Changelog *JiraChangelog `json:"changelog"`
	Comment   *jira.Comment  `json:"comment"`
}

type WebhookFormatter func(*JiraWebhook) (*model.ChatMessage, error)

var handlers map[string]WebhookFormatter

func handleWebhook(config *model.Config, outChan chan *model.ChatMessage,
	w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	if query.Get("key") != config.WebHookKey {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid API key"))
		return
	}

	d := json.NewDecoder(r.Body)
	var data *JiraWebhook
	err := d.Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON decode error: " + err.Error()))
		return
	}

	if handler, ok := handlers[data.Event]; ok {
		message, err := handler(data)
		if err != nil {
			// TODO Log something
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		outChan <- message
	}

	w.WriteHeader(http.StatusOK)
}

func Start(config *model.Config, wg *sync.WaitGroup,
	outChan chan *model.ChatMessage, done chan struct{}) {

	handlers = map[string]WebhookFormatter{
		"issue_updated":   IssueUpdatedFormatter,
		"issue_created":   IssueCreatedFormatter,
		"issue_deleted":   IssueDeletedFormatter,
		"comment_created": CommentFormatter,
	}

	go func() {
		// Create web server to listen for webhooks
		http.HandleFunc("/jirahook", func(w http.ResponseWriter, r *http.Request) {
			handleWebhook(config, outChan, w, r)
		})
		// TODO Support TLS. Unimportant for now only because this runs solely
		// within our own network.
		err := http.ListenAndServe(config.WebHookBind, nil)
		if err != nil {
			panic(err.Error) // TODO Fatal error, but not panic
		}
	}()
}
