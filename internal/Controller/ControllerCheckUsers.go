package Controller

import (
	"Kaban/internal/Service/Handlers"
	"encoding/json"
	"log/slog"
	"net/http"
)

func GetFrom(w http.ResponseWriter, r *http.Request, s *Handlers.HandlerPackCollect) {

	if r.Method != http.MethodGet {
		http.Error(w, "Cant' treat", http.StatusNotFound)
		slog.Info("Not found")
		return
	}
	type AnswerStruct struct {
		StatusRedict string `json:"status_redict"`
	}

	//store := SessionStore()
	seSession, err := store.Get(r, TokenName)
	if err != nil {
		slog.Error("Error check", "Err", err)
		return
	}
	rtToken, _ := seSession.Values[RTCookieName].(string)
	jwts, _ := seSession.Values[JwtCookieName].(string)

	NewJwt, err := s.Auth(rtToken, jwts)
	if err != nil {
		w.Header().Set(ContentType, Json)
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/login"}); err != nil {
			slog.Error("Error decode the json", "Err", err)
			return
		}
		return
	}
	if NewJwt != "" {
		seSession.Values[JwtCookieName] = NewJwt

	}
	w.Header().Set(ContentType, Json)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(AnswerStruct{StatusRedict: "/main"}); err != nil {
		slog.Error("Error decode the json", "Err", err)
		return
	}
	return
}
