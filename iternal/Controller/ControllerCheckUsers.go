package Controller

import (
	"Kaban/iternal/Service/Handlers"
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

func GetFrom(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {

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
	rtToken, ok := seSession.Values["RT"].(string)
	if !ok {
		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/login"}); err != nil {
			slog.Error("Error decode the json", "Err", err)
			return
		}
		return
	}
	jwts, _ := seSession.Values["JWT"].(string)
	slog.Info("JWT", "JWT", jwts)

	NewJwt, err := s.Auth(rtToken, jwts)
	if err != nil {

		w.Header().Set(ContentType, JsonExample)
		w.WriteHeader(http.StatusUnauthorized)

		if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/login"}); err != nil {
			slog.Error("Error decode the json", "Err", err)
			return
		}
		return
	}
	seSession.Values["JWT"] = NewJwt
	w.Header().Set(ContentType, JsonExample)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/main"}); err != nil {
		slog.Error("Error decode the json", "Err", err)
		return
	}
	return
}
