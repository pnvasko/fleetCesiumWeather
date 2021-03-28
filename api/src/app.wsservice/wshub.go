package app_wsservice

import (
	base "app.base"
	models "app.models"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type wsHub struct {
	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	wsclients map[*wsClient]uuid.UUID
	clients   map[uuid.UUID]*wsClient

	subscribes map[string]map[uuid.UUID]*wsClient

	inSubscrib    chan *models.SubscribeRequest
	inBroadcastCh chan *models.WsBroadcastMsg
	inSubscribCh  chan *models.WsSubscribeMsg
	inClientCh    chan *models.WsClientMsg

	wsChanHandler func(*models.WsClientMsg)

	workerPool *base.WorkerPool
	log        *base.Logger
}

func newWsHub(ctx context.Context) (*wsHub, error) {
	var err error
	wsh := &wsHub{
		wsclients:     make(map[*wsClient]uuid.UUID),
		clients:       make(map[uuid.UUID]*wsClient),
		subscribes:    make(map[string]map[uuid.UUID]*wsClient),
		inSubscrib:    make(chan *models.SubscribeRequest, 100),
		inBroadcastCh: make(chan *models.WsBroadcastMsg, 100),
		inSubscribCh:  make(chan *models.WsSubscribeMsg, 100),
		inClientCh:    make(chan *models.WsClientMsg, 100),
	}
	wsh.ctx, wsh.cancel = context.WithCancel(ctx)

	wsh.log = base.GetLog(wsh.ctx, wsh)
	if wsh.workerPool, err = base.GetWorkerPool(wsh.ctx); err != nil {
		return nil, err
	}

	return wsh, nil
}

func (ws *wsHub) InClientCh() chan<- *models.WsClientMsg {
	return ws.inClientCh
}

func (ws *wsHub) InBroadcastCh() chan<- *models.WsBroadcastMsg {
	return ws.inBroadcastCh
}

func (ws *wsHub) InSubscrib() chan<- *models.SubscribeRequest {
	return ws.inSubscrib
}

func (ws *wsHub) InSubscribCh() chan<- *models.WsSubscribeMsg {
	return ws.inSubscribCh
}

func (ws *wsHub) submit(f func()) {
	ws.workerPool.Submit(f)
}

func (ws *wsHub) run() {
	go ws.loop()
	go ws.writer()
}

func (ws *wsHub) close() {
	ws.cancel()
	close(ws.inBroadcastCh)
	close(ws.inSubscribCh)
	close(ws.inClientCh)
	close(ws.inSubscrib)
}

func (ws *wsHub) registration(conn *websocket.Conn) (*wsClient, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	cli, err := newWsClient(ws.ctx, conn)

	if err != nil {
		return nil, err
	}
	_, ok := ws.wsclients[cli]
	if !ok {
		ws.clients[cli.uuid] = cli
		ws.wsclients[cli] = cli.uuid
	}

	return cli, nil
}

func (ws *wsHub) unregistration(cli *wsClient) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	_, ok := ws.wsclients[cli]
	if ok {
		delete(ws.wsclients, cli)
		delete(ws.clients, cli.uuid)
	}

	for sn, _ := range ws.subscribes {
		delete(ws.subscribes[sn], cli.uuid)
	}

	return nil
}

func (ws *wsHub) subscribeSender(id uuid.UUID, subscribe string, data []byte) {
	send := func() {
		cli, ok := ws.clients[id]
		if !ok {
			ws.log.Error(fmt.Sprintf("wsHub.subscribeSender error can't get ws client %s", id.String()))
			return
		}
		if ok := cli.TestSubscribe(subscribe); ok {
			if err := cli.WriterByte(ws.ctx, data); err != nil {
				ws.log.Error(fmt.Sprintf("wsHub.subscribeSender %s error: %s", subscribe, err))
			}
		}
	}
	ws.submit(send)
}

func (ws *wsHub) sender(id uuid.UUID, data []byte) {
	send := func() {
		cli, ok := ws.clients[id]
		if !ok {
			ws.log.Error(fmt.Sprintf("wsHub.sender error can't get ws client %s", id.String()))
			return
		}
		if err := cli.WriterByte(ws.ctx, data); err != nil {
			ws.log.Error(fmt.Sprintf("wsHub.sender %s error: %s", id.String(), err))
			return
		}
	}
	ws.submit(send)
}

func (ws *wsHub) subscribe(msg *models.SubscribeRequest) {
	sub := func() {
		cli, ok := ws.clients[msg.UuidWsClients]
		if !ok {
			ws.log.Error(fmt.Sprintf("wsHub.subscribe error can't get ws client %s", msg.UuidWsClients.String()))
			return
		}

		if err := cli.Subscribe(msg.Name); err != nil {
			ws.log.Error(fmt.Sprintf("wsHub.subscribe error subscribe [%s] for %s", msg.Name, msg.UuidWsClients.String()))
			return
		}

		ws.mu.Lock()
		_, ok = ws.subscribes[msg.Name]
		if !ok {
			ws.subscribes[msg.Name] = make(map[uuid.UUID]*wsClient)
		}
		ws.subscribes[msg.Name][msg.UuidWsClients] = cli

		ws.mu.Unlock()

		resp := models.NewResponse()
		resp.Id = msg.RequestId
		resp.Result = []byte(fmt.Sprintf(`{"result": "ok", "action": "subscribe","chanel": "%s", "timestamp": %d}`, msg.Name, time.Now().Unix()))
		data, _ := resp.MarshalJSON()
		if err := cli.WriterByte(cli.ctx, data); err != nil {
			cli.log.Error("wsHub.subscribe WriterByte error: ", err)
		}
	}
	ws.submit(sub)
}

func (ws *wsHub) unsubscribe(msg *models.SubscribeRequest) {
	unsub := func() {
		cli, ok := ws.clients[msg.UuidWsClients]
		if !ok {
			ws.log.Error(fmt.Sprintf("wsHub.sender error can't get ws client %s", msg.UuidWsClients.String()))
			return
		}
		if err := cli.UnSubscribe(msg.Name); err != nil {
			ws.log.Error(fmt.Sprintf("wsHub.unsubscribe error unsubscribe [%s] for %s", msg.Name, msg.UuidWsClients.String()))
			return
		}

		ws.mu.Lock()

		_, ok = ws.subscribes[msg.Name]
		if !ok {
			ws.subscribes[msg.Name] = make(map[uuid.UUID]*wsClient)
		}

		delete(ws.subscribes[msg.Name], msg.UuidWsClients)

		ws.mu.Unlock()

		resp := models.NewResponse()
		resp.Id = msg.RequestId
		resp.Result = []byte(fmt.Sprintf(`{"result": "ok", "action": "unsubscribe","chanel": "%s", "timestamp": %d}`, msg.Name, time.Now().Unix()))
		data, _ := resp.MarshalJSON()
		if err := cli.WriterByte(cli.ctx, data); err != nil {
			cli.log.Error("wsHub.subscribe WriterByte error: ", err)
		}
		/*
			if err := cli.WriterByte(ws.ctx, data); err != nil {
				ws.log.Error(fmt.Sprintf("wsHub.sender %s error: %s", id.String(), err))
				return
			}
		*/
	}
	ws.submit(unsub)
}

func (ws *wsHub) loop() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		case msg, ok := <-ws.inSubscrib:
			if !ok {
				continue
			}
			if msg.Type == models.SubscribeMsg {
				ws.subscribe(msg)
			} else {
				ws.unsubscribe(msg)
			}

		}
	}
}

func (ws *wsHub) writer() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		case m, ok := <-ws.inBroadcastCh:
			if !ok {
				continue
			}
			for id, _ := range ws.clients {
				ws.sender(id, m.Data)
			}
		case m, ok := <-ws.inSubscribCh:
			if !ok {
				continue
			}
			clients, ok := ws.subscribes[m.Subscribe]
			if !ok {
				continue
			}
			for cliUuid, _ := range clients {
				ws.subscribeSender(cliUuid, m.Subscribe, m.Data)
			}
		case m, ok := <-ws.inClientCh:
			if !ok {
				continue
			}
			ws.sender(m.Uuid, m.Data)
		}
	}
}
