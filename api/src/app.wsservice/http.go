package app_wsservice

import (
	models "app.models"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"nhooyr.io/websocket"
)

func (ws *WsService) newRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", ws.serveFiles)
	r.Get("/ws", ws.wsUpgrade)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", ws.apiInfo)
		r.Get("/ws", ws.wsUpgrade)
	})

	return r
}

func (ws *WsService) apiInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Api Page")
}

func (ws *WsService) wsUpgrade(w http.ResponseWriter, r *http.Request) {
	var wscli *wsClient
	var err error
	ctx, cancel := context.WithCancel(r.Context())

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		ws.log.Error("WsService.wsUpgrade error: ", err)
	}
	if wscli, err = ws.wshub.registration(conn); err != nil {
		ws.log.Error("WsService.wsUpgrade registration error: ", err)
		return
	}
	defer func() {
		fmt.Println("WsService.wsUpgrade defer...")
		cancel()
		if err := conn.Close(websocket.StatusInternalError, "connection closed."); err != nil {
			ws.log.Error("WsService.wsUpgrade close websocketerror: ", err)
		}
		if err := ws.wshub.unregistration(wscli); err != nil {
			ws.log.Error("WsService.wsUpgrade unregistration error: ", err)
		}
	}()
	msg := models.NewResponse()
	msg.Method = "RegistrationDesktop"
	data, _ := msg.MarshalJSON()

	if err := wscli.WriterByte(ctx, data); err != nil {
		ws.log.Error("WsService.wsUpgrade WriterByte error : ", err)
	}
	wscli.Receive()
	fmt.Println("WsService.wsUpgrade finish...")
}

func (ws *WsService) serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	var path string
	if r.URL.Path == "/" {
		path = fmt.Sprintf("%s/%s", ws.config.HttpConfig.StaticPath, "index.html")
	} else {
		path = fmt.Sprintf("%s/%s", ws.config.HttpConfig.StaticPath, r.URL.Path)
	}
	http.ServeFile(w, r, path)
}
