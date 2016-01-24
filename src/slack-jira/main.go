package main

import (
	"commandrouter"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"jiraclient"
	"jirawebhook"
	"model"
	"os"
	"slackclient"
	"sync"
)

func ReadConfig() (c *model.Config, err error) {
	c = new(model.Config)
	err = envconfig.Process("sj", c)
	return
}

func main() {
	var config *model.Config
	var err error

	if config, err = ReadConfig(); err != nil {
		fmt.Printf("Error reading config: %s\n", err.Error())
		os.Exit(1)
	}

	// Channel for sending messages out to Slack.
	slackChatChan := make(chan *model.ChatMessage, 1)

	// This channel is closed when the server is done. At present, there is no
	// reason to ever do this, but the option is here.
	closeChan := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(3)

	// Route incoming messages from Slack to the appropriate destination.
	router := commandrouter.New()

	jiraclient.Start(config, wg, router, slackChatChan, closeChan)
	jirawebhook.Start(config, wg, slackChatChan, closeChan)
	slackclient.Start(config, wg, slackChatChan, router, closeChan)

	wg.Wait()
	router.Close()
}
