package server

import (
	"chat/library"
	"net/http"
	"sync"
)

type Routes struct {
	Handlers map[string]RouterHandler
	Lock     sync.Mutex // a lock for the map
}

type RouterHandler struct {
	ServerType   string
	ControlName  string
	CallBackHttp func(http.ResponseWriter, *http.Request)
	CallBackTcp  func(http.ResponseWriter, *http.Request)
	CallBackWs   func(library.WsConn, []byte) (p []byte)
	// CallBack     reflect.Value
}

func NewRoutes() *Routes {
	r := new(Routes)
	r.Handlers = make(map[string]RouterHandler)
	r.Lock = sync.Mutex{}
	return r
}

func (r *Routes) RouterHttp(name string, serverType string, fn func(http.ResponseWriter, *http.Request)) {
	r.Handlers[name] = RouterHandler{
		serverType, name, fn, nil, nil,
	}
}

func (r *Routes) RouterTcp(name string, serverType string, fn func(http.ResponseWriter, *http.Request)) {
	r.Handlers[name] = RouterHandler{
		serverType, name, nil, fn, nil,
	}
}

func (r *Routes) RouterWs(name string, serverType string, fn func(library.WsConn, []byte) []byte) {
	r.Handlers[name] = RouterHandler{
		serverType, name, nil, nil, fn,
	}
}
