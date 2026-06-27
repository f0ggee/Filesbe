package SessionHandle

import (
	"Kaban/internal/DomainLevel"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/sessions"
)

type SessionConnect struct {
	Activity *sessions.CookieStore
	mut      *sync.RWMutex
}

func (s *SessionConnect) getUserConnect(r *http.Request) (*sessions.Session, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.Activity.Get(r, DomainLevel.TokenName)
}

var SessionControl = &SessionConnect{
	Activity: func() *sessions.CookieStore {

		var store1z, err = hex.DecodeString(os.Getenv("KEY1"))
		if err != nil {
			slog.Error("SessionControl: error to create a cookie", "Err", err)
			panic(err)
		}
		Store := sessions.NewCookieStore(store1z)
		return Store

	}(),
	mut: &sync.RWMutex{},
}
