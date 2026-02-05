package main

import (
	Controller2 "Kaban/iternal/Controller"
	"Kaban/iternal/Service/Handlers"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//Once create the pair of keys
	Handlers.SwapKeys()

	router := mux.NewRouter()

	//router = router.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
	//	if request.Host != "filesbes.com" {
	//		return false
	//	}
	//	slog.Info(request.Host)
	//
	//	return true
	//}).Subrouter()

	//The router will return  static files
	StaticFiles := router.PathPrefix("/Fronted").Subrouter()

	router.HandleFunc("/aboutProject", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "iternal/Service/Fronted/InfoPageAboutApp.html")

	})

	StaticFiles.Handle("/favicon.png", http.FileServer(http.Dir("iternal/Service")))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "iternal/Service/Fronted/Maine.html")
		}

	})

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			slog.Info("Got a ticker", t)
			Handlers.SwapKeys()
		}
	}()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "iternal/Service/Fronted/Login.html")

	})

	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./robots.txt")

	})

	router.HandleFunc("/informationPage", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/InformationPage.html")

	}).Name("NameFile")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "iternal/Service/Fronted/Register.html")
	})
	router.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "iternal/Service/Fronted/Main_Page.html")

	})
	router.HandleFunc("/sitemap.xml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/sitemap.xml")

	})

	router.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/Protecion.html")

	})

	router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/UrlFronted.html")

	}).Name("fileName")

	router.HandleFunc("/login/api", Controller2.Login).Methods("POST")
	router.HandleFunc("/register/api", Controller2.Register).Methods("POST")

	router.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithEncrypt(writer, request)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)
	router.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithNotEncrypt(writer, request)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)

	router.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderNoEncrypt(writer, request, router)

	}).Methods(http.MethodPost)
	router.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderEncrypt(writer, request, router)

	}).Methods(http.MethodPost)
	router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.GetFrom(writer, request)

	}).Methods("GET")
	router.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.CUrlUp(writer, request)

	}).Methods(http.MethodGet)

	//##
	server := http.Server{
		Addr:                         ":8080", // I must change on 443
		Handler:                      router,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            6 * time.Second,
		WriteTimeout:                 0,
		IdleTimeout:                  60 * time.Second,
		MaxHeaderBytes:               1 << 20,
	}

	//err := server.ListenAndServeTLS("/etc/letsencrypt/live/filesbes.com/fullchain.pem", "/etc/letsencrypt/live/filesbes.com/privkey.pem")
	//if err != nil {
	//	slog.Error("Err cant' do this", "err", err)
	//	return
	//}

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Server couldn't start", err)
		return

	}

}
