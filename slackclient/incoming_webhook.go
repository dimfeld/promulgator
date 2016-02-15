package slackclient

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/dimfeld/promulgator/commandrouter"
	"github.com/dimfeld/promulgator/model"
)

func StartIncomingWebhook(config *model.Config, commandrouter *commandrouter.Router) {
	http.HandleFunc("/slackhook", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logIn.Warnf("Form decode error: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		token := r.PostForm.Get("token")
		if token != config.SlackSlashCommandKey {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		response := make(chan string, 1)

		ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		processIncomingSlack(ctx, commandrouter,
			r.PostForm.Get("text"), r.PostForm.Get("user_name"), r.PostForm.Get("channel_name"),
			true, response)

		select {
		case s := <-response:
			w.WriteHeader(http.StatusOK)
			if s != "" {
				w.Write([]byte(s))
			} else {
				w.Write([]byte("Success"))
			}

		case <-ctx.Done():
			w.WriteHeader(http.StatusGatewayTimeout)
		}
	})
}
