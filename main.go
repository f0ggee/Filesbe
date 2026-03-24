package main

import (
	Controller2 "Kaban/internal/Controller"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"

	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ControllingTokens"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Generating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
	CryptoChecking "Kaban/internal/InfrastructureLayer/Crypto/Checking"
	"Kaban/internal/InfrastructureLayer/Crypto/Decription"
	"Kaban/internal/InfrastructureLayer/Crypto/Encryption"
	CryptoGenerater "Kaban/internal/InfrastructureLayer/Crypto/Generating"
	"Kaban/internal/InfrastructureLayer/DataConverting"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Checking"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Reading"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Writinig"
	"Kaban/internal/InfrastructureLayer/FileKeyInteration/HandleFileInfo"
	"Kaban/internal/InfrastructureLayer/FileKeyInteration/HandlerFile"
	"Kaban/internal/InfrastructureLayer/GrpcManage/HandlingRequests"
	"Kaban/internal/InfrastructureLayer/GrpcManage/PacketChecking"
	"Kaban/internal/InfrastructureLayer/GrpcManage/SendingRequest"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RedisInteration/DeletingRedis"
	"Kaban/internal/InfrastructureLayer/RedisInteration/ReadingRedis"
	"Kaban/internal/InfrastructureLayer/RedisInteration/WritingRedis"
	"Kaban/internal/InfrastructureLayer/s3Interation/DeleterS3"
	"Kaban/internal/Service/Handlers"
	"Kaban/internal/Service/Helpers"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"github.com/gorilla/mux"
)

func main() {
	SettingSlog()

	memguard.CatchInterrupt()
	defer memguard.Purge()

	db, err := DatabaseControl.Connect()
	if err != nil {
		slog.Error("Err_from_register 1 ", err)
		return
	}
	defer db.Close()
	cfg, err := Helpers.S3Helper()
	if err != nil {
		return
	}
	redisConn := RedisInteration.ConnectToRedis()
	defer redisConn.Close()

	ManagingAuthTokens := ControllingTokens.ManageTokens{}
	GeneratingAuthTokens := Generating.CreatingTokens{}
	CheckingAuthToken := ValidatingTokens.Checking{}
	CryptoEncryption := Encryption.Encrypter{}
	CryptoDecryption := Decription.DecryptionData{}
	CryptoGenerate := CryptoGenerater.Generating{}
	CryptoCheck := CryptoChecking.Checking{}
	DbCheck := Checking.CheckerDb{Db: db}
	DbReading := Reading.Read{Db: db}
	DbWriting := Writinig.Writer{Db: db}
	ConverterJson := DataConverting.ConvertingData{}
	ProcessedFile := HandlerFile.ProcessingFile{}
	ProcessedFileInfo := HandleFileInfo.ProcessingFileInfo{}
	PacketValidate := PacketChecking.PacketValidating{}

	GrpcHandlingRequests := HandlingRequests.HandlerGrpcRequest{
		CryptoEncrypt:    &CryptoEncryption,
		CryptoDecrypt:    CryptoDecryption,
		CryptoValidate:   &CryptoCheck,
		ValidationPacket: PacketValidate,
	}

	SendingGrcp := SendingRequest.SenderRequests{}

	DeleterRds := DeletingRedis.DeleterRedis{Re: redisConn}
	ReaderRedis := ReadingRedis.RedisReader{Re: redisConn}
	WriterRedis := WritingRedis.Writing{Re: redisConn}
	//CheckerRedis := RedisChecking.ValidationRedis{Re: redisConn}

	S3Deleter := DeleterS3.DeleterS3{Conf: cfg}

	HandlerPack := Handlers.HandlerPackCollect{
		S3: Handlers.S3Controlling{
			Deleter:   &S3Deleter,
			S3Connect: cfg,
		},
		Crypto: Handlers.HandlerPackCrypto{
			Validate: &CryptoCheck,
			Decrypt:  &CryptoDecryption,
			Encrypt:  &CryptoEncryption,
			Generate: &CryptoGenerate,
		},
		FileInfo: Handlers.HandlerFileManagerPack{
			FileInfo:     ProcessedFileInfo,
			FileManaging: ProcessedFile,
		},
		AuthTokens: Handlers.HandlerPackAuthTokens{
			Manage:          ManagingAuthTokens,
			GeneratingToken: GeneratingAuthTokens,
			Checking:        CheckingAuthToken,
		},
		DatabaseControlling: Handlers.DatabaseControlling{
			Writer:  &DbWriting,
			Reader:  &DbReading,
			Checker: &DbCheck,
		},
		RedisControlling: Handlers.RedisControlling{
			Deleter: &DeleterRds,
			Reader:  &ReaderRedis,
			Writer:  &WriterRedis,
		},
		Grpc: Handlers.HandlerGrpc{
			GrpcSendingRequest: &SendingGrcp,
			ProcessingRequests: GrpcHandlingRequests,
		},
		Convert: Handlers.Converter{Converting: ConverterJson},
	}
	Sa := Handlers.NewHandlerPackCollect(HandlerPack.S3, HandlerPack.Crypto, HandlerPack.FileInfo, HandlerPack.AuthTokens, HandlerPack.DatabaseControlling, HandlerPack.RedisControlling, HandlerPack.Grpc, HandlerPack.Convert)

	slog.Info("Starting server", Sa.Crypto.Decrypt.SayHello("EE"))
	router := mux.NewRouter()
	newRouter := router.PathPrefix("/").Subrouter()
	newRouter.Use(Controller2.CheckBots)
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
	TimeSwaping := Sa.SwapKeyFirst()

	ticker := time.NewTicker(time.Until(time.Now().Add(TimeSwaping)))
	defer ticker.Stop()

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

	newRouter.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithEncrypt(writer, request, Sa)

		//Handlers.Delete(ch)

	}).Methods(http.MethodGet)
	newRouter.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {

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

func SettingSlog() {
	handler := slog.New(slog.NewTextHandler(os.Stdout, nil))
	child := handler.With(
		"Time", time.Now(),
	)

	slog.SetDefault(child)
}
