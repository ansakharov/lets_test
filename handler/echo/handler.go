package echo_handler

import (
	"net/http"
)

// Handler prints request.
func Handler(greet string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(greet + r.URL.Query().Encode()))
	}
	return http.HandlerFunc(fn)
}
