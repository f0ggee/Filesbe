package main

import (
	Controller2 "Kaban/iternal/Controller"
	"Kaban/iternal/InfrastructureLayer/s3Interation"
	"Kaban/iternal/Service/Connect_to_BD"
	"Kaban/iternal/Service/Handlers"
	"Kaban/iternal/Service/Helpers"
	"log/slog"
	"net/http"
	"os"
	"time"

	"Kaban/iternal/InfrastructureLayer/TokenInteraction"

	"Kaban/iternal/InfrastructureLayer/FileKeyInteration"
	"Kaban/iternal/InfrastructureLayer/RedisInteration"
	"Kaban/iternal/InfrastructureLayer/TokenInteraction/manageTokensImpl"
	"Kaban/iternal/InfrastructureLayer/UserInteraction"
	"Kaban/iternal/InfrastructureLayer/encryptionKeyInteration"

	"github.com/awnumar/memguard"
	"github.com/gorilla/mux"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	handler := slog.New(slog.NewTextHandler(os.Stdout, nil))
	child := handler.With(
		"Time", time.Now(),
	)
	slog.SetDefault(child)

	memguard.CatchInterrupt()
	defer memguard.Purge()

	db, err := Connect_to_BD.Connect()
	if err != nil {
		slog.Error("Err_from_register 1 ", err)
		return
	}
	cfg, err := Helpers.S3Helper()
	if err != nil {
		return
	}

	redisConn := RedisInteration.ConnectToRedis()
	defer redisConn.Close()

	TokensRealization := TokenInteraction.ControlTokens{A: nil}
	DatabaseRealization := UserInteraction.DB{Db: db}
	manageTokensImplRealization := manageTokensImpl.ManageTokensImpl{}
	s3Connect := s3Interation.ConntrolerForS3{}
	RedisStruct := RedisInteration.RedisInterationLayer{
		Re: redisConn,
	}
	InfoMange := FileKeyInteration.FileInfoController{}
	encryptKey := encryptionKeyInteration.EncryptionKey{}

	HandlerPack := &Handlers.HandlerPack{
		Tokens:    &TokensRealization,
		Database:  &DatabaseRealization,
		TokenImpl: manageTokensImplRealization,
		S3Conn:    &s3Connect,
		S3Connect: cfg,
		RedisConn: &RedisStruct,
		FileInfo:  &InfoMange,
		Choose:    &encryptKey,
	}
	Sa := Handlers.CollectorPack(*HandlerPack)

	router := mux.NewRouter()

	//router = router.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
	//	if request.Host != "filesbes.com" {
	//		return false
	//	}
	//	slog.Info(request.Host)
	//
	//	return true
	//}).Subrouter()

	slog.Info("Starting Server", Sa.S.FileInfo.SayHi())
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

	Sa.SwapKeyFirst()

	go func() {
		for t := range ticker.C {
			slog.Info("Got a ticker", t)
			Sa.SwapKeys()
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

	router.HandleFunc("/login/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.Login(writer, request, Sa)

	}).Methods("POST")
	router.HandleFunc("/register/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.Register(writer, request, Sa)

	}).Methods("POST")

	router.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithEncrypt(writer, request, Sa)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)
	router.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithNotEncrypt(writer, request, Sa)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)

	router.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderNoEncrypt(writer, request, router, Sa)

	}).Methods(http.MethodPost)
	router.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderEncrypt(writer, request, router, Sa)

	}).Methods(http.MethodPost)
	router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.GetFrom(writer, request, Sa)

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

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Server couldn't start", err)
		return

	}

}
