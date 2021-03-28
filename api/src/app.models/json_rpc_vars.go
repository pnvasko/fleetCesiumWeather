package app_models

var (
	startRpcJsonResponseTmp  = []byte(`{"jsonrpc":"`)
	idRpcJsonResponseTmp     = []byte(`", "id":`)
	resultRpcJsonResponseTmp = []byte(`, "result":`)
	methodRpcJsonResponseTmp = []byte(`, "method":"`)
	qRpcJsonResponseTmp      = []byte(`"`)
	paramsRpcJsonResponseTmp = []byte(`, "params":`)
	errorRpcJsonResponseTmp  = []byte(`, "error":`)
	endRpcJsonResponseTmp    = []byte(`}`)

	methodJsonRequestTmp   = []byte(`, "method":"`)
	commaJsonTmp           = []byte(`"`)
	startRpcJsonErrorTmp   = []byte(`{"code":`)
	messageRpcJsonErrorTmp = []byte(`, "message":"`)
	dataRpcJsonErrorTmp    = []byte(`", "data":`)
	endRpcJsonErrorTmp     = []byte(`}`)
	noneRpcJsonErrorTmp    = []byte(`null`)
)
