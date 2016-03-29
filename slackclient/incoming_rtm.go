package slackclient

import (
	"sync"

	"github.com/nlopes/slack"
	"github.com/spacemonkeygo/spacelog"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

var rtmLogger *spacelog.Logger

func routeMessage(api *slack.Client, commandrouter *commandrouter.Router, msg *slack.MessageEvent) {

}

func StartIncomingRTM(wg *sync.WaitGroup, config *model.Config, api *slack.Client,
	commandrouter *commandrouter.Router, done chan struct{}) {

	rtmLogger = spacelog.GetLoggerNamed("slack-rtm")

	// wg.Add(1)
	//
	// go (func() {
	// 	rtm := api.NewRTM()
	// 	go rtm.ManageConnection()
	//
	// 	for {
	// 		select {
	// 		case msg := <-rtm.IncomingEvents:
	// 			if ev, ok := msg.Data.(*slack.MessageEvent); ok {
	// 				// Handle the message
	// 				go routeMessage(api, commandrouter, ev)
	// 			}
	// 		case <-done:
	// 			err := rtm.Disconnect()
	// 			if err != nil {
	// 				rtmLogger.Errorf("Shutting down Slack RTM: %s", err.Error())
	// 			}
	// 			wg.Done()
	// 			return
	// 		}
	// 	}
	// })()

}
