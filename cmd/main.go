package main

import (
	Controller2 "Kaban/iternal/Controller"
	"Kaban/iternal/Service/Handlers"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	Handlers.SwapKeys()

	router := mux.NewRouter()
	api := mux.NewRouter()
	router.Use(Controller2.LoggingRequest)
	api.Use(Controller2.CheckJwtTokenLifeTime)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "iternal/Service/Fronted/Maine.html")
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

		http.ServeFile(w, r, "iternal/Service/Fronted/login.html")

	})
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "iternal/Service/Fronted/Register.html")
	})
	router.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "iternal/Service/Fronted/Main_Page.html")

	})
	router.HandleFunc("/wait", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/WaitDownload.html")

	})
	router.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/Protecion.html")

	})

	router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "iternal/Service/Fronted/UrlFronted.html")

	}).Name("fileName")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла", err)

	}

	router.HandleFunc("/login/api", Controller2.Loging).Methods("POST")
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
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     slog.NewLogLogger(nil, slog.LevelInfo),
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
	}

	//err = server.ListenAndServeTLS("/etc/letsencrypt/live/filesbes.com/fullchain.pem", "/etc/letsencrypt/live/filesbes.com/privkey.pem")
	//if err != nil {
	//	slog.Error("Err cant' do this", "err", err)
	//	return
	//}

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Server couldn't start", err)
		return

	}

}
