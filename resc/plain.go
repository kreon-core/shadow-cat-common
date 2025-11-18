package resc

import (
	"net/http"

	"github.com/kreon-core/shadow-cat-common/logc"
)

func PlainText(w http.ResponseWriter, statusCode int, payload string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(payload)); err != nil {
		logc.Warn().Err(err).Msg("Failed to write plain text response")
		return
	}
}
