// Package web contains web server functions and types
package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeroibot/pack/clock"
	"github.com/zeroibot/pack/fail"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/str"
)

const okMessage string = "OK"

type Response struct {
	Data any
	OK   bool
	Msg  string
}

// Heartbeat is a health check web endpoint
func Heartbeat(w http.ResponseWriter, _ *http.Request) {
	SendJSON(w, http.StatusOK, Response{clock.DateTimeNow(), true, okMessage})
}

// SendJSON sends a JSON response to the client
func SendJSON[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", my.Err500)
	}
}

// SendError sends a Response to the client indicating an error occurred
func SendError(w http.ResponseWriter, rq *my.Request, err error) {
	displayOutput(rq, err)
	sendErrorResponse(w, rq, err)
}

// SendActionResponse sends a Response to the client indicating whether the action was successful
func SendActionResponse(w http.ResponseWriter, rq *my.Request, err error) {
	displayOutput(rq, err)
	if err == nil {
		SendJSON(w, rq.Status, Response{nil, true, okMessage})
	} else {
		sendErrorResponse(w, rq, err)
	}
}

// SendDataResponse sends a data Response to the client
func SendDataResponse(w http.ResponseWriter, rq *my.Request, data any, err error) {
	displayOutput(rq, err)
	if err == nil {
		SendJSON(w, rq.Status, Response{data, true, okMessage})
	} else {
		sendErrorResponse(w, rq, err)
	}
}

// sendErrorResponse sends a Response to the client indicating an error occurred
func sendErrorResponse(w http.ResponseWriter, rq *my.Request, err error) {
	status, message := getStatusMessage(rq, err)
	SendJSON(w, status, Response{nil, false, message})
}

// displayOutput displays the request output and error message, if applicable
func displayOutput(rq *my.Request, err error) {
	out := str.NewBuilder()
	if rq != nil {
		out.Add(rq.Output())
	}
	if err != nil {
		out.AddFmt("Error: %s", err.Error())
	}
	output := out.Build("\n")
	if output != "" {
		fmt.Println("Output: " + output)
	}
}

// getStatusMessage returns the status code and error message
func getStatusMessage(rq *my.Request, err error) (int, string) {
	message, _ := fail.PublicMessage(err)
	status := my.Err400
	if rq != nil {
		status = rq.Status
	}
	if errors.Is(err, fail.MissingParams) {
		status = my.Err400
	} else if errors.Is(err, fail.MissingSession) {
		status = my.Err401
	} else if errors.Is(err, fail.NotAuthorized) {
		status = my.Err403
	} else if errors.Is(err, fail.NotFoundItem) {
		status = my.Err404
	}
	return status, message
}
