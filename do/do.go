// Package do contains Task types and methods
package do

import (
	"errors"
	"net/http"

	"github.com/zeroibot/pack/dict"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/root"
	"github.com/zeroibot/pack/sys"
	"github.com/zeroibot/pack/web"
)

// Note: we use a global My instance for this package, because it makes setting up the task more convenient
// The initial implementation injected the App instance, but this makes 3 packages (do, root, web) be aware of an App type,
// and it lead to more complex code. By using the global My instance for this package, it simplifies the task handlers.
var myInstance *my.Instance = nil
var errNoMyInstance = errors.New("no My instance")

type DataFn[T any] = func(*my.Request) (T, error)
type ActionFn = func(*my.Request) error

type CmdParamsFn = func(*my.Request, []string) error
type WebParamsFn = func(*my.Request, *http.Request) error

type Data[T any] struct {
	Name string
	Fn   DataFn[T]
	Cmd  CmdParamsFn
	Web  WebParamsFn
}

type Action struct {
	Name string
	Fn   ActionFn
	Cmd  CmdParamsFn
	Web  WebParamsFn
}

type ForkData[T any] struct {
	Name   string
	Fork   map[string]DataFn[T]
	WebKey func(*http.Request) string
	Cmd    CmdParamsFn
	Web    WebParamsFn
}

type ForkAction struct {
	Name   string
	Fork   map[string]ActionFn
	WebKey func(*http.Request) string
	Cmd    CmdParamsFn
	Web    WebParamsFn
}

// SetMyInstance sets the My instance
func SetMyInstance(i *my.Instance) {
	myInstance = i
}

// CmdHandler returns a Data root command handler
func (t Data[T]) CmdHandler() root.CmdHandler {
	return func(params []string) {
		rq, err := cmdBasic(t.Name, params, t.Cmd)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		data, err := t.Fn(rq)
		sys.DisplayData(rq, data, err)
	}
}

// CmdHandler returns a ForkData root command handler
func (t ForkData[T]) CmdHandler() root.CmdHandler {
	return func(params []string) {
		rq, key, err := cmdFork(t.Name, params, t.Cmd)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		fn, ok := t.Fork[key]
		if !ok {
			sys.DisplayError(fail.InvalidOption)
			return
		}
		data, err := fn(rq)
		sys.DisplayData(rq, data, err)
	}
}

// CmdHandler returns an Action root command handler
func (t Action) CmdHandler() root.CmdHandler {
	return func(params []string) {
		rq, err := cmdBasic(t.Name, params, t.Cmd)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		err = t.Fn(rq)
		sys.DisplayOutput(rq, err)
	}
}

// CmdHandler returns a ForkAction root command handler
func (t ForkAction) CmdHandler() root.CmdHandler {
	return func(params []string) {
		rq, key, err := cmdFork(t.Name, params, t.Cmd)
		if err != nil {
			sys.DisplayError(err)
			return
		}
		fn, ok := t.Fork[key]
		if !ok {
			sys.DisplayError(fail.InvalidOption)
			return
		}
		err = fn(rq)
		sys.DisplayOutput(rq, err)
	}
}

// WebHandler returns a Data web handler
func (t Data[T]) WebHandler() web.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if myInstance == nil {
			web.SendError(w, nil, errNoMyInstance)
			return
		}
		rq, err := myInstance.NewRequest(t.Name)
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

// WebHandler returns a ForkData web handler
func (t ForkData[T]) WebHandler() web.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if myInstance == nil {
			web.SendError(w, nil, errNoMyInstance)
			return
		}
		rq, err := myInstance.NewRequest(t.Name)
		if err != nil {
			web.SendError(w, rq, err)
			return
		}
		key := t.WebKey(r)
		if dict.NoKey(t.Fork, key) {
			web.SendError(w, rq, fail.InvalidOption)
			return
		}
		if t.Web != nil {
			err = t.Web(rq, r)
			if err != nil {
				web.SendError(w, rq, err)
				return
			}
		}
		fn := t.Fork[key]
		data, err := fn(rq)
		web.SendDataResponse(w, rq, data, err)
	}
}

// WebHandler returns an Action web handler
func (t Action) WebHandler() web.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if myInstance == nil {
			web.SendError(w, nil, errNoMyInstance)
			return
		}
		rq, err := myInstance.NewRequest(t.Name)
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

// Common: basic command handlers
func cmdBasic(name string, params []string, cmdParams CmdParamsFn) (*my.Request, error) {
	if myInstance == nil {
		return nil, errNoMyInstance
	}
	rq, err := myInstance.NewRequest(name)
	if err != nil {
		return rq, err
	}
	if cmdParams != nil {
		err = cmdParams(rq, params)
		if err != nil {
			return rq, err
		}
	}
	return rq, nil
}

// Common: fork command handlers
func cmdFork(name string, params []string, cmdParams CmdParamsFn) (*my.Request, string, error) {
	if myInstance == nil {
		return nil, "", errNoMyInstance
	}
	rq, err := myInstance.NewRequest(name)
	if err != nil {
		return rq, "", err
	}
	key, params := params[0], params[1:]
	if cmdParams != nil {
		err = cmdParams(rq, params)
		if err != nil {
			return rq, "", err
		}
	}
	return rq, key, nil
}
