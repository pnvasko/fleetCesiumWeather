package app_models

import (
	base "app.base"
	"bytes"
)

type DebugFleetInfoDbMsg struct {
	Id    int64             `json:"id"`
	Fleet string            `json:"fleet"`
	Route FleetInfoResponse `json:"route"`
}

func (dbMsg *DebugFleetInfoDbMsg) Load(data []byte) error {
	return base.JsonProcessor.NewDecoder(bytes.NewReader(data)).Decode(dbMsg)
}

func (dbMsg *DebugFleetInfoDbMsg) GetWsSubscribeMsg() ([]byte, error) {
	// data := FleetInfoResponse{}
	data, err := base.JsonProcessor.Marshal(dbMsg.Route)
	return data, err
}
