package server

import (
	"chat/app/control"
)

func NewControl(serviceType string, r *Routes) {
	contr := map[string]map[string]interface{}{
		"http": {
			"/http/test": control.HttpTest,
		},
		"tcp": {
			"/tcp/test": control.TcpTest,
		},
		"ws": {
			"/ws/test": control.WsTest,
		},
	}

	if len(contr[serviceType]) > 0 {
		for n, h := range contr[serviceType] {
			r.Router(n, serviceType, h)
		}
	}

}
