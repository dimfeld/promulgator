package jiraclient

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	jira "github.com/andygrunwald/go-jira"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

const (
	Comment = iota
	Assign
	Resolve
	Close
	Reopen
)

var commands []commandrouter.Command = []commandrouter.Command{
	{Comment, "comment", false, false, "Add a comment to an issue -- comment <tag> <comment text>"},
	{Comment, "c", false, false, "Add a comment to an issue -- c <tag> <comment text>"},
	{Assign, "assign", false, false, "Assign an issue to a user -- assign <tag> <user>"},
	{Assign, "a", false, false, "Assign an issue to a user -- a <tag> <user>"},
	{Resolve, "resolve", false, false, "Resolve an issue -- resolve <tag>"},
	{Resolve, "r", false, false, "Resolve an issue -- r <tag>"},
	{Close, "close", false, false, "Close an issue -- close <tag>"},
	{Close, "c", false, false, "Close an issue -- c <tag>"},
	{Reopen, "reopen", false, false, "Reopen an issue -- reopen <tag>"},
}

type parsedCommand struct {
	command string
	issue   string
	rest    string
}

func parseCommon(s string) parsedCommand {
	var issue string
	var rest string

	words := strings.SplitN(s, " ", 3)

	if len(words) > 1 {
		issue = words[1]
	}

	if len(words) > 2 {
		rest = words[2]
	}

	return parsedCommand{words[0], issue, rest}
}

type JiraCommands struct {
	Client *jira.Client
}

func (jc *JiraCommands) AddComment(fromUser string, cmd parsedCommand) string {
	comment := jira.Comment{
		Body: fmt.Sprintf("%s commented:\n%s", fromUser, cmd.rest),
	}
	url := fmt.Sprintf("/rest/api/2/issue/%s/comment", cmd.issue)
	req, err := jc.Client.NewRequest("POST", url, &comment)
	if err != nil {
		// TODO log here
		return "Internal error, see logs"
	}

	resp, err := jc.Client.Do(req, nil)
	if err != nil {
		// TODO log here
		return "Error contacting Jira server"
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return "Success"
	case http.StatusNotFound:
		return "Issue not found"
	default:
		// TODO Log here
		return "Jira server internal error, see logs"
	}
}

func (jc *JiraCommands) Assign(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Resolve(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Close(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Reopen(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Process(command int, message *model.ChatMessage) string {
	parsed := parseCommon(message.Text)
	if parsed.issue == "" {
		return "Missing issue tag"
	}

	switch command {
	case Comment:
		return jc.AddComment(message.FromUser, parsed)
	case Assign:
		return jc.Assign(message.FromUser, parsed)
	case Resolve:
		return jc.Resolve(message.FromUser, parsed)
	case Close:
		return jc.Close(message.FromUser, parsed)
	case Reopen:
		return jc.Reopen(message.FromUser, parsed)
	default:
		// TODO Log about unexpected command tag
		return "Internal error"
	}
}

// Start the Jira client goroutine
func Start(config *model.Config, wg *sync.WaitGroup,
	commandrouter *commandrouter.Router,
	outChan chan *model.ChatMessage, done chan struct{}) {

	wg.Add(1)

	jiraClient, err := jira.NewClient(nil, config.JiraUrl)
	if err != nil {
		panic(err)
	}

	processor := &JiraCommands{
		Client: jiraClient,
	}

	incoming, err := commandrouter.AddDestination("JiraCommands", commands)
	if err != nil {
		panic(err)
	}

	go func() {
	Loop:
		for {
			select {
			case match := <-incoming:
				match.Response <- processor.Process(match.Tag, match.Message)
			case <-done:
				break Loop
			}
		}
		wg.Done()
	}()

}
