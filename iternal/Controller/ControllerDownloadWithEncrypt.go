package Controller

import (
	"Kaban/iternal/Service/Handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	return name

}
func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Status method don't allow", http.StatusBadRequest)
		return
	}

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Error           string `json:"Error"`
		UrlToRedict     string `json:"UrlToRedict"`
	}
	name := getNameFromUrl(r)

	err := Handlers.DownloadEncrypt(w, r.Context(), name)
	if err != nil {
		http.Redirect(w, r, "/informationPage", http.StatusFound)
		return
	}

	return

}
