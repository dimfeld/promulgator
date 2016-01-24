package slackclient

import (
	"commandrouter"
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.ChatMessage,
	commandrouter *commandrouter.Router, done chan struct{}) {

}
