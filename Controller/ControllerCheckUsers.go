package Controller

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		formatForOutput := "01/02/2006 03:04 PM"

		t := time.Now()
		slog.Info("Info",
			slog.Group("Request"),
			slog.String("URL", request.URL.Path),
			slog.String("Time", t.Format(formatForOutput)),
			slog.String("Ip", request.RemoteAddr),
		)
		next.ServeHTTP(writer, request)
	})
}

func CheckJwtTokenLifyTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(writer, request)

	})

}

func GetFrom(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Cant' treat", http.StatusNotFound)
		slog.Info("Not found")
		return
	}

	store := Store()
	seSession, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("Error check", "Err", err)

		return

	}
	rtToken, _ := seSession.Values["RT"].(string)

	jwts, _ := seSession.Values["JWT"].(string)
	s1, s2, err, ok := Auth(rtToken, jwts, seSession)
	slog.Info("Data", s1, s2)

	if err != nil {
		slog.Error("as")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}
	if !ok {
		return
	}

	http.Redirect(w, r, "/main", http.StatusMovedPermanently)

	clientIP := r.RemoteAddr

	host, _, err := net.SplitHostPort(clientIP)
	if err != nil {
		slog.Info("ee")
		return
	}
	clientIP = host

}
