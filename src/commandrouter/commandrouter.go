package commandrouter

import (
	"errors"
	"model"
	"regexp"
	"strings"
)

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
}

// cmd is internal data needed to route a command match
type cmd struct {
	Tag         int
	Destination Destination
	MatchAll    bool
}

type CommandList []cmd

type regexpCommand struct {
	regexp *regexp.Regexp
	cmd    cmd
}

type Destination struct {
	Channel chan Match
}

type DestinationList []Destination

type Match struct {
	// The tag supplied when the command was added.
	Tag     int
	Message *model.ChatMessage
	// RegexpMatch contains the matching subgroups, if any, when a Regexp
	// command matches.
	RegexpMatch []string
}

type Router struct {
	destinations   []Destination
	commands       map[string]CommandList
	regexpCommands []regexpCommand

	done chan struct{}
}

func New() *Router {
	return &Router{
		destinations:   DestinationList{},
		commands:       map[string]CommandList{},
		regexpCommands: []regexpCommand{},
		done:           make(chan struct{}),
	}
}

func (r *Router) Close() {
	close(r.done)
}

func (r *Router) addCommand(d Destination, c Command) error {
	cmd := cmd{Tag: c.Tag, Destination: d, MatchAll: c.MatchAll}

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
			cmdList = CommandList{cmd}
		} else {
			cmdList = append(cmdList, cmd)
		}
		r.commands[c.Command] = cmdList
	}

	return nil
}

func (r *Router) AddDestination(name string, commands []Command) (chan Match, error) {
	ci := make(chan Match, 1)

	dest := Destination{ci}
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
func (r *Router) Route(msg *model.ChatMessage) (bool, error) {
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
				}
				c.Destination.Channel <- match
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
			rc.cmd.Destination.Channel <- match
			found = true
		}
	}

	return found, nil
}
