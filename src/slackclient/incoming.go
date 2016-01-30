package slackclient

import (
	"commandrouter"
	"github.com/nlopes/slack"
	"model"
	"sync"
)

func StartIncoming(wg *sync.WaitGroup, config *model.Config, api *slack.Client,
	commandrouter *commandrouter.Router,
	done chan struct{}) {

}
