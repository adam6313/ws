package ws

import (
	"log"
	"net/http"

	"github.com/axolotlteam/thunder/st"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	// Subprotocols:    []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
}

// WebSocket -
type WebSocket struct {
	Conn   *websocket.Conn
	Out    chan []byte
	In     chan []byte
	Close  chan bool
	Events map[string]EventHandler
}

// NewWebSocket -
func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	subprotocols := r.Header["Sec-Websocket-Protocol"]
	upgrader.Subprotocols = subprotocols

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("An error occured while upgrading the connection: %v", err)
		return nil, st.ErrorConnectFailed
	}

	ws := &WebSocket{
		Conn:   conn,
		Out:    make(chan []byte),
		In:     make(chan []byte),
		Close:  make(chan bool),
		Events: make(map[string]EventHandler),
	}

	go ws.Reader()
	go ws.Writer()

	return ws, nil
}

// Reader -
func (ws *WebSocket) Reader() {
	defer func() {
		ws.Conn.Close()
	}()

	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS Message Error: %v", err)
			}
			log.Printf("WS Message Errors: %v", err)
			ws.Close <- true
			break
		}
		event, err := NewEvent(message)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
		} else {
			log.Printf("message: %v", event)
		}

		if action, ok := ws.Events[event.Name]; ok {
			action(event)
		}

	}
}

// Writer -
func (ws *WebSocket) Writer() {
	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				ws.Conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
			}

			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			w.Close()
		}
	}
}

// On -
func (ws *WebSocket) On(eventName string, action EventHandler) *WebSocket {
	ws.Events[eventName] = action
	return ws
}
