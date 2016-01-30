package slackclient

import (
	"commandrouter"
	"github.com/nlopes/slack"
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.ChatMessage,
	commandrouter *commandrouter.Router, done chan struct{}) {

	api := slack.New(config.SlackKey)

	StartOutgoing(wg, config, api, inChan, done)
}
