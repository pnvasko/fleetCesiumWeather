package app_wsservice

import (
	base "app.base"
	models "app.models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

var (
	subscribeCmd   = []byte(`"method":"Subscribe"`)
	unsubscribeCmd = []byte(`"method":"Unsubscribe"`)
	maxPacketSize  = 1024 * 1024
)

type wsClient struct {
	mu          sync.Mutex
	uuid        uuid.UUID
	subscribers map[string]struct{}

	ctx    context.Context
	cancel context.CancelFunc

	workerPool     *base.WorkerPool
	rpcChanHandler handlerInCh

	conn *websocket.Conn
	log  *base.Logger
}

type handlerOutCh func() <-chan *models.RpcJsonResponse

// type handlerInCh func(*models.RpcJsonRequest)
type handlerInCh func(interface{})

func newWsClient(ctx context.Context, conn *websocket.Conn) (*wsClient, error) {
	var err error
	var ok bool

	wsc := &wsClient{
		uuid:        uuid.Must(uuid.NewRandom()),
		conn:        conn,
		subscribers: make(map[string]struct{}),
	}

	wsc.ctx, wsc.cancel = context.WithCancel(ctx)
	wsc.log = base.GetLog(wsc.ctx, wsc)
	if wsc.workerPool, err = base.GetWorkerPool(wsc.ctx); err != nil {
		return nil, err
	}

	rpcChanHandler := wsc.ctx.Value("rpcInCh")
	// if wsc.rpcChanHandler, ok = rpcChanHandler.(func(*models.RpcJsonRequest)); !ok {
	if wsc.rpcChanHandler, ok = rpcChanHandler.(func(interface{})); !ok {
		return nil, fmt.Errorf("newWsClient error get handlerInCh from context")
	}

	/*
		hInCh, ok := wsc.ctx.Value("rpcInCh").(handlerInCh)
		if !ok {
			return nil, fmt.Errorf("newWsClient error get rpcChan from context")
		}
		wsc.rpcChan = hInCh()

	*/
	return wsc, nil
}

func (cli *wsClient) submit(f func()) {
	cli.workerPool.Submit(f)
	return
}

func (cli *wsClient) Subscribe(name string) error {
	cli.mu.Lock()
	if _, ok := cli.subscribers[name]; !ok {
		cli.subscribers[name] = struct{}{}
	}
	cli.mu.Unlock()

	return nil
}

func (cli *wsClient) UnSubscribe(name string) error {
	cli.mu.Lock()
	cli.mu.Unlock()
	if _, ok := cli.subscribers[name]; ok {
		delete(cli.subscribers, name)
	}
	return nil
}

func (cli *wsClient) TestSubscribe(name string) bool {
	cli.mu.Lock()
	cli.mu.Unlock()
	_, ok := cli.subscribers[name]
	return ok
}

func (cli *wsClient) Receive() {
	buffer := make([]byte, maxPacketSize)
	for {
		select {
		case <-cli.ctx.Done():
			return
		default:
			_, reader, err := cli.conn.Reader(cli.ctx)
			if err != nil {
				cli.log.Error("wsClient.Receive conn.Reader error: ", err)
				return
			}
			n, err := reader.Read(buffer)
			if err != nil && err != io.EOF {
				cli.log.Error("wsClient.Receive Read(buffer) error: ", err)
				return
			}
			if n > 0 {
				// var out bytes.Buffer
				tmp := make([]byte, n)
				copy(tmp, buffer)
				cli.submit(func() {
					req := &models.RpcJsonRequest{}
					if err := req.Load(tmp); err != nil {
						cli.log.Error("wsClient.Receive RpcJsonRequest load error: ", err)
						return
					}
					req.UuidWsClients = cli.uuid
					switch req.Method {
					case "Subscribe", "UnSubscribe":
						sr := models.SubscribeRequest{}
						if err := sr.Load(req.Params); err != nil {
							cli.log.Error("wsClient.Receive Subscribe/UnSubscribe SubscribeRequest load params error: ", err)
							resp := models.PrepareRpcErrorResponse(req.Id, 0, cli.uuid)
							data, _ := resp.MarshalJSON()
							if err := cli.WriterByte(cli.ctx, data); err != nil {
								cli.log.Error("wsClient.Receive Subscribe/UnSubscribe WriterByte error: ", err)
							}
							return
						}
						actionType := models.SubscribeMsg
						if req.Method == "UnSubscribe" {
							actionType = models.UnSubscribeMsg
						}
						cli.rpcChanHandler(&models.SubscribeRequest{
							RequestId:     req.Id,
							UuidWsClients: cli.uuid,
							Type:          actionType,
							Name:          sr.Name,
						})
					default:
						cli.rpcChanHandler(req)
					}
				})
			}
		}
	}
}

func (cli *wsClient) WriterByte(ctx context.Context, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return cli.conn.Write(ctx, websocket.MessageText, msg)
}

func (cli *wsClient) Writer(ctx context.Context, msg interface{}) error {
	w, err := cli.conn.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}
	defer func() {
		if err := w.Close(); err != nil {
			cli.log.Error("wsClient.Writer error close: ", err)
		}
	}()

	e := base.JsonProcessor.NewEncoder(w)
	if err := e.Encode(msg); err != nil {
		return fmt.Errorf("Server.writer failed to encode json: %w", err)
	}

	return nil
}
