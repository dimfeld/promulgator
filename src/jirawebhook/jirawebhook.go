package jirawebhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
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
	Id    string           `json:"id"`
}

type JiraWebhook struct {
	Timestamp uint64         `json:"timestamp"`
	Event     string         `json:"webhookEvent"`
	User      *jira.Assignee `json:"user"`
	Issue     *jira.Issue    `json:"issue"`
	Changelog *JiraChangelog `json:"changelog"`
	Comment   *jira.Comment  `json:"comment"`
	JiraURL   string         `json:"-"`
}

type WebhookFormatter func(*JiraWebhook) (*model.ChatMessage, error)

var handlers map[string]WebhookFormatter

func handleWebhook(config *model.Config, outChan chan *model.ChatMessage,
	w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	if query.Get("key") != config.WebHookKey {
		fmt.Printf("Bad webhook key %s\n", query.Get("key"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid API key"))
		return
	}

	buf, _ := ioutil.ReadAll(r.Body)

	d := json.NewDecoder(bytes.NewReader(buf))
	data := &JiraWebhook{}
	if err := d.Decode(data); err != nil {
		fmt.Println("JSON decode error: " + err.Error())
		fmt.Println(string(buf))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON decode error: " + err.Error()))
		return
	}

	// TODO Proper debug logging
	// buf, _ := json.MarshalIndent(data, "|", "  ")
	// fmt.Println(string(buf))
	if handler, ok := handlers[data.Event]; ok {
		// Set this since the formatters will probably want it for hyperlinks.
		data.JiraURL = config.JiraUrl
		message, err := handler(data)
		if err != nil {
			// TODO Log something
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if message != nil {
			outChan <- message
		}
	}
	w.WriteHeader(http.StatusOK)
}

func Start(config *model.Config, wg *sync.WaitGroup,
	outChan chan *model.ChatMessage, done chan struct{}) {

	wg.Add(1)

	handlers = map[string]WebhookFormatter{
		"jira:issue_updated": IssueUpdatedFormatter,
		// "jira:issue_created":   IssueCreatedFormatter,
		// "jira:issue_deleted":   IssueDeletedFormatter,
	}

	go func() {
		// Create web server to listen for webhooks
		http.HandleFunc("/jirahook", func(w http.ResponseWriter, r *http.Request) {
			handleWebhook(config, outChan, w, r)
		})
		// TODO Support TLS. Unimportant for now only because this runs solely
		// within our own network.
		fmt.Printf("Listening on %s\n", config.WebHookBind)
		err := http.ListenAndServe(config.WebHookBind, nil)
		if err != nil {
			panic(err.Error()) // TODO Fatal error, but not panic
		}
	}()
}
