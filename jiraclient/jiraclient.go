package jiraclient

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

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
	FixVersion
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
	{FixVersion, "fixversion", false, false, "Set an issue's fix version -- fixversion <id> <version>"},
}

type parsedCommand struct {
	command string
	issue   string
	rest    string
}

func parseCommon(input string) (parsedCommand, error) {
	var command string
	var issue string
	var rest string

	nextWord := func(s string) (firstWordEnd int, secondWordBegin int) {
		i := 0
		// Skip past the current word
		for i < len(s) && s[i] != ' ' {
			i++
		}
		firstWordEnd = i

		// Past all spaces, to the next word
		for i < len(s) && s[i] == ' ' {
			i++
		}
		secondWordBegin = i
		return
	}

	firstEnd, secondBegin := nextWord(input)
	command = input[:firstEnd]

	if secondBegin < len(input) {
		input = input[secondBegin:]
		secondEnd, restBegin := nextWord(input)
		issue = input[0:secondEnd]

		if restBegin < len(input) {
			rest = input[restBegin:]
		}

		if strings.IndexAny(issue, "/#?&") != -1 {
			return parsedCommand{}, errors.New("Invalid character in issue ID")
		}
	}

	return parsedCommand{command, issue, rest}, nil
}

func errorDetail(resp *http.Response, err error, key string) string {
	if err != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return "Issue not found"
		case http.StatusBadRequest:
			if errDetail, ok := err.(*jira.ErrorResponse); ok {
				if e, ok := errDetail.Errors[key]; ok {
					return e
				}
			}
		default:
			// TODO Log here
			return "Jira server internal error, see logs"
		}
	}

	return "Internal error, see logs"
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
	_, resp, err := jc.Client.Issue.AddComment(cmd.issue, &comment)
	if err != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return fmt.Sprintf("Issue %s not found", cmd.issue)
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
		return errorDetail(resp, err, "assignee")
	}

	if username == "" {
		return fmt.Sprintf("Removed assignee from %s", cmd.issue)
	} else {
		return fmt.Sprintf("Assigned %s to %s", cmd.issue, username)
	}
}

// TODO Putting these off for now since we have to navigate the transition
// system and this is also specific to both the Jira config and the issue type.
// Likely, this will take some extra configuration to map short command names to
// longer transition names, or something along those lines.
func (jc *JiraCommands) Resolve(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Close(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) Reopen(fromUser string, cmd parsedCommand) string {
	return "Unimplemented!"
}

func (jc *JiraCommands) SetFixVersion(fromUser string, cmd parsedCommand) string {
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			FixVersions: []*jira.FixVersion{
				&jira.FixVersion{Name: cmd.rest},
			},
		},
	}

	url := fmt.Sprintf("/rest/api/2/issue/%s", cmd.issue)
	req, err := jc.Client.NewRequest("PUT", url, &issue)
	if err != nil {
		// TODO log here
		return "Internal error, see logs"
	}

	resp, err := jc.Client.Do(req, nil)
	if err != nil {
		return errorDetail(resp, err, "fixVersions")
	}

	return fmt.Sprintf("Set fix version of %s to %s", cmd.issue, cmd.rest)
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
	case FixVersion:
		return jc.SetFixVersion(message.FromUser, parsed)
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

	httpClient := http.Client{}
	httpClient.Timeout = time.Duration(config.RequestTimeout) * time.Millisecond
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
