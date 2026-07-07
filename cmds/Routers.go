package cmds

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetAboutProjectUrl(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/aboutProject", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/InfoPageAboutApp.html")
	})
}
func GetDefaultRouter(router *mux.Router) *mux.Route {
	return router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "internal/Service/Fronted/Maine.html")
		}

	})
}

func GetPhotoRequest(StaticFiles *mux.Router) *mux.Route {
	return StaticFiles.Handle("/favicon.png", http.FileServer(http.Dir("internal/Service/Fronted/favicon.png")))
}

func GetLoginRequest(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/Service/Fronted/Login.html")
	})
}

func SetRobotsRequest(router *mux.Router) *mux.Route {
	return router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pkg/robots.txt")
	})
}
func GetInformationPage(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/informationPage", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/InformationPage.html")

	}).Name("NameFile")
}

func GetRegisterUrl(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/Service/Fronted/Register.html")
	})
}
func GetMainUrl(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Main_Page.html")

	})
}
func GetSitemap(router *mux.Router) *mux.Route {
	return router.HandleFunc("/sitemap.xml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/sitemap.xml")

	})
}
func GetProtectRequest(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/Protecion.html")

	})
}
func GetUrlRequest(router *mux.Router) *mux.Route {
	return router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/UrlFronted.html")

	}).Name("fileName")
}
