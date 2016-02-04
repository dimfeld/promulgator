package commandrouter

import (
	"github.com/dimfeld/promulgator/model"
)

func chatBuffer(input, output chan *model.ChatMessage, done chan struct{}) {
	pending := []*model.ChatMessage{}
	for {
		var maybeOutput chan *model.ChatMessage
		var outputItem *model.ChatMessage
		if len(pending) != 0 {
			// Try to send only if there's an item to send.
			outputItem = pending[0]
			maybeOutput = output
		}

		select {
		case newItem := <-input:
			if newItem == nil {
				// Input channel was closed
				close(output)
				return
			}
			pending = append(pending, newItem)

		case maybeOutput <- outputItem:
			pending = pending[1:]

		case <-done:
			close(output)
			return
		}
	}
}

// ChatBuffer provides an unbounded buffer for ChatMessage objects
func ChatBuffer(input chan *model.ChatMessage, done chan struct{}) chan *model.ChatMessage {
	output := make(chan *model.ChatMessage, 1)
	go chatBuffer(input, output, done)
	return output
}

func matchBuffer(input, output chan Match, done chan struct{}) {
	pending := []Match{}
	for {
		var maybeOutput chan Match
		var outputItem Match
		if len(pending) != 0 {
			// Try to send only if there's an item to send.
			outputItem = pending[0]
			maybeOutput = output
		}

		select {
		case newItem := <-input:
			if newItem.Message == nil {
				// input channel was closed
				close(output)
				return
			}
			pending = append(pending, newItem)

		case maybeOutput <- outputItem:
			// Sent the item, so remove it from the list
			pending = pending[1:]

		case <-done:
			close(output)
			return
		}
	}
}

// MatchBuffer provides an unbounded buffer for Match objects
func MatchBuffer(input chan Match, done chan struct{}) chan Match {
	output := make(chan Match, 1)
	go matchBuffer(input, output, done)
	return output
}
