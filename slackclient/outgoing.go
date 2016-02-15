package slackclient

import (
	"sync"

	"github.com/nlopes/slack"

	"github.com/dimfeld/promulgator/model"
)

func sendMessage(sendWg *sync.WaitGroup, config *model.Config, api *slack.Client,
	msg *model.ChatMessage) {

	params := slack.NewPostMessageParameters()
	params.Attachments = msg.Attachments
	params.Username = config.SlackUser
	params.IconURL = config.SlackIcon
	params.LinkNames = 1
	params.AsUser = false

	channel := msg.Channel
	if channel == "" {
		// For now just print everything to the main channel. When we have RTM
		// support up and running we might also want the ability to write back
		// to a user DM channel.
		channel = config.SlackDefaultChannel
	}

	logOut.Debugf("Posting to channel %s %s %+v", channel, msg.Text, params)
	_, _, err := api.PostMessage(channel, msg.Text, params)
	if err != nil {
		logOut.Errorf("Error writing chat message: %s\n", err.Error())
	}
	sendWg.Done()
}

func StartOutgoing(wg *sync.WaitGroup, config *model.Config,
	api *slack.Client,
	inChan chan *model.ChatMessage,
	done chan struct{}) {

	wg.Add(1)

	go func() {
		sendWg := &sync.WaitGroup{}
		for {
			select {
			case msg := <-inChan:
				sendWg.Add(1)
				go sendMessage(sendWg, config, api, msg)

			case <-done:
				// Wait for all existing sends to finish before marking ourselves done.
				sendWg.Wait()
				wg.Done()
			}
		}
	}()
}
