// handlers/handlers.go
package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, world"
	status := http.StatusOK

	w.WriteHeader(status)
	w.Write([]byte(msg))

	slog.Info(strconv.Itoa(status), "path", r.URL.Path, "msg", msg)
}
