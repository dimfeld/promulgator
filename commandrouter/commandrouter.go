package commandrouter

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/net/context"

	"github.com/dimfeld/promulgator/model"
)

// Command defines a single command, how it should be detected, and what value
// should be used to notify the client.
type Command struct {
	// Tag is the value sent to the destination when a command matches.
	Tag int
	// The command to match on. If the command is not a regexp (see below), the
	// command must be a single word and will only be matched if it occurs at
	// the beginning of the message.
	Command string
	// IsRegExp is true if Command should be compiled as a regex.
	IsRegexp bool
	// MatchAll is true if the command should match on all messages. If false,
	// the command will match only on messages directed to the bot user.
	MatchAll bool
	Help     string
}

// cmd is internal data needed to route a command match
type cmd struct {
	Tag         int
	destination destination
	MatchAll    bool
}

type commandList []cmd

type regexpCommand struct {
	regexp *regexp.Regexp
	cmd    cmd
}

type destination struct {
	Channel chan Match
}

type destinationList []destination

// Match contains all the information needed for a client to process an incoming
// message.
type Match struct {
	// The tag supplied when the command was added.
	Tag     int
	Message *model.ChatMessage
	// RegexpMatch contains the matching subgroups, if any, when a Regexp
	// command matches.
	RegexpMatch []string
	Response    chan string
}

// Router is the command router itself.
type Router struct {
	destinations   []destination
	commands       map[string]commandList
	regexpCommands []regexpCommand

	done chan struct{}
}

// New creates a new Router
func New() *Router {
	return &Router{
		destinations:   destinationList{},
		commands:       map[string]commandList{},
		regexpCommands: []regexpCommand{},
		done:           make(chan struct{}),
	}
}

// Close shuts down the router and closes all it outgoing channels.
func (r *Router) Close() {
	for _, d := range r.destinations {
		close(d.Channel)
	}

	close(r.done)
	r.done = nil
}

func (r *Router) addCommand(d destination, c Command) error {
	cmd := cmd{Tag: c.Tag, destination: d, MatchAll: c.MatchAll}

	if c.IsRegexp {
		re, err := regexp.Compile(c.Command)
		if err != nil {
			return err
		}
		rc := regexpCommand{
			regexp: re,
			cmd:    cmd,
		}
		r.regexpCommands = append(r.regexpCommands, rc)
	} else {
		if strings.IndexAny(c.Command, " \r\n") != -1 {
			return errors.New("Commands must not have whitespace. Use a regex instead: " + c.Command)
		}

		cmdList := r.commands[c.Command]
		if cmdList == nil {
			cmdList = commandList{cmd}
		} else {
			cmdList = append(cmdList, cmd)
		}
		r.commands[c.Command] = cmdList
	}

	return nil
}

// AddDestination adds a new destination and associated commands to the router.
// name is an arbitrary string used to identify the destination.
func (r *Router) AddDestination(name string, commands []Command) (chan Match, error) {
	ci := make(chan Match, 1)

	dest := destination{ci}
	r.destinations = append(r.destinations, dest)

	for _, c := range commands {
		if err := r.addCommand(dest, c); err != nil {
			close(ci)
			return nil, err
		}
	}

	// Start a new buffering goroutine to handle sending to this destination.
	co := MatchBuffer(ci, r.done)

	return co, nil
}

// Route processes a ChatMessage and routes it to the correct destination, if any.
func (r *Router) Route(ctx context.Context, msg *model.ChatMessage, responseChan chan string) (bool, error) {
	if r.done == nil {
		return false, errors.New("Router is closed")
	}

	found := false
	firstSpace := strings.Index(msg.Text, " ")
	var firstWord string
	if firstSpace == -1 {
		// Just one word, but this might be ok
		firstWord = msg.Text
	} else {
		firstWord = msg.Text[:firstSpace]
	}
	if cmdList, ok := r.commands[firstWord]; ok {
		// There is a command for this word. Send it to all the destinations.
		for _, c := range cmdList {
			if c.MatchAll || msg.ToBot {
				match := Match{
					Tag:     c.Tag,
					Message: msg,
					// Send responseChan to all matches. This does mean that the
					// match handler should send from inside a select to avoid
					// blocking.
					// TODO It also means that the first response to be sent
					// is the only one processed, which isn't ideal. I should
					// come up with a better system that amalgamates all the
					// responses into one when there are multiple matches.
					Response: responseChan,
				}
				c.destination.Channel <- match
				found = true
			}
		}
	}

	for _, rc := range r.regexpCommands {
		if !rc.cmd.MatchAll && !msg.ToBot {
			// Don't even bother matching this regexp if it only applies to
			// messages addressed to us, and this message is not.
			continue
		}
		reMatch := rc.regexp.FindStringSubmatch(msg.Text)
		if reMatch != nil {
			match := Match{
				Tag:     rc.cmd.Tag,
				Message: msg,
			}
			rc.cmd.destination.Channel <- match
			found = true
		}
	}

	return found, nil
}
