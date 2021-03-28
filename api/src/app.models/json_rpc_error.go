package app_models

import (
	"bytes"
	"fmt"
	"strconv"
)

type RpcJsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []byte `json:"data,omitempty"`
}

func GetRpcJsonError(code int, data []byte) *RpcJsonError {
	var msg string
	switch code {
	case ParseError:
		msg = parseErrorMsg
	case InvalidRequest:
		msg = invalidRequestMsg
	case MethodNotFound:
		msg = methodNotFoundMsg
	case InvalidParams:
		msg = invalidParamsMsg
	case InternalError:
		msg = internalErrorMsg
	case ServerError:
		msg = serverErrorMsg
	case UnauthorizedError:
		msg = UnauthorizedErrorMsg
	case notAllowed:
		msg = notAllowedErrorMsg
	default:
		msg = "Unknow error"
	}
	return &RpcJsonError{
		Code:    int(code),
		Message: msg,
		Data:    data,
	}
}

func (e RpcJsonError) Error() string {
	return e.Message
}

func (e *RpcJsonError) MarshalJSON() ([]byte, error) {
	var bmsg bytes.Buffer

	bmsg.Write(startRpcJsonErrorTmp)
	bmsg.WriteString(fmt.Sprintf("%d", e.Code))
	bmsg.Write(messageRpcJsonErrorTmp)
	bmsg.WriteString(e.Message)
	bmsg.Write(dataRpcJsonErrorTmp)
	if len(e.Data) > 0 {
		bmsg.Write(e.Data)
	} else {
		bmsg.Write(noneRpcJsonErrorTmp)
	}

	bmsg.Write(endRpcJsonErrorTmp)

	return bmsg.Bytes(), nil
}

func NewByteErrorResponse(id int64, err string) []byte {
	var out bytes.Buffer
	out.Write([]byte(`{"id":`))
	out.WriteString(strconv.FormatInt(id, 10))
	out.Write([]byte(`,"error":"`))
	out.WriteString(err)
	out.Write([]byte(`","jsonrpc":"2.0"}`))
	return out.Bytes()
}
