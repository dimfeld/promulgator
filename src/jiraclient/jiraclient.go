package jiraclient

import (
	"model"
	"sync"
)

func Start(config *model.Config, wg *sync.WaitGroup,
	inChan chan *model.TrackerUpdate,
	outChan chan *model.ChatMessage, done chan struct{}) {

}
