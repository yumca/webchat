package server

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	Routes *Routes
	// WsServLinks *WsLinks
}

// http server
// NewHttpServer creates a new Server object.
func NewHttpServer(r *Routes, nws *NetworkServer) {
	hs := &HttpServer{
		r,
	}
	NewControl("http", r)
	for k, h := range r.Handlers {
		if h.ServerType == "http" {
			http.HandleFunc(k, hs.httpHandleFunc)
		}
	}
	fmt.Println("HttpServer Start")
	nws.started = true
	go http.ListenAndServe(nws.address, nil)
}

func (hs *HttpServer) httpHandleFunc(w http.ResponseWriter, r *http.Request) {
	routerHandler, ok := hs.Routes.Handlers[r.RequestURI]
	if ok {
		routerHandler.CallBackHttp(w, r)
	}
}
