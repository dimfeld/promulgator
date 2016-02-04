package slackclient

import (
	"github.com/nlopes/slack"
	"sync"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.ChatMessage,
	commandrouter *commandrouter.Router, done chan struct{}) {

	api := slack.New(config.SlackKey)

	StartOutgoing(wg, config, api, inChan, done)
}
