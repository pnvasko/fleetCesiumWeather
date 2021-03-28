package app_jsonrpc

import "time"

const (
	DefaultConcurrency     = 8 * 1024
	DefaultRequestTimeout  = 60 * time.Second
	DefaultPendingMessages = 32 * 1024
	DefaultFlushDelay      = -1
	DefaultBufferSize      = 64 * 1024
	DefaultMethodPrefix    = "RPC"
	DefaultSourceChanel    = "rpc.in"

	IncorrectRequestTypeError     = 1
	IncorrectRequestTypeErrorMsg  = "Incorrect request type"
	UnknownServiceNameError       = 2
	UnknownServiceNameErrorMsg    = "Unknown service name"
	UnknownMethodError            = 3
	UnknownMethodErrorMsg         = "Unknown method"
	UnexpectedRequestTypeError    = 4
	UnexpectedRequestTypeErrorMsg = "Unexpected request type for method"
)
