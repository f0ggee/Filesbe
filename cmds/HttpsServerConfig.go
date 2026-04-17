package cmds

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ServerConfig(r *mux.Router) *http.Server {

	server := http.Server{
		Addr:                         ":443",
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
