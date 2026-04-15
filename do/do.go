// Package do contains Task types and methods
package do

import (
	"net/http"

	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/root"
	"github.com/zeroibot/pack/sys"
	"github.com/zeroibot/pack/web"
)

type DataFn[T any] = func(*my.Request) (T, error)
type ActionFn = func(*my.Request) error

type CmdParamsFn = func(*my.Request, []string) error
type WebParamsFn = func(*my.Request, *http.Request) error

type App interface {
	MyInstance() *my.Instance
}

type Data[T any, A App] struct {
	Name string
	Fn   DataFn[T]
	Cmd  CmdParamsFn
	Web  WebParamsFn
}

type Action[A App] struct {
	Name string
	Fn   ActionFn
	Cmd  CmdParamsFn
	Web  WebParamsFn
}

// CmdHandler returns a Data root command handler
func (t Data[T, A]) CmdHandler() root.CmdHandler[A] {
	return func(this A, params []string) {
		rq, err := this.MyInstance().NewRequest(t.Name)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		if t.Cmd != nil {
			err = t.Cmd(rq, params)
			if err != nil {
				sys.DisplayError(err)
				return
			}
		}
		data, err := t.Fn(rq)
		sys.DisplayData(rq, data, err)
	}
}

// CmdHandler returns an Action root command handler
func (t Action[A]) CmdHandler() root.CmdHandler[A] {
	return func(this A, params []string) {
		rq, err := this.MyInstance().NewRequest(t.Name)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		if t.Cmd != nil {
			err = t.Cmd(rq, params)
			if err != nil {
				sys.DisplayError(err)
				return
			}
		}
		err = t.Fn(rq)
		sys.DisplayOutput(rq, err)
	}
}

// WebHandler returns a Data web handler
func (t Data[T, A]) WebHandler() web.HandlerFn[A] {
	return func(this A) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			rq, err := this.MyInstance().NewRequest(t.Name)
			if err != nil {
				web.SendError(w, rq, err)
				return
			}
			if t.Web != nil {
				err = t.Web(rq, r)
				if err != nil {
					web.SendError(w, rq, err)
					return
				}
			}
			data, err := t.Fn(rq)
			web.SendDataResponse(w, rq, data, err)
		}
	}
}

// WebHandler returns an Action web handler
func (t Action[A]) WebHandler() web.HandlerFn[A] {
	return func(this A) web.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			rq, err := this.MyInstance().NewRequest(t.Name)
			if err != nil {
				web.SendError(w, rq, err)
				return
			}
			if t.Web != nil {
				err = t.Web(rq, r)
				if err != nil {
					web.SendError(w, rq, err)
					return
				}
			}
			err = t.Fn(rq)
			web.SendActionResponse(w, rq, err)
		}
	}
}
