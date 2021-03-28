package app_models

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
)

type RpcJsonResponse struct {
	Id            string          `json:"id,omitempty"`
	Result        []byte          `json:"result,omitempty"`
	Error         *RpcJsonError   `json:"error,omitempty"`
	Method        string          `json:"method,omitempty"`
	Params        json.RawMessage `json:"params,omitempty"`
	Jsonprc       string          `json:"jsonrpc"`
	UuidWsClients uuid.UUID
}

func NewResponse() *RpcJsonResponse {
	return &RpcJsonResponse{
		Jsonprc: "2.0",
		Error:   &RpcJsonError{},
	}
}

func (e *RpcJsonResponse) MarshalJSON() ([]byte, error) {
	var bmsg bytes.Buffer

	bmsg.Write(startRpcJsonResponseTmp)
	bmsg.WriteString(JsonrpcVersion)
	bmsg.Write(idRpcJsonResponseTmp)
	bmsg.Write(qRpcJsonResponseTmp)
	bmsg.WriteString(e.Id)
	bmsg.Write(qRpcJsonResponseTmp)
	if e.Method != "" {
		bmsg.Write(methodRpcJsonResponseTmp)
		bmsg.WriteString(e.Method)
		bmsg.Write(qRpcJsonResponseTmp)
	}
	if len(e.Params) > 0 {
		bmsg.Write(paramsRpcJsonResponseTmp)
		bmsg.Write(e.Params)
	}
	bmsg.Write(resultRpcJsonResponseTmp)
	if len(e.Result) > 0 {
		bmsg.Write(e.Result)
	} else {
		bmsg.Write(noneRpcJsonErrorTmp)
	}

	if e.Error != nil {
		if e.Error.Code != 0 {
			berr, _ := e.Error.MarshalJSON()
			bmsg.Write(errorRpcJsonResponseTmp)
			bmsg.Write(berr)
		}
	}
	bmsg.Write(endRpcJsonResponseTmp)

	return bmsg.Bytes(), nil
}

func PrepareRpcErrorResponse(id string, code int, uuidWsClients uuid.UUID) *RpcJsonResponse {
	resp := &RpcJsonResponse{
		Id:            id,
		UuidWsClients: uuidWsClients,
		Jsonprc:       JsonrpcVersion, // GetRpcJsonError(InvalidRequest, nil),
		Error:         GetRpcJsonError(code, nil),
	}

	return resp
}
