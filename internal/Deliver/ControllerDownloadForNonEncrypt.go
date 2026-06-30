package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func DownloadWithNotEncrypt(w http.ResponseWriter, r *http.Request, s *Application.HandlerPackCollect) {
	type JsonAnswer struct {
		StatusOperation string   `json:"StatusOperation"`
		Error           []string `json:"Error"`
		Url             string   `json:"Url"`
	}
	if r.Method != http.MethodGet {
		//TODO add handling the error
	}

	//TODO remove the line below
	name := getNameFromUrl(r)

	err, _ := s.DownloadWithNonEncrypt(w, name, r.Context())

	//TODO add handling the error
	return

}
