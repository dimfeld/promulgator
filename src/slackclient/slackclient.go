package slackclient

import (
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.ChatMessage, outChan chan *model.TrackerUpdate, done chan struct{}) {

}
