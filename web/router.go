package web

import (
	"fmt"
	"net/http"
)

type Handler = http.HandlerFunc
type HandlerMap = map[string]Handler
type HandlerTask = map[string]webTask

type webTask interface {
	WebHandler() Handler
}

// RegisterRoutes adds handlers to the router
func RegisterRoutes(mux *http.ServeMux, verb string, handlers HandlerMap, middlewares ...Middleware) int {
	var middleware Middleware = nil
	if len(middlewares) == 1 {
		middleware = middlewares[0]
	} else if len(middlewares) > 1 {
		middleware = StackMiddlewares(middlewares...)
	}
	for path, handler := range handlers {
		pattern := fmt.Sprintf("%s %s", verb, path)
		if middleware != nil {
			mux.Handle(pattern, middleware(handler))
		} else {
			mux.HandleFunc(pattern, handler)
		}
	}
	return len(handlers)
}

// RegisterTasks adds handlers to the routers from the tasks
func RegisterTasks(mux *http.ServeMux, verb string, tasks HandlerTask, middlewares ...Middleware) int {
	handlers := make(HandlerMap)
	for path, task := range tasks {
		handlers[path] = task.WebHandler()
	}
	return RegisterRoutes(mux, verb, handlers, middlewares...)
}
