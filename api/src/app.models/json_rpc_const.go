package app_models

const (
	JsonrpcVersion       = "2.0"
	contentType          = "Content-Type"
	contentTypeJSON      = "application/json"
	UnauthorizedError    = 1
	UnauthorizedErrorMsg = "Action is not authorized"
	notAllowed           = 2
	notAllowedErrorMsg   = "Action is not allowed"
	ParseError           = -32700
	parseErrorMsg        = "Parse error"
	InvalidRequest       = -32600
	invalidRequestMsg    = "Invalid Request"
	MethodNotFound       = -32601
	methodNotFoundMsg    = "Method not found"
	InvalidParams        = -32602
	invalidParamsMsg     = "Invalid params"
	InternalError        = -32603
	internalErrorMsg     = "Internal error"
	ServerError          = -32000
	serverErrorMsg       = "Server error"
)
