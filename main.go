package main

import (
	"Kaban/cmds"
	Controller2 "Kaban/internal/Controller"
	"Kaban/internal/Controller/Middlewares"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/KeysManager"
	"Kaban/internal/InfrastructureLayer/RedisInteration/RedisChecking"
	"Kaban/internal/InfrastructureLayer/s3Interation"
	"Kaban/internal/InfrastructureLayer/s3Interation/S3Downloader"
	"Kaban/internal/InfrastructureLayer/s3Interation/S3Uploader"
	"Kaban/internal/Service/Helpers"
	"fmt"
	"sync"

	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ControllingTokens"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
	"Kaban/internal/InfrastructureLayer/Crypto/Checking"
	"Kaban/internal/InfrastructureLayer/Crypto/Decription"
	"Kaban/internal/InfrastructureLayer/Crypto/Encryption"
	"Kaban/internal/InfrastructureLayer/Crypto/Generating"
	"Kaban/internal/InfrastructureLayer/DataConverting"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Reading"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Validator"
	"Kaban/internal/InfrastructureLayer/DatabaseControl/Writinig"
	"Kaban/internal/InfrastructureLayer/FileKeyInteration/HandleFileInfo"
	"Kaban/internal/InfrastructureLayer/FileKeyInteration/HandlerFile"
	"Kaban/internal/InfrastructureLayer/GrpcManage/HandlingRequests"
	"Kaban/internal/InfrastructureLayer/GrpcManage/PacketChecking"
	"Kaban/internal/InfrastructureLayer/GrpcManage/Sender"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RedisInteration/DeletingRedis"
	"Kaban/internal/InfrastructureLayer/RedisInteration/ReadingRedis"
	"Kaban/internal/InfrastructureLayer/RedisInteration/WritingRedis"
	"Kaban/internal/InfrastructureLayer/s3Interation/DeleterS3"
	"Kaban/internal/Service/Handlers"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"github.com/gorilla/mux"
)

func main() {
	cmds.SettingSlog()

	memguard.CatchInterrupt()
	defer memguard.Purge()

	db, err := DatabaseControl.Connect()
	if err != nil {
		slog.Error("Error connect to database", err)
		return
	}
	defer db.Close()
	cfg, err := Helpers.S3Helper()
	if err != nil {
		return
	}
	redisConn := RedisInteration.ConnectToRedis()
	defer redisConn.Close()

	OldS3Connect, err := Helpers.Inzelire()
	if err != nil {
		slog.Error("Error connect to s3 old ", "Error", err)
		return
	}

	ManagingAuthTokens := ControllingTokens.ManageTokens{}
	GeneratingAuthTokens := Creating.CreatingTokens{}
	CheckingAuthToken := ValidatingTokens.Checking{}
	CryptoEncryption := Encryption.Encrypter{}
	CryptoDecryption := Decription.DecryptionData{}
	CryptoGenerate := Generating.Generating{}
	CryptoCheck := Checking.Validating{}
	DbCheck := Validator.CheckerDb{Db: db}
	DbReading := Reading.Read{Db: db}
	DbWriting := Writinig.Writer{Db: db}
	ConverterJson := DataConverting.ConvertingData{}
	ProcessedFile := HandlerFile.ProcessingFile{}
	ProcessedFileInfo := HandleFileInfo.ProcessingFileInfo{}
	PacketValidate := PacketChecking.PacketValidating{}
	KeysController := &KeysManager.Updater{
		Mu:            &sync.RWMutex{},
		NewPrivateKey: &memguard.LockedBuffer{},
		OldPrivateKey: &memguard.LockedBuffer{},
		OurPrivateKey: os.Getenv("Our_Private_Key"),
		MasterKey:     os.Getenv("Publick_Key_Master_Server"),
	}
	GrpcHandlingRequests := HandlingRequests.HandlerGrpcRequest{
		CryptoEncrypt:    &CryptoEncryption,
		CryptoDecrypt:    CryptoDecryption,
		CryptoValidate:   &CryptoCheck,
		ValidationPacket: PacketValidate,
		Keys:             KeysController,
	}

	SendingGrcp := Sender.SenderRequests{}

	DeleterRds := DeletingRedis.DeleterRedis{Re: redisConn}
	ReaderRedis := ReadingRedis.RedisReader{Re: redisConn}
	WriterRedis := WritingRedis.Writing{Re: redisConn}

	CheckerRedis := RedisChecking.ValidationRedis{Re: redisConn}

	S3Information := s3Interation.Variables{
		Bucket:     os.Getenv("BUCKET"),
		S3Connect:  cfg,
		OldConnect: OldS3Connect,
	}

	S3Deleter := DeleterS3.DeleterS3{
		S3Info: S3Information,
	}
	S3Uploading := S3Uploader.Uploading{
		S3Info: S3Information}

	S3Download := S3Downloader.S3Download{S3Info: S3Information}

	HandlerPack := Handlers.HandlerPackCollect{
		S3: Handlers.S3Controlling{
			Deleter:    &S3Deleter,
			Uploader:   &S3Uploading,
			S3Download: S3Download,
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
			Deleter:      &DeleterRds,
			Reader:       &ReaderRedis,
			Writer:       &WriterRedis,
			CheckerRedis: &CheckerRedis,
		},
		Grpc: Handlers.HandlerGrpc{
			GrpcSendingRequest: &SendingGrcp,
			ProcessingRequests: GrpcHandlingRequests,
		},
		Convert: Handlers.Converter{Converting: ConverterJson},
		Keys:    Handlers.KeysControlling{ControllerKey: KeysController},
	}
	Sa := Handlers.NewHandlerPackCollect(HandlerPack.S3, HandlerPack.Crypto, HandlerPack.FileInfo, HandlerPack.AuthTokens, HandlerPack.DatabaseControlling, HandlerPack.RedisControlling, HandlerPack.Grpc, HandlerPack.Convert, HandlerPack.Keys)

	router := mux.NewRouter()

	router.Use(Middlewares.Logging)
	newRouter := router.PathPrefix("/").Subrouter()
	newRouter.Use(Middlewares.CheckBots)
	StaticFiles := router.PathPrefix("/Fronted").Subrouter()

	router.HandleFunc("/aboutProject", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/InfoPageAboutApp.html")

	})

	StaticFiles.Handle("/favicon.png", http.FileServer(http.Dir("internal/Service")))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "internal/Service/Fronted/Maine.html")
		}

	})

	KeysController.FillOldKey()
	TimeSwaping := Sa.SwapKeyFirst()

	fmt.Println(TimeSwaping)

	ticker := time.NewTicker(time.Until(time.Now().Add(TimeSwaping)))
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			slog.Info("Got a ticker", t)
			Sa.SwapKeys()
		}
	}()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "internal/Service/Fronted/Login.html")

	})

	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./robots.txt")

	})

	router.HandleFunc("/informationPage", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/InformationPage.html")

	}).Name("NameFile")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/Service/Fronted/Register.html")
	})
	router.HandleFunc("/main", func(writer http.ResponseWriter, request *http.Request) {

		http.ServeFile(writer, request, "internal/Service/Fronted/Main_Page.html")

	})
	router.HandleFunc("/sitemap.xml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/sitemap.xml")

	})

	router.HandleFunc("/protect", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/Protecion.html")

	})

	router.HandleFunc("/URL/{name}/{bool}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "internal/Service/Fronted/UrlFronted.html")

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

		Controller2.BuildUrl(writer, request)

	}).Methods(http.MethodGet)

	serverConfig := cmds.ServerConfig(router)
	defer serverConfig.Close()

	err = serverConfig.ListenAndServe()
	if err != nil {
		slog.Error("Server couldn't start", err)
		return

	}
}
