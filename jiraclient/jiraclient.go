package jiraclient

import (
	"errors"
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
	SetFixVersion
)

var commands []commandrouter.Command = []commandrouter.Command{
	{Comment, "comment", false, false, "Add a comment to an issue -- comment <id> <comment text>"},
	{Comment, "c", false, false, "Add a comment to an issue -- c <id> <comment text>"},
	{Assign, "assign", false, false, "Assign an issue to a user -- assign <id> <user>"},
	{Assign, "a", false, false, "Assign an issue to a user -- a <id> <user>"},
	// {Resolve, "resolve", false, false, "Resolve an issue -- resolve <id>"},
	// {Resolve, "r", false, false, "Resolve an issue -- r <id>"},
	// {Close, "close", false, false, "Close an issue -- close <id>"},
	// {Close, "c", false, false, "Close an issue -- c <id>"},
	// {Reopen, "reopen", false, false, "Reopen an issue -- reopen <id>"},
	// {SetFixVersion, "fixversion", false, false, "Set an issue's fix version -- fixversion <id> <version>"},
}

type parsedCommand struct {
	command string
	issue   string
	rest    string
}

func parseCommon(s string) (parsedCommand, error) {
	var issue string
	var rest string

	words := strings.SplitN(s, " ", 3)

	if len(words) > 1 {
		issue = strings.ToUpper(words[1])
		if strings.IndexAny(issue, "/#?&") != -1 {
			return parsedCommand{}, errors.New("Invalid character in issue ID")
		}
	}

	if len(words) > 2 {
		rest = words[2]
	}

	return parsedCommand{words[0], issue, rest}, nil
}

type JiraCommands struct {
	Client *jira.Client
}

func (jc *JiraCommands) AddComment(fromUser string, cmd parsedCommand) string {
	if cmd.rest == "" {
		return "No comment text provided"
	}

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
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "Issue not found"
		default:
			// TODO Log here
			return "Jira server internal error, see logs"
		}
	}

	return fmt.Sprintf("Added comment to issue %s", cmd.issue)
}

func (jc *JiraCommands) Assign(fromUser string, cmd parsedCommand) string {
	// TODO Look up our database of Jira/Slack users once it exists.
	username := cmd.rest
	if username == "" {
		return "No assignee provided. Use `none` to remove the assignee"
	} else if username == "none" {
		username = ""
	}

	url := fmt.Sprintf("/rest/api/2/issue/%s/assignee", cmd.issue)
	assignee := jira.Assignee{
		Name: username,
	}
	req, err := jc.Client.NewRequest("PUT", url, &assignee)
	if err != nil {
		// TODO log here
		return "Internal error, see logs"
	}

	resp, err := jc.Client.Do(req, nil)
	if err != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "Issue or user not found"
		default:
			// TODO Log here
			return "Jira server internal error, see logs"
		}
	}

	if username == "" {
		return fmt.Sprintf("Removed assignee from %s", cmd.issue)
	} else {
		return fmt.Sprintf("Assigned %s to %s", cmd.issue, username)
	}
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
	parsed, err := parseCommon(message.Text)
	if err != nil {
		return err.Error()
	}

	if parsed.issue == "" {
		return "Missing issue id"
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

	//httpClient := NewOAuthClient(${1:key}, ${2:accessToken}, ${3:accessSecret}, config.JiraUrl)
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
