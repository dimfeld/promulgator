package main

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/spacemonkeygo/spacelog"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/jiraclient"
	"github.com/dimfeld/promulgator/jirawebhook"
	"github.com/dimfeld/promulgator/model"
	"github.com/dimfeld/promulgator/slackclient"
)

func readConfig() (c *model.Config, err error) {
	c = new(model.Config)
	err = envconfig.Process("PROMULGATOR", c)
	if c.SlackSlashCommandKey == "" {
		c.SlackSlashCommandKey = c.SlackKey
	}

	if !strings.HasSuffix(c.JiraUrl, "/") {
		c.JiraUrl = c.JiraUrl + "/"
	}
	return
}

func main() {
	var config *model.Config
	var err error

	spacelog.MustSetup("promulgator", spacelog.SetupConfig{Output: "stderr", Stdlevel: "warn"})
	mainLogger := spacelog.GetLoggerNamed("main")

	if config, err = readConfig(); err != nil {
		mainLogger.Critf("Error reading config: %s\n", err.Error())
		os.Exit(1)
	}

	// Channel for sending messages out to Slack.
	slackOutChan := make(chan *model.ChatMessage, 1)

	// This channel is closed when the server is done. At present, there is no
	// reason to ever do this, but the option is here.
	closeChan := make(chan struct{})

	wg := &sync.WaitGroup{}

	// Route incoming messages from Slack to the appropriate destination.
	router := commandrouter.New()

	jiraclient.Start(config, wg, router, slackOutChan, closeChan)
	jirawebhook.Start(config, wg, slackOutChan, closeChan)
	slackclient.Start(config, wg, slackOutChan, router, closeChan)

	go func() {
		// TODO Support TLS. Unimportant for now only because this runs solely
		// within our own network, using nginx for TLS termination.
		err := http.ListenAndServe(config.WebHookBind, nil)
		if err != nil {
			mainLogger.Crit(err.Error())
			os.Exit(1)
		}
		mainLogger.Noticef("Listening on %s", config.WebHookBind)
	}()

	wg.Wait()
	router.Close()
}
