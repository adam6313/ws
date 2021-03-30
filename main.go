package main

import (
	"encoding/json"
	"ws/lib/ws"

	"github.com/gin-gonic/gin"
)

func main() {

}

func example(c *gin.Context) {
	// new websocket
	w, err := ws.NewWebSocket(c.Writer, c.Request)
	if err != nil {
		return
	}

	// listen `message`, can listen other
	w.On("message", func(e *ws.Event) {
		d, err := json.Marshal(e.Data)
		if err != nil {
			return
		}

		// send message to websocket
		w.Out <- (&ws.Event{
			Name: "message",
			Data: string(d),
		}).Raw()

	})
}
