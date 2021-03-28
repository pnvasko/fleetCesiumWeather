package app_models

import (
	base "app.base"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
)

const (
	SubscribeMsg = iota
	UnSubscribeMsg
)

type SubscribeRequest struct {
	RequestId     string
	UuidWsClients uuid.UUID
	Type          int
	Name          string `json:"name"`
}

func (sr *SubscribeRequest) Load(data json.RawMessage) error {
	return base.JsonProcessor.NewDecoder(bytes.NewReader(data)).Decode(sr)
}
