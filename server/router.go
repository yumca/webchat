package server

import (
	"reflect"
	"sync"
)

type Routes struct {
	Handlers map[string]RouterHandler
	Lock     sync.Mutex // a lock for the map
}

type RouterHandler struct {
	ServerType  string
	ControlName string
	CallBack    reflect.Value
}

func NewRoutes() *Routes {
	r := new(Routes)
	r.Handlers = make(map[string]RouterHandler)
	r.Lock = sync.Mutex{}
	return r
}

func (r *Routes) Router(name string, serverType string, fn interface{}) {
	r.Handlers[name] = RouterHandler{
		serverType, name, reflect.ValueOf(fn),
	}
}

func (r *Routes) PassedArgs(callback *RouterHandler, args ...interface{}) []reflect.Value {
	funcType := callback.CallBack.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}
