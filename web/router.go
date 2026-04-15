package web

import (
	"fmt"
	"net/http"
)

type Handler = http.HandlerFunc
type HandlerFn[A any] = func(*A) Handler
type HandlerMap[A any] = map[string]HandlerFn[A]
type HandlerTask[A any] = map[string]task[A]

type task[A any] interface {
	Web() HandlerFn[A]
}

// RegisterRoutes adds handlers to the router
func RegisterRoutes[A any](this *A, mux *http.ServeMux, verb string, handlers HandlerMap[A], middlewares ...Middleware) int {
	var middleware Middleware = nil
	if len(middlewares) > 0 {
		middleware = StackMiddlewares(middlewares...)
	}
	for path, handlerFn := range handlers {
		pattern := fmt.Sprintf("%s %s", verb, path)
		handler := handlerFn(this)
		if middleware != nil {
			mux.Handle(pattern, middleware(handler))
		} else {
			mux.HandleFunc(pattern, handler)
		}
	}
	return len(handlers)
}

// RegisterTasks adds handlers to the routers from the tasks
func RegisterTasks[A any](this *A, mux *http.ServeMux, verb string, tasks HandlerTask[A], middlewares ...Middleware) int {
	handlers := make(HandlerMap[A])
	for path, task := range tasks {
		handlers[path] = task.Web()
	}
	return RegisterRoutes(this, mux, verb, handlers, middlewares...)
}
