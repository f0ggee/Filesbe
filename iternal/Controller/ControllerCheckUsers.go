package Controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		formatForOutput := "01/02/2006 03:04 PM"

		t := time.Now()

		if request.Header.Get("User-Agent") == "" {
			return
		}

		if request != nil {
			slog.String("Method ", request.Method)
			slog.String("Time", t.Format(formatForOutput))
			slog.String("URL", request.URL.String())
			slog.String("Host", request.Host)
		}

		next.ServeHTTP(writer, request)
	})
}

func GetFrom(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Cant' treat", http.StatusNotFound)
		slog.Info("Not found")
		return
	}
	type AnswerStruct struct {
		StatusRedict string `json:"status_redict"`
	}

	store := Store()
	seSession, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("Error check", "Err", err)

		return

	}
	rtToken, _ := seSession.Values["RT"].(string)

	jwts, _ := seSession.Values["JWT"].(string)
	_, _, err, ok := Auth(rtToken, jwts, seSession)

	if err != nil {

		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/login"}); err != nil {
			slog.Error("Error decode the json", "Err", err)
			return
		}
	}
	if !ok {
		return
	}

	w.Header().Set(ContentType, JsonExample)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/main"}); err != nil {
		slog.Error("Error decode the json", "Err", err)
		return
	}
	return
}
