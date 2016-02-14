package slackclient

import (
	"sync"

	"github.com/nlopes/slack"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

func StartIncomingRTM(wg *sync.WaitGroup, config *model.Config, api *slack.Client,
	commandrouter *commandrouter.Router, done chan struct{}) {

}
