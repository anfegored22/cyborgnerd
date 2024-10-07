package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Msg, ClientMsg string
}

func (e *Error) Error() string {
	return e.Msg
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func ToStandar(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := h(w, r)
		if e != nil {
			fmt.Print("<p>Error</p>")
		}
	}
}
