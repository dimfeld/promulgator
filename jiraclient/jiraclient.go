package jiraclient

import (
	"sync"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

// Start the Jira client goroutine
func Start(config *model.Config, wg *sync.WaitGroup,
	commandrouter *commandrouter.Router,
	outChan chan *model.ChatMessage, done chan struct{}) {

}
