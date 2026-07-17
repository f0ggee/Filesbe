package cmds

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetAboutProjectUrlRouter(getRequest *mux.Router) *mux.Route {
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

func GetLoginPageRouter(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/Service/Fronted/Login.html")
	})
}

func SetRobotsRouter(router *mux.Router) *mux.Route {
	return router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pkg/robots.txt")
	})
}
func GetInformationPageRouter(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/informationPage", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/InformationPage.html")

	}).Name("NameFile")
}

func GetRegisterPageRouter(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/Service/Fronted/Register.html")
	})
}
func GetMainPageRouter(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Main_Page.html")

	})
}
func GetSitemapRouter(router *mux.Router) *mux.Route {
	return router.HandleFunc("/sitemap.xml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/sitemap.xml")

	})
}
func GetProtectPageRouter(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/Protecion.html")

	})
}
func GetUrlPageRouter(router *mux.Router) *mux.Route {
	return router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/UrlFronted.html")

	}).Name("fileName")
}
