package slackclient

import (
	"golang.org/x/net/context"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

func processIncomingSlack(ctx context.Context, commandrouter *commandrouter.Router,
	command string, user string, inChannel string, toBot bool, responseChan chan string) {

	msg := model.ChatMessage{
		FromUser: user,
		Channel:  inChannel,
		Text:     command,
		ToBot:    toBot,
	}

	handled, err := commandrouter.Route(ctx, &msg, responseChan)
	if err != nil {
		logIn.Errorf("commandrouter: %s", err.Error())
		// TODO Make this message configurable?
		select {
		case responseChan <- "Internal error, please see bot logs":
		default:
		}
	} else if !handled {
		responseChan <- "I didn't recognize that command"
	}
}
