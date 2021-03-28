package app_jsonrpc

/*
import (
	models "app.models"
	"context"
	"nhooyr.io/websocket"
	"sync"
)

func (s *JsonRpcServer) Handle(ctx context.Context, conn *websocket.Conn) {
	var wg sync.WaitGroup
	var err error
	basectx, basecancel := context.WithCancel(ctx)

	readerStop := make(chan struct{})
	writerStop := make(chan struct{})

	readerqueue := make(chan interface{})
	writerqueue := make(chan interface{})
	responses := make(chan *models.RpcJsonResponse)

	defer func() {
		wg.Wait()
		defer basecancel()
		close(responses)
		close(writerqueue)
		close(readerqueue)
	}()


}
*/
