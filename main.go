package main

import (
	"Kaban/cmds"
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

	apps := cmds.CollectApplicationBuilders(cmds.CollectorsApp{
		FileTransferringCollector: fileTransferringCollector,
		RedisCollector:            redisCollector,
		ParserCollector:           parserCollector,
		CryptoCollector:           cryptoCollector,
		SessionKeys:               sessionKeys,
		DatabaseCollector:         databaseCollector,
		AuthTokensCollector:       authTokensCollector,
	})

	controlers := cmds.CollectControllers(cmds.CollectorsController{
		Apps:                apps,
		AuthTokensCollector: authTokensCollector,
		SessionCollector:    sessionCollector,
		Router:              router,
		ParserCollector:     parserCollector,
	})

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
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	go func() {
		for _ = range ticker.C {
			exchangeTime := planExchange(redisCollector, ctx, GetPacketData)
			ticker = time.NewTicker(exchangeTime)
		}
	}()
	var serverConfig = cmds.ServerConfig(router)
	runtime.GC()
	cmds.SetRouters(router, &controlers)
	slog.Info("The server started at", "Configure", serverConfig.Addr)
	if err = serverConfig.ListenAndServe(); err != nil {
		slog.Error("Server couldn't start", "Error", err)
		return

	}
}

func planExchange(redisCollector cmds.RedisCollector, ctx context.Context, GetPacketData receive.GetNewKey) time.Duration {
	slog.Info("Main: starting planning exchange", "TIME", time.Now().Hour())
	packet, err := redisCollector.Read.GetKey(ctx)
	if err != nil {
		return 0
	}
	exchangeTime, err := GetPacketData.GetPacketData(packet)
	if err != nil {
		return 0
	}
	return exchangeTime
}
