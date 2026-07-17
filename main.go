package main

import (
	"Kaban/cmds"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Helpers"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
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
	router, getRequest, postRequest, StaticFiles := Routers()
	serverConfig := cmds.ServerConfig(router)
	defer serverConfig.Close()

	///Collectors
	S3Collector := cmds.GetS3Collector(cfg, OldS3Connect)
	ParserCollector := cmds.GetRepoParsersCollector()
	FileCollector := cmds.GetFileControlCollector()
	EncrypterKeysCollector := cmds.GetEncrypterKeysCollector(Key1, Key2)
	CryptoCollector := cmds.GetNewCryptoCollector(cmds.NewCryptoCollectorInput{Decode: ParserCollector.Decode})
	PrivateKey := []byte(os.Getenv("Our_Private_Key"))
	MasterKey := []byte(os.Getenv("Public_Key_Master_Server"))
	ServerKeysCollector := cmds.GetNewServerKeysCollector(PrivateKey, MasterKey)
	AuthCollector := cmds.GetAuthTokensCollector([]byte(os.Getenv("KEY1")))
	DatabaseCollector := cmds.GetDatabaseManageCollector(db)
	DeliverPackagesCollector := cmds.GetDeliverPackagesCollector(RepoSessionHandle.GetCookieStore())
	SessionCollector := cmds.GetSessionCollector(RepoSessionHandle.GetCookieStore())
	GrpcCollector := cmds.GetGrpcCollector(CryptoCollector.Encrypt, CryptoCollector.Decrypt, ParserCollector.Decode, CryptoCollector.Validate, *EncrypterKeysCollector, *ServerKeysCollector)
	RedisCollector := cmds.GetRedisCollector(redisConn)

	//Application builders
	DownloadApplicationBuilder := cmds.GetDownloadApplicationBuilder(S3Collector, RedisCollector, FileCollector)
	DownloadEncryptApplicationBuilder := cmds.GetDownloadEncryptApplicationBuilder(RedisCollector, S3Collector, EncrypterKeysCollector, CryptoCollector, FileCollector)
	LoginApplicationBuilder := cmds.GetLoginApplicationBuilder(CryptoCollector, &DatabaseCollector, AuthCollector)
	RegisterApplicationBuilder := cmds.GetRegisterApplicationBuilder(&DatabaseCollector, AuthCollector, CryptoCollector)
	UploaderApplicationBuilder := cmds.GetUploadApplicationBuilder(CryptoCollector, FileCollector, ParserCollector, S3Collector, RedisCollector)
	UploaderEncryptApplicationBuilder := cmds.GetUploadEncryptApplicationBuilder(cmds.UploadEncryptBuilderIncomeData{
		F:          FileCollector,
		S3:         S3Collector,
		Z:          ParserCollector,
		C:          CryptoCollector,
		ServerKeys: ServerKeysCollector,
		R:          RedisCollector,
	})
	//Controllers builders
	ControllerDownloadBuilder := cmds.GetControllerDownloadBuilder(&cmds.RegisterControllerBuilderIncomeData{
		Answ:    DeliverPackagesCollector.DownloadCollector.Answ,
		UrlWork: DeliverPackagesCollector.DownloadCollector.Answ,
		App:     DownloadApplicationBuilder,
	})
	ControllerDownloadEncryptBuilder := cmds.GetEncryptDownloadControllerBuilder(cmds.EncryptDownloadControllerIncomeData{
		Answ:     DeliverPackagesCollector.DownloadEncryptCollector.Answ,
		UrlBuild: DeliverPackagesCollector.DownloadEncryptCollector.Answ,
		App:      DownloadEncryptApplicationBuilder,
	})
	ControllerUploadEncryptBuilder := cmds.GetUploaderEncrypterControllerBuilder(&cmds.UploaderBuilderIncomeData{
		R:           router,
		ReadSession: &SessionCollector.Session,
		AuthCheck:   AuthCollector.Validate,
		Answers:     DeliverPackagesCollector.UploaderEncrypterCollector.Answ,
		Build:       DeliverPackagesCollector.UploaderEncrypterCollector.Answ,
		A:           *UploaderEncryptApplicationBuilder,
	})
	ControllerUploadBuilder := cmds.GetUploaderControllerBuilder(cmds.UploaderEncryptBuilderIncomeData{
		S:       DeliverPackagesCollector.UploaderCollector.Answ,
		Builder: DeliverPackagesCollector.UploaderCollector.Answ,
		Session: &SessionCollector.Session,
		Auth:    AuthCollector.Validate,
		App:     *UploaderApplicationBuilder,
	})
	ControllerLoginBuilder := cmds.GetControllerLoginBuilder(cmds.NewLoginIncomeData{
		S:      &DeliverPackagesCollector.LoginCollector.Answ,
		Sess:   &SessionCollector.Session,
		Parses: ParserCollector.Decode,
		App:    LoginApplicationBuilder,
	})
	ControllerRegisterBuilder := cmds.GetControllerRegisterBuilder(cmds.RegisterBuilderIncomeData{
		Answ:    DeliverPackagesCollector.RegisterCollector.Answ,
		Session: &SessionCollector.Session,
		D:       ParserCollector.Decode,
		App:     *RegisterApplicationBuilder,
	})
	ControllerUrlUploader := cmds.GetControllerUrlUploaderBuilder(DeliverPackagesCollector.UrlBuilderCollector.Answ)
	ControllerCheckAuth := cmds.GetControllerCheckAuthBuilder(cmds.CheckUserBuilderIncomeData{
		Session: &SessionCollector.Session,
		Auth:    AuthCollector.Validate,
	})
	ControllerProtocolManageBuilder := cmds.GetProtocolManageBuilder(cmds.ProtocolManageBuilder{
		EncrypterKeys: *EncrypterKeysCollector,
		ServerKeys:    *ServerKeysCollector,
		C:             CryptoCollector,
		Parser:        ParserCollector,
		GrpcConn:      *GrpcCollector,
	})

	cmds.GetAboutProjectUrlRouter(getRequest)
	cmds.GetDefaultRouter(router)
	cmds.GetPhotoRequest(StaticFiles)
	cmds.GetLoginPageRouter(postRequest)
	cmds.SetRobotsRouter(router)
	cmds.GetInformationPageRouter(getRequest)
	cmds.GetRegisterPageRouter(postRequest)
	cmds.GetMainPageRouter(getRequest)
	cmds.GetSitemapRouter(router)
	cmds.GetProtectPageRouter(postRequest)
	cmds.GetUrlPageRouter(router)
	cmds.GetDownloadApi(getRequest, ControllerDownloadBuilder)
	cmds.GetEncryptDownloadApi(getRequest, ControllerDownloadEncryptBuilder)
	cmds.GetLoginApi(postRequest, ControllerLoginBuilder)
	cmds.GetRegisterApiRouter(postRequest, ControllerRegisterBuilder)
	cmds.GetUploaderApiRouter(postRequest, ControllerUploadBuilder)
	cmds.GetEncryptUploaderApiRouter(postRequest, ControllerUploadEncryptBuilder)
	cmds.GetMainApiRouter(getRequest, ControllerCheckAuth)
	cmds.GetDoUrlApiRouter(getRequest, ControllerUrlUploader)

	TimeSwaping := ControllerProtocolManageBuilder.Start.GetExchangerInitializer()
	slog.Info("This time", "Time", TimeSwaping)
	ticker := time.NewTicker(TimeSwaping)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			slog.Time("Func Ticker: Got a ticker", t)
			Time := ControllerProtocolManageBuilder.Processing.GetPlanningExchanger()
			slog.Duration("Func Ticker: Time", Time)
			ticker = time.NewTicker(Time)
		}
	}()

	slog.Info("The server started at", "Configure", serverConfig.Addr)
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
