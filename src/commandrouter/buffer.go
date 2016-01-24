package commandrouter

import (
	"model"
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
			pending = append(pending, newItem)

		case maybeOutput <- outputItem:
			pending = pending[1:]

		case <-done:
			close(output)
			return
		}
	}
}

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

func MatchBuffer(input chan Match, done chan struct{}) chan Match {
	output := make(chan Match, 1)
	go matchBuffer(input, output, done)
	return output
}
