package app_models

import "github.com/google/uuid"

type WsClientMsg struct {
	Uuid uuid.UUID
	Data []byte
}

type WsSubscribeMsg struct {
	Subscribe string
	Data      []byte
}

type WsBroadcastMsg struct {
	Data []byte
}
