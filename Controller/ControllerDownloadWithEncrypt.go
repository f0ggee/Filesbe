package Controller

import (
	"Kaban/Service/Handlers"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func CheckDenyList(RT string) bool {

	DenyLists := DenyList
	if _, ok := DenyLists[RT]; !ok {
		return true
	}
	return false
}

func CookieGetInControllerDownloader(w http.ResponseWriter, r *http.Request) (string, error) {
	store := Store()

	session, err := store.Get(r, "token6")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return "", err
	}

	RT, ok := session.Values["RT"].(string)
	if !ok {
		slog.Info("Hoa")
		http.Redirect(w, r, "/login", http.StatusFound)

	}
	JWT, ok := session.Values["JWT"].(string)
	if !ok {
		slog.Info("Hoa")
		http.Redirect(w, r, "/login", http.StatusFound)

	}

	_, SC, err, _ := Auth(RT, JWT, session)
	if err != nil {
		slog.Error("Err in Auth server", "Err", err)
		return "", err
	}

	return SC, err
}
func getNameFromUrl(r *http.Request) string {
	vars := mux.Vars(r)

	name := vars["name"]
	slog.Info(name)
	return name

}
func DownloadWithEncrypt(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Status method don't allow", http.StatusBadRequest)
		return
	}

	type Answer struct {
		StatusOperation string `json:"StatusOperation"`
		Errr            error  `json:"Errr"`
	}
	name := getNameFromUrl(r)

	_, er := CookieGetInControllerDownloader(w, r)
	if er != nil {

		if errors.Is(er, errors.New("token in deny list")) {
			slog.Error("Error because in deny list")
			w.Header().Set("Content-Type", "application/json")
			errs := json.NewEncoder(w).Encode(Answer{StatusOperation: "Don't start", Errr: errors.New("can't validate")})
			if errs != nil {
				slog.Error("error in decode json", errs)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return

		}
		slog.Error("Can't get cookie because", "err", er)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		err := json.NewEncoder(w).Encode(Answer{StatusOperation: "Don't start", Errr: errors.New("can't validate")})
		if err != nil {
			slog.Error("error in decode json", "err", err)
			return
		}

		return
	}

	err := Handlers.DownloadEncrypt(w, r.Context(), name)
	if err != nil {

		w.Header().Set("Content-Type", JAson)
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(Answer{
			StatusOperation: "BREAK",
			Errr:            errors.New("can't download file"),
		}); err != nil {
			slog.Error("Error Encode json", "err", err)
			return
		}
		return

	}

}
