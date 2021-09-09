package server

import (
	"chat/library"
	"sync"
)

type NetService struct {
	Routes  *Routes
	Service map[string]*NetworkServer
	wg      *sync.WaitGroup
}

type Server interface{}

type NetworkServer struct {
	ServiceType string
	Server      *Server
	started     bool
	address     string
}

func NewNetService(conf library.Config) (Service *NetService, err error) {
	Service = new(NetService)
	Service.Service = make(map[string]*NetworkServer)
	Service.wg = &sync.WaitGroup{}
	Service = &NetService{
		NewRoutes(),
		make(map[string]*NetworkServer),
		&sync.WaitGroup{},
	}
	if len(conf.Server) > 0 {
		for k, addr := range conf.Server {
			Service.Register(k, addr)
		}
	}
	// if conf.Redis.Stat == "on" {
	// 	rdb, err := Rdlink()
	// 	if err == nil {
	// 		Serv.RdbWorkerKey = conf.Redis.RedisPrefix + conf.Setting.WorkerName + ":Worker"
	// 		Serv.rdb = rdb
	// 	}
	// }
	return
}

// 注册服务 开始程序
func (s *NetService) Register(serviceType, address string) {
	// defer s.wg.Done()
	nsr := &NetworkServer{
		serviceType,
		nil,
		false,
		address,
	}
	s.Service[serviceType] = nsr
}

func (s *NetService) Start() {
	if len(s.Service) > 0 {
		for k, _ := range s.Service {
			s.NewServer(k)
		}
	}
}

func (s *NetService) NewServer(ServType string) {
	switch ServType {
	case "http":
		NewHttpServer(s.Routes, s.Service[ServType])
	case "tcp":
		NewTcpServer(s.Routes, s.Service[ServType])
	case "ws":
		NewWsServer(s.Routes, s.Service[ServType])
	}
}
func (s *NetService) Stop() {
}
