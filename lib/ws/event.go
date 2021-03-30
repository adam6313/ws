package ws

import (
	"github.com/axolotlteam/thunder/st"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// EventHandler -
type EventHandler func(*Event)

// Event -
type Event struct {
	Name string      `json:"event"`
	Data interface{} `json:"data"`
}

// NewEvent -
func NewEvent(rawData []byte) (*Event, error) {
	e := new(Event)

	err := json.Unmarshal(rawData, e)
	if err != nil {
		return nil, st.ErrorDataParseFailed
	}

	return e, nil
}

// Raw -
func (e *Event) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}
