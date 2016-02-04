package slackclient

import (
	"github.com/nlopes/slack"
	"sync"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

func StartIncoming(wg *sync.WaitGroup, config *model.Config, api *slack.Client,
	commandrouter *commandrouter.Router,
	done chan struct{}) {

}
