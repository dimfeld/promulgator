package slackclient

import (
	"sync"

	"github.com/nlopes/slack"
	"github.com/spacemonkeygo/spacelog"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

var logIn *spacelog.Logger
var logOut *spacelog.Logger

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.ChatMessage,
	commandrouter *commandrouter.Router, done chan struct{}) {

	logIn = spacelog.GetLoggerNamed("slack-incoming")
	logOut = spacelog.GetLoggerNamed("slack-outgoing")

	api := slack.New(config.SlackKey)

	StartOutgoing(wg, config, api, inChan, done)

	// Not ready yet
	//StartIncomingWebhook(config, commandrouter)
}
