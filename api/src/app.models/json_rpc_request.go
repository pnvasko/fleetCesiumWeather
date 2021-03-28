package app_models

import (
	base "app.base"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

type RpcJsonRequest struct {
	Id            string          `json:"id,omitempty"`
	Jsonrpc       string          `json:"jsonrpc"`
	Method        string          `json:"method"`
	Params        json.RawMessage `json:"payload,omitempty"`
	UuidWsClients uuid.UUID
}

func NewRpcJsonRequest() (*RpcJsonRequest, error) {
	jr := &RpcJsonRequest{
		Id:      uuid.Must(uuid.NewRandom()).String(),
		Jsonrpc: JsonrpcVersion,
	}
	return jr, nil
}

func (rjr *RpcJsonRequest) isValidVersion() bool {
	return rjr.Jsonrpc == JsonrpcVersion
}

func (rjr *RpcJsonRequest) isMethodEmpty() bool {
	return len(rjr.Method) == 0
}

func (rjr *RpcJsonRequest) Validate() error {
	if !rjr.isValidVersion() {
		return errors.New("invalid json-rpc version")
	}

	if rjr.isMethodEmpty() {
		return errors.New("method is empty")
	}
	return nil
}

func (rjr *RpcJsonRequest) Load(data []byte) error {
	return base.JsonProcessor.NewDecoder(bytes.NewReader(data)).Decode(rjr)
}

func (rjr *RpcJsonRequest) MarshalJSON() ([]byte, error) {
	var bmsg bytes.Buffer
	bmsg.Write(startRpcJsonResponseTmp)
	bmsg.WriteString(JsonrpcVersion)
	bmsg.Write(idRpcJsonResponseTmp)
	bmsg.Write(commaJsonTmp)
	bmsg.WriteString(rjr.Id)
	bmsg.Write(commaJsonTmp)
	bmsg.Write(methodJsonRequestTmp)
	bmsg.WriteString(rjr.Method)
	bmsg.Write(commaJsonTmp)
	bmsg.Write(paramsRpcJsonResponseTmp)
	bmsg.Write(rjr.Params)
	bmsg.Write(endRpcJsonResponseTmp)

	return bmsg.Bytes(), nil
}

func (rjr *RpcJsonRequest) SimpleMarshalJSON() []byte {
	data, _ := rjr.MarshalJSON()
	return data
}
