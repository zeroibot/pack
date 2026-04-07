package sys

import (
	"fmt"

	"github.com/zeroibot/pack/ds"
	"github.com/zeroibot/pack/my"
	"github.com/zeroibot/pack/str"
)

// DisplayOutput prints request logs and error; prints OK if no error
func DisplayOutput(rq *my.Request, err error) {
	if rq != nil {
		displayOutput(rq)
	}
	if err == nil {
		fmt.Println(okMessage)
	} else {
		DisplayError(err)
	}
}

// DisplayResult prints request logs and error
func DisplayResult(rq *my.Request, err error) {
	if rq != nil {
		displayOutput(rq)
	}
	if err != nil {
		DisplayError(err)
	}
}

// DisplayError prints error
func DisplayError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// DisplayData prints data, request logs, and error
func DisplayData[T any](rq *my.Request, result ds.Result[T]) {
	if rq != nil {
		displayOutput(rq)
	}
	if result.IsError() {
		DisplayError(result.Error())
	} else {
		output, err := str.IndentedJSON(result.Value(), 2)
		if err == nil {
			fmt.Println(output)
		} else {
			DisplayError(err)
		}
	}
}

// DisplayList prints list items, request logs, and error
func DisplayList[T any](rq *my.Request, result ds.Result[[]T]) {
	if rq != nil {
		displayOutput(rq)
	}
	if result.IsError() {
		DisplayError(result.Error())
	} else {
		items := result.Value()
		for i, item := range items {
			fmt.Printf("%d: %v\n", i+1, item)
		}
		fmt.Println("Count:", len(items))
	}
}

// Common: displays the output of the request
func displayOutput(rq *my.Request) {
	output := rq.Output()
	if output != "" {
		fmt.Println(output)
	}
}
