package logger

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/sysmanage-web/core/logger"
)

func loadLoggerApi(r chi.Router) {
	r.Post("/getLogEntry", func(w http.ResponseWriter, r *http.Request) {
		// Fetch from logger.LogMap
		console := logger.LogMap.Get(r.URL.Query().Get("id"))

		if console.IsDone {
			console.LastLog = append(console.LastLog, "\n\n=====\nTask completed successfully")
			w.Header().Set("X-Is-Done", "1")
		}

		w.Header().Set("X-Last-Updated", console.LastUpdate.Format(time.RFC3339))

		bytes, err := json.Marshal(console.LastLog)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal log entry."))
			return
		}

		w.Write([]byte(bytes))
	})
}
