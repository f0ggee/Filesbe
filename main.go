package main

import (
	"Kaban/cmds"
	"Kaban/internal/Deliver"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/FileControls"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/s3Repo"
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Helpers"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("cannot load env file", "Error", err)
	}
	Key1 := &memguard.LockedBuffer{}
	Key2 := &memguard.LockedBuffer{}
	cmds.SettingSlog()
	memguard.CatchInterrupt()
	defer memguard.Purge()

	db, err := DatabaseControl.Connect()
	if err != nil {
		slog.Error("Error connect to database", "error", err)
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
	S3Collector := cmds.GetS3Collector(cfg, OldS3Connect)
	ParserCollector := cmds.GetRepoParsersCollector()
	FileCollector := cmds.GetFileControlCollector()
	EncrypterCollector := cmds.GetEncrypterKeysCollector(Key1, Key2)
	CryptoCollector := cmds.GetNewCryptoCollector(cmds.NewCryptoCollectorInput{
		OurPrivateKey:   []byte(os.Getenv("Our_Private_Key")),
		MasterPublicKet: []byte(os.Getenv("Publick_Key_Master_Server")),
		Decode:          ParserCollector.Decode,
	})
	AuthCollector := cmds.GetAuthTokensCollector([]byte(os.Getenv("KEY1")))
	DatabaseCollector := cmds.GetDatabaseManageCollector(db)
	DeliverPackagesCollector := cmds.GetDeliverPackagesCollector(RepoSessionHandle.GetCookieStore())
	GrcpCollector := cmds.GetGrpcCollector(CryptoCollector.Encrypt, CryptoCollector.Decrypt, ParserCollector.Decode, CryptoCollector.Validate, EncrypterCollector.Keys, CryptoCollector.Keys)
	FileControlsCollector := cmds.GetFileControlsCollector()
	RedisCollector := cmds.GetRedisCollector(redisConn)
	ControllerDonwloadBuilder := cmds.GetControllerDownloadBuilder(DeliverPackagesCollector)

	router, getRequest, postRequest, StaticFiles := Routers()
	cmds.GetAboutProjectUrl(getRequest)
	cmds.GetDefaultRouter(router)
	cmds.GetPhotoRequest(StaticFiles)
	cmds.GetLoginRequest(postRequest)
	cmds.SetRobotsRequest(router)
	cmds.GetInformationPage(getRequest)
	cmds.GetRegisterUrl(postRequest)
	cmds.GetMainUrl(getRequest)
	cmds.GetSitemap(router)
	cmds.GetProtectRequest(postRequest)
	cmds.GetUrlRequest(router)
	cmds.GetDownloadRequest(getRequest)
	cmds.GetEncryptDownload(getRequest)
	cmds.GetLoginApi(postRequest)
	cmds.GetRegisterApi(postRequest)
	cmds.GetDownloadApi(postRequest, router)
	cmds.GetEncryptDownloadApi(postRequest, router)
	cmds.GetMainApi(router)
	cmds.GetDoUrlApi(router)

	KeysController.FillOldKey()
	TimeSwaping := Sa.SwapKeyFirst()
	slog.Info("This time", "Time", TimeSwaping)
	ticker := time.NewTicker(TimeSwaping)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			slog.Time("Func Ticker: Got a ticker", t)
			Time := Sa.SwapKeys()
			slog.Duration("Func Ticker: Time", Time)
			ticker = time.NewTicker(Time)
		}
	}()

	serverConfig := cmds.ServerConfig(router)
	defer serverConfig.Close()

	slog.Info("The server started at ", "Configure", serverConfig.Addr)
	if err = serverConfig.ListenAndServe(); err != nil {
		slog.Error("Server couldn't start", "Error", err)
		return

	}
}

func Routers() (*mux.Router, *mux.Router, *mux.Router, *mux.Router) {
	router := cmds.GetRouter()
	cmds.SetLogging(router)
	getRequest := cmds.GetGetRouter(router)
	postRequest := cmds.GetPostRouter(router)
	StaticFiles := router.PathPrefix("/Fronted").Subrouter()
	return router, getRequest, postRequest, StaticFiles
}
