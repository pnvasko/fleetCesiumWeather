package app_wsservice

import (
	base "app.base"
	jsonrpc "app.jsonrpc"
	models "app.models"
	services "app.services"
	"github.com/jackc/pgx/pgxpool"
	"time"

	"context"
	"fmt"
	"net/http"
)

type WsService struct {
	config *Config

	ctx    context.Context
	cancel context.CancelFunc

	wshub      *wsHub
	workerPool *base.WorkerPool
	rpcServer  *jsonrpc.JsonRpcServer

	pgxpool *pgxpool.Pool
	ps      *postgresSubscribes
	http    *http.Server
	log     *base.Logger
}

const defaultRpcService = "BaseRPC"

func NewWsService(ctx context.Context, path string) (*WsService, error) {
	var err error

	ws := &WsService{
		config: newConfig(),
	}

	if err = ws.config.Load(path); err != nil {
		return nil, fmt.Errorf("NewService error load config: %s", err)
	}

	ws.log = base.NewLoger("WsService", ws.config.LogPath, ws.config.Debug)

	ws.ctx, ws.cancel = context.WithCancel(ctx)
	ws.ctx = context.WithValue(ws.ctx, "config", ws.config)
	ws.ctx = context.WithValue(ws.ctx, "logger", ws.log)

	ws.ctx = context.WithValue(ws.ctx, "rpcInCh", ws.rpcInChHandler)
	ws.ctx = context.WithValue(ws.ctx, "wsInCh", ws.wsInChHandler)

	ws.workerPool, err = base.NewWorkerPool(ws.ctx, ws.config.WorkerPoolSize)
	if err != nil {
		return nil, err
	}
	ws.ctx = context.WithValue(ws.ctx, "worker_pool", ws.workerPool)

	if ws.pgxpool, err = initDBDriver(ctx, ws.config); err != nil {
		return nil, fmt.Errorf("Error connect to DB server: %s", err)
	}
	ws.ctx = context.WithValue(ws.ctx, "pgxpool", ws.pgxpool)

	ws.rpcServer, err = jsonrpc.NewJsonRpcServer(ws.ctx, "RPC")
	if err != nil {
		return nil, err
	}

	// if err := ws.rpcServer.AddService(defaultRpcService, reflect.TypeOf((*services.RpcServices)(nil)).Elem()); err != nil {
	rfs, err := services.NewRpcFleetService(ws.ctx)
	if err != nil {
		return nil, err
	}
	if err := ws.rpcServer.AddService("Fleet", rfs); err != nil {
		return nil, err
	}
	if err := ws.rpcServer.SetNatsHandlerFuncHandler(); err != nil {
		return nil, err
	}

	ws.http = &http.Server{Addr: fmt.Sprintf(":%d", ws.config.HttpConfig.Port), Handler: ws.newRouter()}

	if ws.wshub, err = newWsHub(ws.ctx); err != nil {
		return nil, err
	}

	if ws.ps, err = newPostgresSubscribe(ws.ctx); err != nil {
		return nil, err
	}

	ws.ps.addHandler("db.tbfm.updated", ws.tbfmUpdatedHandler)
	ws.ps.addHandler("db.debug.updated", ws.debugUpdatedHandler)

	return ws, nil
}

func (ws *WsService) debugUpdatedHandler(msg *postgresPubSubContext) error {
	dbMsg := models.DebugFleetInfoDbMsg{}
	if err := dbMsg.Load([]byte(msg.Msg.Payload)); err != nil {
		ws.log.Error("WsService.debugUpdatedHandler error load dbmsg payload: ", err)
		return err
	}

	data, _ := dbMsg.GetWsSubscribeMsg()
	fleetInfoData, _ := models.NewRpcJsonRequest()
	fleetInfoData.Method = "FleetInfoData"
	fleetInfoData.Params = data

	wsMsg := &models.WsSubscribeMsg{
		Subscribe: dbMsg.Fleet,
		Data:      fleetInfoData.SimpleMarshalJSON(),
	}

	ws.wshub.InSubscribCh() <- wsMsg

	return nil
}

func (ws *WsService) tbfmUpdatedHandler(msg *postgresPubSubContext) error {
	fmt.Println("WsService.tbfmUpdatedHandler Channel: ", msg.Msg.Channel)
	fmt.Println("WsService.tbfmUpdatedHandler Payload: ", msg.Msg.Payload)
	return nil
}

// func (ws *WsService) rpcInChHandler(msg *models.RpcJsonRequest) {
func (ws *WsService) rpcInChHandler(msg interface{}) {
	switch msg.(type) {
	case *models.RpcJsonRequest:
		ws.rpcServer.InCh() <- msg.(*models.RpcJsonRequest)
	case *models.SubscribeRequest:
		ws.wshub.InSubscrib() <- msg.(*models.SubscribeRequest)
	default:
		ws.log.Error("WsService.rpcInChHandler error msg type: %T", msg)
	}
}

func (ws *WsService) wsInChHandler(msg *models.WsClientMsg) {
	fmt.Println("WsService.wsInChHandler")
	ws.wshub.InClientCh() <- msg
}

func (ws *WsService) loop() {
	broadcastTicker := time.NewTicker(time.Duration(180) * time.Second)
	defer broadcastTicker.Stop()

	fleetDataTicker := time.NewTicker(time.Duration(120) * time.Second)
	defer fleetDataTicker.Stop()

	resp, _ := models.NewRpcJsonRequest()
	resp.Method = "Broadcast"
	resp.Params = []byte(`{"broadcast": "ping"}`)
	msg := &models.WsBroadcastMsg{
		Data: resp.SimpleMarshalJSON(),
	}

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-broadcastTicker.C:
			ws.wshub.inBroadcastCh <- msg
		case <-fleetDataTicker.C:
			fleetData, _ := models.NewRpcJsonRequest()
			fleetData.Method = "FleetData"
			fleetData.Params = []byte(fmt.Sprintf(`{"fleet": "%s","lat":%f,"lng":%f,"alt":%f,"speed":%f,"timestamp": %d}`,
				"SWA2558",
				35.555,
				-120.22,
				9000.0,
				850.5,
				time.Now().Unix(),
			))
			fleetMsg := &models.WsSubscribeMsg{
				Subscribe: "SWA2558",
				Data:      fleetData.SimpleMarshalJSON(),
			}
			ws.wshub.inSubscribCh <- fleetMsg
		}
	}
}

func (ws *WsService) Run() error {
	ws.log.Debug("WsService.Run...")
	go ws.rpcServer.Loop()
	ws.wshub.run()

	go ws.loop()
	go ws.ps.loop()

	go func() {
		defer ws.cancel()
		ws.log.Info("WsService.Run http server listening on: ", ws.config.HttpConfig.Port)
		if err := ws.http.ListenAndServe(); err != http.ErrServerClosed {
			ws.log.Fatalf("WsService.Run ListenAndServe error:", err)
		}
	}()

	return nil
}

func (ws *WsService) Complete() <-chan struct{} {
	return ws.ctx.Done()
}

func (ws *WsService) Shutdown() {
	ws.log.Debug("WsService.Shutdown start...")
	complete := make(chan struct{})
	ws.ps.close(complete)
	<-complete
	ws.wshub.close()
	_ = ws.rpcServer.Stop()
	ws.workerPool.Release()
	ws.log.Debug("WsService.Shutdown finish.")
	ws.cancel()
}
