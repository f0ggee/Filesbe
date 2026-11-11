package main

import (
	"Kaban/cmds"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/Protocol/receive"
	"context"
	"errors"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/awnumar/memguard"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cmds.SettingSlog()
	memguard.CatchInterrupt()
	defer memguard.Purge()
	connections := cmds.GetConnections()
	defer connections.CloseConnections()
	RepoEncrypterKeys.SetKeys()
	grpcCollector := cmds.GetGrpcCollector(connections.GrpcConnect)
	router := mux.NewRouter()
	cryptoCollector := cmds.GetCryptoCollector()
	databaseCollector := cmds.GetDatabaseCollector(connections.Db)
	fileTransferringCollector := cmds.GetFileTransferringCollector()
	parserCollector := cmds.GetNewParsersCollector()
	redisCollector := cmds.GetNewRedisCollector(connections.Conn)
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

	controllers := cmds.CollectControllers(cmds.CollectorsController{
		Apps:                apps,
		AuthTokensCollector: authTokensCollector,
		SessionCollector:    sessionCollector,
		Router:              router,
		ParserCollector:     parserCollector,
	})
	protocolFirst := cmds.GetEncryptedPacketBuilder(cmds.IncomeEncryptedPacketBuilder{
		CryptoGraph: cryptoCollector,
		Parser:      parserCollector.Parser,
	})

	preparingData, err := protocolFirst.GetEncryptedPacket()
	if err != nil {
		return
	}
	grpcMaker, err := grpcCollector.Reqs.Make(context.Background())
	if err != nil {
		return
	}
	data, err := grpcMaker.SetNewKeyRequest(preparingData)
	if err != nil {
		return
	}
	GetPacketData := cmds.GetCheckPacketBuilder(cmds.IncomeCheckPacketBuilder{
		Parser:      parserCollector.Parser,
		CryptoGraph: cryptoCollector,
	})

	timeDuration, err := GetPacketData.GetPacketData(data)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(timeDuration)
	defer ticker.Stop()

	stopSignal := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Go(func() {
		for _ = range ticker.C {
			exchangeTime, err := tryPacket(redisCollector, GetPacketData)
			if err != nil {
				stopSignal <- struct{}{}
				return
			}
			ticker = time.NewTicker(exchangeTime)
		}
	})

	var serverConfig = cmds.ServerConfig(router)
	runtime.GC()
	cmds.SetRouters(router, &controllers)

	wg.Go(func() {
		select {
		case _, _ = <-stopSignal:
			return

		default:

		}
		slog.Info("The server started at", "Configure", serverConfig.Addr)
		if err = serverConfig.ListenAndServe(); err != nil {
			slog.Error("Server couldn't start", "Error", err)
			return

		}
	})

	wg.Wait()
}

var (
	ErrorGetPacket     = errors.New("error to get a new packet")
	ErrorDecryptPacket = errors.New("error to decrypt a packet")
)

func tryPacket(redisCollector cmds.RedisCollector, GetPacketData receive.GetNewKey) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	slog.Info("start the plan exchange", slog.Time("Time", time.Now()))
	packet, err := redisCollector.Read.GetKey(ctx)
	if err != nil {
		return 0, ErrorGetPacket
	}
	exchangeTime, err := GetPacketData.GetPacketData(packet)
	if err != nil {
		return 0, ErrorDecryptPacket
	}
	return exchangeTime, nil
}
