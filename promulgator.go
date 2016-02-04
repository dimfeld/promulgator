package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"sync"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/jiraclient"
	"github.com/dimfeld/promulgator/jirawebhook"
	"github.com/dimfeld/promulgator/model"
	"github.com/dimfeld/promulgator/slackclient"
)

func readConfig() (c *model.Config, err error) {
	c = new(model.Config)
	err = envconfig.Process("PROMULGATOR", c)
	return
}

func main() {
	var config *model.Config
	var err error

	if config, err = readConfig(); err != nil {
		fmt.Printf("Error reading config: %s\n", err.Error())
		os.Exit(1)
	}

	// Channel for sending messages out to Slack.
	slackChatChan := make(chan *model.ChatMessage, 1)

	// This channel is closed when the server is done. At present, there is no
	// reason to ever do this, but the option is here.
	closeChan := make(chan struct{})

	wg := &sync.WaitGroup{}

	// Route incoming messages from Slack to the appropriate destination.
	router := commandrouter.New()

	jiraclient.Start(config, wg, router, slackChatChan, closeChan)
	jirawebhook.Start(config, wg, slackChatChan, closeChan)
	slackclient.Start(config, wg, slackChatChan, router, closeChan)

	wg.Wait()
	router.Close()
}
