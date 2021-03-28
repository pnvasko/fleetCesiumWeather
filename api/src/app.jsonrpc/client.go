package app_jsonrpc

import (
	base "app.base"
	models "app.models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type AsyncResult struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type JsonRpcClient struct {
	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	PendingRequests int
	FlushDelay      time.Duration
	RequestTimeout  time.Duration

	requestsChan    chan *models.RpcJsonRequest
	pendingRequests map[string]map[string]chan *models.RpcJsonRequest

	log *base.Logger
}

func NewJsonRpcClient(ctx context.Context) *JsonRpcClient {
	cli := &JsonRpcClient{
		PendingRequests: DefaultPendingMessages,
		FlushDelay:      DefaultFlushDelay,
		RequestTimeout:  DefaultRequestTimeout,
	}

	cli.ctx, cli.cancel = context.WithCancel(ctx)
	cli.log = base.GetLog(cli.ctx, cli)
	cli.requestsChan = make(chan *models.RpcJsonRequest, cli.PendingRequests)
	cli.pendingRequests = make(map[string]map[string]chan *models.RpcJsonRequest, cli.PendingRequests)

	return cli
}

func (cli *JsonRpcClient) GetClientMsgId() string {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	return uuid.Must(uuid.NewRandom()).String()
}

func (cli *JsonRpcClient) addToPendingRequests(uuid, msgId string) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	_, ok := cli.pendingRequests[uuid]
	if !ok {
		cli.pendingRequests[uuid] = make(map[string]chan *models.RpcJsonRequest)
	}
	_, ok = cli.pendingRequests[uuid][msgId]
	if !ok {
		cli.pendingRequests[uuid][msgId] = make(chan *models.RpcJsonRequest)
	}
}

func (cli *JsonRpcClient) deleteFromPendingRequests(uuid, msgId string) {
	cli.mu.Lock()
	defer func() {
		if len(cli.pendingRequests[uuid]) == 0 {
			delete(cli.pendingRequests, uuid)
		}
		cli.mu.Unlock()
	}()

	_, ok := cli.pendingRequests[uuid]
	if !ok {
		return
	}

	_, ok = cli.pendingRequests[uuid][msgId]
	if !ok {
		return
	}
	close(cli.pendingRequests[uuid][msgId])
	delete(cli.pendingRequests[uuid], msgId)
}

func (cli *JsonRpcClient) putToPendingRequests(msg *models.RpcJsonRequest) {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	uuidWsClients := msg.UuidWsClients.String()

	_, ok := cli.pendingRequests[uuidWsClients]
	if !ok {
		return
	}

	ch, ok := cli.pendingRequests[uuidWsClients][msg.Id]
	if !ok {
		return
	}
	ch <- msg
}

func (cli *JsonRpcClient) publish(ch string, msg interface{}) error {
	cli.log.Debug("JsonRpcClient.publish ch:", ch)
	cli.log.Debug("JsonRpcClient.publish msg:", msg)
	return nil
}

func (cli *JsonRpcClient) eventProcessing(msg *models.RpcJsonRequest) {
	cli.log.Debug("JsonRpcClient.eventProcessing:")
	cli.requestsChan <- msg
}

func (cli *JsonRpcClient) RequestsChan() <-chan *models.RpcJsonRequest {
	return cli.requestsChan
}

func (cli *JsonRpcClient) Loop() {
	for {
		select {
		case <-cli.ctx.Done():
			return
		case msg, ok := <-cli.requestsChan:
			if ok {
				cli.putToPendingRequests(msg)
			}
		}
	}
}

func (cli *JsonRpcClient) closePendingRequests() error {
	for _, cli := range cli.pendingRequests {
		for _, req := range cli {
			close(req)
		}
	}

	return nil
}

func (cli *JsonRpcClient) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	cli.cancel()
	close(cli.requestsChan)
	if err := cli.closePendingRequests(); err != nil {
		cli.log.Error("JsonRpcClient.Stop error close pendingRequests chains")
	}
}

func (cli *JsonRpcClient) Call(request interface{}) (*models.RpcJsonRequest, error, error) {
	return cli.CallTimeout(request, cli.RequestTimeout, false)
}

func (cli *JsonRpcClient) Send(request interface{}) (*models.RpcJsonRequest, error, error) {
	return cli.CallTimeout(request, cli.RequestTimeout, true)
}

func (cli *JsonRpcClient) CallTimeout(request interface{}, timeout time.Duration, skipResponse bool) (*models.RpcJsonRequest, error, error) {
	return cli.callAsync(request, timeout, false, false)
}

func (cli *JsonRpcClient) callAsync(request interface{}, timeout time.Duration, skipResponse bool, usePool bool) (*models.RpcJsonRequest, error, error) {
	ctx, cancel := context.WithTimeout(cli.ctx, timeout)
	req, ok := request.(*models.RpcJsonRequest)
	if !ok {
		return nil, fmt.Errorf("RpcClient.callAsync unknow request type"), fmt.Errorf("unknow request type")
	}
	// req.ReturnCh = cli.inChName
	data, err := req.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("RpcClient.callAsync marshal error: ", err), fmt.Errorf("prepare function call")
	}
	uuidWsClients := req.UuidWsClients.String()
	cli.addToPendingRequests(uuidWsClients, req.Id)
	defer func() {
		cli.deleteFromPendingRequests(uuidWsClients, req.Id)
	}()
	fmt.Println("JsonRpcClient.callAsync data:", string(data))
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				return nil, err, fmt.Errorf("function call deadline")
			}
			return nil, fmt.Errorf("JsonRpcClient.callAsync unknow error"), fmt.Errorf("unknow")
		case msg, ok := <-cli.pendingRequests[uuidWsClients][req.Id]:
			cancel()
			if ok {
				return msg, nil, nil
			}
		}
	}
}
