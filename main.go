package main

import (
	"Kaban/Controller"
	"Kaban/Service/Handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	Handlers.SwapKeys()

	router := mux.NewRouter()
	api := mux.NewRouter()
	router.Use(Controller.LoggingRequest)
	api.Use(Controller.CheckJwtTokenLifyTime)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "Service/Fronted/Maine.html")
		}

	})

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			fmt.Println("Got a tick", t)
			Handlers.SwapKeys()
		}
	}()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "Service/Fronted/login.html")

	})
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Service/Fronted/Register.html")
	})
	router.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "Service/Fronted/Main_Page.html")

	})
	router.HandleFunc("/wait", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "Service/Fronted/WaitDownload.html")

	})
	router.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "Service/Fronted/Protecion.html")

	})

	router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "Service/Fronted/UrlFronted.html")

	}).Name("fileName")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла", err)

	}

	router.HandleFunc("/login/api", Controller.Loging).Methods("POST")
	router.HandleFunc("/register/api", Controller.Register).Methods("POST")

	router.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller.DownloadWithEncrypt(writer, request)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)
	router.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller.DownloadWithNotEncrypt(writer, request)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)

	router.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller.FileUploaderNoEncrypt(writer, request, router)

	}).Methods(http.MethodPost)
	router.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller.FileUploaderEncrypt(writer, request, router)

	}).Methods(http.MethodPost)
	router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller.GetFrom(writer, request)

	}).Methods("GET")
	router.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller.CUrlUp(writer, request)

	}).Methods(http.MethodGet)

	server := http.Server{
		Addr:                         ":8080",
		Handler:                      router,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            6 * time.Second,
		WriteTimeout:                 0,
		IdleTimeout:                  60 * time.Second,
		MaxHeaderBytes:               1 << 20,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     slog.NewLogLogger(nil, slog.LevelInfo),
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
	}

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Err cant' do this", "err", err)
		return
	} else {
		slog.Info("Start")
	}

}
