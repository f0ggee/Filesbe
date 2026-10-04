package main

import (
	"Kaban/cmds"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DatabaseControl"
	gr "Kaban/internal/InfrastructureLayer/Grpc"
	"Kaban/internal/InfrastructureLayer/RedisInteration"
	"Kaban/internal/Protocol/receive"
	"Kaban/internal/Protocol/send"
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/awnumar/memguard"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	defer RedisInteration.RedisConn.Close()
	defer gr.GrpcConn.Close()
	cmds.SettingSlog()
	memguard.CatchInterrupt()
	defer memguard.Purge()
	db, err := DatabaseControl.Connect()
	if err != nil {
		slog.Error("Error connect to database", "error", err)
		return
	}
	defer db.Close()
	var databaseConn = DatabaseControl.DatabaseConn{
		Db: db,
	}
	grpcCollector := cmds.GetGrpcCollector()

	router := cmds.GetRouter()
	cryptoCollector := cmds.GetCryptoCollector()
	databaseCollector := cmds.GetDatabaseCollector(databaseConn)
	fileTransferringCollector := cmds.GetFileTransferringCollector()
	parserCollector := cmds.GetNewParsersCollector()
	redisCollector := cmds.GetNewRedisCollector()
	sessionKeys := cmds.GetNewSessionKeysCollector()
	authTokensCollector := cmds.GetNewTokensAuthCollector()
	sessionCollector := cmds.GetNewSessionCollector()

	downloadApplication := cmds.GetDownloadApplicationBuilder(cmds.DownloadBuilderApplication{
		Trans: fileTransferringCollector,
		Red:   redisCollector,
		Parse: parserCollector,
	})
	downloadEncryptApplication := cmds.GetDownloadEncryptApplicationBuilder(cmds.DownloadEncryptApplication{
		Red:          redisCollector,
		Crypt:        cryptoCollector,
		Keys:         sessionKeys,
		Transferring: fileTransferringCollector,
	})
	loginApplication := cmds.GetLoginApplicationBuilder(cmds.LoginApplication{
		DB:     databaseCollector,
		Crypt:  cryptoCollector,
		Tokens: authTokensCollector,
	})
	registerApplication := cmds.GetRegisterBuilder(cmds.RegisterApplication{
		DB:     databaseCollector,
		Tokens: authTokensCollector,
		Crypt:  cryptoCollector,
	})
	uploadApplication := cmds.GetUploadBuilder(cmds.UploadApplication{
		Crypt: cryptoCollector,
		Parse: parserCollector,
		S3:    fileTransferringCollector.S3FileTransferringCollector,
		Red:   redisCollector,
	})
	uploadEncrypt := cmds.GetUploadEncryptBuilder(cmds.UploadEncrypt{
		Parse:        parserCollector,
		Crypt:        cryptoCollector,
		Transferring: fileTransferringCollector,
		Red:          redisCollector,
	})

	controllerDownload := cmds.GetControllerDownloadBuilder(downloadApplication)
	controllerEncryptDownload := cmds.GetControllerDownloadEncryptBuilder(&downloadEncryptApplication)
	controllerUploader := cmds.GetControllerFileUploadBuilder(cmds.ControllerFileUploadBuilder{
		Token:    &authTokensCollector.Rf,
		Sessions: sessionCollector,
		App:      &uploadApplication,
	})

	controllerUploadEncrypt := cmds.GetControllerFileUploadEncryptBuilder(cmds.ControllerFileUploadEncrypt{
		Token:    &authTokensCollector.Jwt,
		Sessions: sessionCollector,
		App:      uploadEncrypt,
	}, router)
	controllerLogin := cmds.GetControllerLoginBuilder(cmds.ControllerLoginBuilder{
		Sessions: sessionCollector,
		Parser:   parserCollector,
		App:      &loginApplication,
	})
	controllerRegister := cmds.GetControllerRegisterBuilder(cmds.RegisterController{
		Sessions: sessionCollector,
		App:      registerApplication,
		Parser:   parserCollector,
	})
	controllerUrlBuild := cmds.GetControllerUrlBuildBuilder()

	var protocolFirst = send.FirstExchange{
		CryptoGenerate: cryptoCollector.Generate,
		Encoder:        parserCollector.Parser,
		Aes:            cryptoCollector.AesGcmRealization,
		Rsa:            &cryptoCollector.RsaRealization,
	}

	preparingData, err := protocolFirst.GetEncryptedPacket()
	if err != nil {
		panic(err)
	}
	grpcMaker, err := grpcCollector.Reqs.Make(context.Background())
	if err != nil {
		panic(err)
	}
	data, err := grpcMaker.SetNewKeyRequest(preparingData)
	if err != nil {
		panic(err)
	}

	var GetPacketData = receive.GetNewKey{
		Decoder:        parserCollector.Parser,
		CryptoValidate: &cryptoCollector.Validate,
		TempKey:        sessionKeys.Keys,
	}

	timeDuration, err := GetPacketData.GetPacketData(data)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(timeDuration)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			slog.Time("Func Ticker: Got a ticker", t)
			Time := ControllerProtocolManageBuilder.Processing.GetPlanningExchanger()
			slog.Duration("Func Ticker: Time", Time)
			ticker = time.NewTicker(Time)
		}
	}()

	runtime.GC()
	slog.Info("The server started at", "Configure", serverConfig.Addr)
	if err = serverConfig.ListenAndServe(); err != nil {
		slog.Error("Server couldn't start", "Error", err)
		return

	}
}

func SetExchanger() {

}
