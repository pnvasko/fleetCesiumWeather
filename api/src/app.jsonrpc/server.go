package app_jsonrpc

import (
	base "app.base"
	models "app.models"
	"context"
	"fmt"
)

type JsonRpcServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	inQueue       chan *models.RpcJsonRequest
	outQueue      chan *models.RpcJsonResponse
	wsChanHandler func(*models.WsClientMsg)

	dispatcher *dispatcherRpc
	workerPool *base.WorkerPool

	log *base.Logger
}

func NewJsonRpcServer(ctx context.Context, prefix string) (*JsonRpcServer, error) {
	var err error
	var ok bool
	srv := &JsonRpcServer{
		inQueue:  make(chan *models.RpcJsonRequest, DefaultPendingMessages),
		outQueue: make(chan *models.RpcJsonResponse, DefaultPendingMessages),
	}
	srv.ctx, srv.cancel = context.WithCancel(ctx)

	srv.log = base.GetLog(srv.ctx, srv)
	srv.workerPool, err = base.GetWorkerPool(srv.ctx)
	if err != nil {
		return nil, err
	}

	methodPrefix := prefix
	if prefix == "" {
		methodPrefix = DefaultMethodPrefix
	}

	srv.dispatcher = newDispatcher(srv.ctx, methodPrefix)

	wsChanHandler := srv.ctx.Value("wsInCh")
	if srv.wsChanHandler, ok = wsChanHandler.(func(*models.WsClientMsg)); !ok {
		return nil, fmt.Errorf("NewJsonRpcServer error get wsChanHandler from context")
	}

	return srv, nil
}

func (srv *JsonRpcServer) OutCh() <-chan *models.RpcJsonResponse {
	return srv.outQueue
}

func (srv *JsonRpcServer) InCh() chan<- *models.RpcJsonRequest {
	return srv.inQueue
}

func (srv *JsonRpcServer) AddService(serviceName string, service interface{}) error {
	return srv.dispatcher.addService(serviceName, service)
}

func (srv *JsonRpcServer) SetNatsHandlerFuncHandler() error {
	var err error
	srv.dispatcher.handler, err = srv.dispatcher.newJsonHandlerFunc()
	if err != nil {
		return err
	}
	return nil
}

func (srv *JsonRpcServer) Loop() {
	for {
		select {
		case <-srv.ctx.Done():
			return
		case msg, ok := <-srv.inQueue:
			if !ok {
				continue
			}
			srv.submit(func() {
				resp := srv.dispatcher.handler("BaseRPC", msg)
				data, _ := resp.MarshalJSON()
				cliMsg := &models.WsClientMsg{
					Uuid: msg.UuidWsClients,
					Data: data,
				}
				srv.wsChanHandler(cliMsg)
			})
		}
	}
}

func (srv *JsonRpcServer) Stop() error {
	srv.cancel()
	return nil
}

func (srv *JsonRpcServer) submit(f func()) {
	srv.workerPool.Submit(f)
}
