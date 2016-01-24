package jiraclient

import (
	"commandrouter"
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	commandrouter *commandrouter.Router,
	outChan chan *model.ChatMessage, done chan struct{}) {

}
