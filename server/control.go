package server

import (
	"chat/app/control"
	"chat/library"
	"net/http"
)

func NewControl(serviceType string, r *Routes) {
	conhttp := map[string]func(http.ResponseWriter, *http.Request){
		"/http/test": control.HttpTest,
	}
	contcp := map[string]func(http.ResponseWriter, *http.Request){
		"/tcp/test": control.TcpTest,
	}
	conws := map[string]func(library.WsConn, []byte) []byte{
		"/ws/test": control.WsTest,
	}

	switch serviceType {
	case "http":
		if len(conhttp) > 0 {
			for n, h := range conhttp {
				r.RouterHttp(n, serviceType, h)
			}
		}
	case "tcp":
		if len(contcp) > 0 {
			for n, h := range contcp {
				r.RouterTcp(n, serviceType, h)
			}
		}
	case "ws":
		if len(conws) > 0 {
			for n, h := range conws {
				r.RouterWs(n, serviceType, h)
			}
		}
	}

}
