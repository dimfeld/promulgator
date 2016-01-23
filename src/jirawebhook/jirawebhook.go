package jirawebhook

import (
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	outChan chan *model.ChatMessage, done chan struct{}) {

}
