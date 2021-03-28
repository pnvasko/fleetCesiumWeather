package app_services

import (
	base "app.base"
	"context"
	"encoding/json"
)

type RpcFleetService struct {
	ctx    context.Context
	cancel context.CancelFunc

	log *base.Logger
}

func NewRpcFleetService(ctx context.Context) (*RpcFleetService, error) {
	rfs := &RpcFleetService{}
	rfs.ctx, rfs.cancel = context.WithCancel(ctx)
	rfs.log = base.GetLog(rfs.ctx, rfs)

	return rfs, nil
}

func (rfs *RpcFleetService) RPCGetFleetScheduled(req json.RawMessage) ([]byte, error) {
	rfs.log.Debug("RpcFleetService.RPCGetFleetScheduled start")
	data := []byte(`{"test": 1, "next": "low"}`)
	return data, nil
}
