package cmds

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRouters(r *mux.Router, controllers *Controllers) {
	apiRouters(r, *controllers)
	userRouters(r)

}
func apiRouters(r *mux.Router, c Controllers) {
	postRouter := r.Methods(http.MethodPost).Subrouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()

	postRouter.HandleFunc("/login/api", func(writer http.ResponseWriter, request *http.Request) {
		c.login.LoginController(writer, request)
	})

	postRouter.HandleFunc("/register/api", func(writer http.ResponseWriter, request *http.Request) {
		c.register.Register(writer, request)
	})
	postRouter.HandleFunc("/upload/api", func(writer http.ResponseWriter, request *http.Request) {
		c.uploader.FileUploader(writer, request)
	})
	postRouter.HandleFunc("/uploadEncrypt/api", func(writer http.ResponseWriter, request *http.Request) {
		c.encryptUpload.FileUploaderEncrypt(writer, request)
	})
	getRouter.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {
		c.urlBuild.SetUrl(writer, request)
	})

	getRouter.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {
		c.download.Download(writer, request)
	})

	getRouter.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {
		c.encryptDownload.DownloadEncrypt(writer, request)
	})
}

func userRouters(r *mux.Router) {
	getRouter := r.Methods(http.MethodGet).Subrouter()
	staticRouters := r.PathPrefix("/static").Subrouter()
	getRouter.HandleFunc("/aboutProject", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/InfoPageAboutApp.html")
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "internal/Service/Fronted/Maine.html")
		}
	})

	staticRouters.HandleFunc("/favicon.png", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "/internal/Service/Fronted/favicon.png")
	})

	getRouter.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/Login.html")
	})
	staticRouters.HandleFunc("/robots.txt", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "/pcf/robots.txt")
	})

	getRouter.HandleFunc("/informationPage", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/InformationPage.html")

	})
	getRouter.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Register.html")

	})
	getRouter.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Main_Page.html")
	})
	staticRouters.HandleFunc("/sitemap.xml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/sitemap.xml")
	})
	getRouter.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Protecion.html")
	})
	getRouter.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/UrlFronted.html")

	}).Name("fileName")
}
