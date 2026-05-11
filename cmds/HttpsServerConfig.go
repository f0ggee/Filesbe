package cmds

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func ServerConfig(r *mux.Router) *http.Server {

	Port := os.Getenv("PORT")
	slog.Info("Our new port", "Port", Port)
	server := http.Server{
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            6 * time.Second,
		WriteTimeout:                 0,
		IdleTimeout:                  60 * time.Second,
		MaxHeaderBytes:               1 << 20,
	}

	return &server
}
