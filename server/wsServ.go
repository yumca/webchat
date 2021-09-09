package server

import (
	"chat/library"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	Routes      *Routes
	WsServLinks *WsLinks
}

type WsLinks struct {
	WsConnectFd   int
	WsConnections map[int]library.WsConn
	Lock          sync.Mutex // a lock for the map
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ws server
// NewWsServer creates a new Server object.
func NewWsServer(r *Routes, nws *NetworkServer) {
	wss := &WsServer{
		r,
		&WsLinks{
			WsConnectFd:   0,
			WsConnections: make(map[int]library.WsConn),
			Lock:          sync.Mutex{},
		},
	}
	NewControl("ws", r)
	for k, h := range r.Handlers {
		if h.ServerType == "ws" {
			http.HandleFunc(k, wss.wsHandleFunc)
		}
	}
	fmt.Println("WsServer Start")
	nws.started = true
	go http.ListenAndServe(nws.address, nil)
}

func (wss *WsServer) wsHandleFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(ws)
	if err != nil {
		log.Println(err)
		return
	}
	wsConn := wss.bindWsConnection(ws, w, r)
	routerHandler, ok := wss.Routes.Handlers[r.RequestURI]
	if !ok {
		return
	}
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1001) {
				delete(wss.WsServLinks.WsConnections, wsConn.Fd)
				log.Printf("客户端：%d 退出", wsConn.Fd)
			}
			log.Println(err)
			return
		}
		result := routerHandler.CallBackWs(wsConn, p)
		if err := ws.WriteMessage(messageType, result); err != nil {
			log.Println(err)
			return
		}
	}
}

func (wss *WsServer) bindWsConnection(ws *websocket.Conn, w http.ResponseWriter, r *http.Request) library.WsConn {
	wss.WsServLinks.Lock.Lock()
	wss.WsServLinks.WsConnectFd = wss.WsServLinks.WsConnectFd + 1
	fd := wss.WsServLinks.WsConnectFd
	wss.WsServLinks.Lock.Unlock()
	wss.WsServLinks.WsConnections[fd] = library.WsConn{
		Fd:             fd,
		Conn:           ws,
		ResponseWriter: w,
		Request:        r,
	}
	return wss.WsServLinks.WsConnections[fd]
}

func (wss *WsServer) GetWsServLinks(fd int) (ok bool, connection library.WsConn) {
	connection, ok = wss.WsServLinks.WsConnections[fd]
	return
}
