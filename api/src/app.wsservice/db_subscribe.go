package app_wsservice

import (
	base "app.base"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/pgxpool"
	"os"
	"strings"
	"sync"
)

func quoteIdentifier(s string) string {
	return `"` + strings.Replace(s, `"`, `""`, -1) + `"`
}

type postgresPubSubContext struct {
	Msg *pgconn.Notification
}
type postgresSubscribes struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	channels []string
	handlers map[string]func(ctx *postgresPubSubContext) error
	loopDone chan struct{}
	shutdown chan struct{}

	sendToEventBus func(channel string, msg []byte) error
	defaultHandler func(ctx *postgresPubSubContext) error

	pgxpool *pgxpool.Pool
	log     *base.Logger
}

func newPostgresSubscribe(ctx context.Context) (*postgresSubscribes, error) {
	var ok bool
	ps := &postgresSubscribes{
		handlers: make(map[string]func(ctx *postgresPubSubContext) error),
		shutdown: make(chan struct{}),
	}

	ps.ctx, ps.cancel = context.WithCancel(ctx)

	ps.log = base.GetLog(ps.ctx, ps)
	if ps.pgxpool, ok = ctx.Value("pgxpool").(*pgxpool.Pool); !ok {
		return nil, fmt.Errorf("error get pgxpool from context")
	}

	return ps, nil
}

func (ps *postgresSubscribes) addHandler(channel string, handler func(ctx *postgresPubSubContext) error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.handlers[channel]
	if !ok {
		ps.handlers[channel] = handler
	} else {
		ps.log.Error("Channel handler exist.")
	}
}

func (ps *postgresSubscribes) loop() {
	ps.log.Debug("postgresSubscribes.loop start...")
	conn, err := ps.pgxpool.Acquire(ps.ctx)
	if err != nil {
		ps.log.Error("postgresSubscribes.loop error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	for channel := range ps.handlers {
		fmt.Println("channel:", channel)
		if _, err = conn.Exec(ps.ctx, fmt.Sprintf("listen %s", quoteIdentifier(channel))); err != nil {
			ps.log.Error("postgresSubscribes.loop error subscribe channel:", err)
			os.Exit(1)
		}
		// if err := conn.Conn().Listen(channel); err != nil { }
	}
	ps.loopDone = make(chan struct{})
	subscribeCtx, subscribeCancel := context.WithCancel(ps.ctx)
notificationLook:
	for {
		select {
		case <-ps.ctx.Done():
			ps.log.Debug("postgresSubscribes.loop done")
			close(ps.loopDone)
			break notificationLook
		case <-ps.shutdown:
			ps.log.Debug("postgresSubscribes.loop shutdown")
			subscribeCancel()
			close(ps.loopDone)
			break notificationLook
		default:
			notification, err := conn.Conn().WaitForNotification(subscribeCtx)
			if err != nil {
				ps.log.Error("postgresSubscribes.loop error waiting notification:", err)
				continue
			}
			handler, ok := ps.handlers[notification.Channel]
			if !ok {
				ps.log.Error("postgresSubscribes.loop error get handler for notification channel: ", notification.Channel)
				continue
			}
			psCtx := &postgresPubSubContext{
				Msg: notification,
			}
			if err := handler(psCtx); err != nil {
				ps.log.Error("postgresSubscribes.loop handler error:", err)
			}
		}
	}
}

func (ps *postgresSubscribes) close(complete chan struct{}) {
	ps.log.Debug("postgresSubscribes.close start...")
	defer close(complete)

	close(ps.shutdown)
	<-ps.loopDone
	ps.cancel()

	conn, err := ps.pgxpool.Acquire(ps.ctx)
	if err != nil {
		ps.log.Error("postgresSubscribes.close error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	for channel := range ps.handlers {
		fmt.Println("channel,", channel)
		if _, err = conn.Exec(ps.ctx, fmt.Sprintf("unlisten %s", quoteIdentifier(channel))); err != nil {
			ps.log.Error("postgresSubscribes.close error unsubscribe channel:", err)
			os.Exit(1)
		}
	}
	ps.log.Debug("postgresSubscribes.close finish.")

}
