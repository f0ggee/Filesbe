package main

import (
	"MasterServer_/DipendsInjective"
	"MasterServer_/Dto"
	InftarctionLevel "MasterServer_/InfrastructureLevel"
	"MasterServer_/InfrastructureLevel/CryptoImpl"
	"MasterServer_/InfrastructureLevel/CryptoImpl/CryprtoGenerator"
	CryptoKey2 "MasterServer_/InfrastructureLevel/CryptoImpl/CryptoKey"
	"MasterServer_/InfrastructureLevel/CryptoImpl/Decryptor"
	"MasterServer_/InfrastructureLevel/CryptoImpl/Encrypter"
	"MasterServer_/InfrastructureLevel/GlobalProces"
	"MasterServer_/InfrastructureLevel/Grpc/GrpcHandleData"
	"MasterServer_/InfrastructureLevel/Grpc/PacketValidation"
	"MasterServer_/InfrastructureLevel/MemguardManipulation"
	"MasterServer_/InfrastructureLevel/RedisUse"
	"MasterServer_/InfrastructureLevel/rsaKeyManipulation"
	"MasterServer_/InfrastructureLevel/serveManage/ConverterData"
	"MasterServer_/InfrastructureLevel/serveManage/GettingInfo"
	Cmds "MasterServer_/cmds"
	GrcpCmds "MasterServer_/cmds/GrpcConn"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("cannot load env file", "Error", err.Error())

	}
}

func main() {

	fmt.Println("s1", os.Getenv("SERVER_1"))

	Dto.Keys.NewPrivateKey, _ = memguard.NewBufferFromReader(rand.Reader, 2048)
	Dto.Keys.OldPrivateKey, _ = memguard.NewBufferFromReader(rand.Reader, 2048)
	Dto.Keys.MasterServerKey = os.Getenv("OUR_KEY")

	fmt.Println("sas", os.Getenv("SERVER_2"))
	handler := slog.New(slog.NewTextHandler(os.Stdout, nil))
	child := handler.With(
		"Time", time.Now(),
		"ServersCount", InftarctionLevel.ServersCount,
	)
	slog.SetDefault(child)

	memguard.CatchInterrupt()
	defer memguard.Purge()
	ConnectRedis := RedisUse.RedisConnect()
	defer ConnectRedis.Close()

	CryptoGenerate := CryprtoGenerator.CryprtoGenerating{}
	Decrypted := Decryptor.Decrypting{}
	Encrypting := Encrypter.Encryption{}

	CryptoKey := CryptoKey2.CryptoManging{M: CryptoImpl.CryptoData{MasterServerSecret: os.Getenv("OUR_KEY")}}

	GrpcHandlingData := GrpcHandleData.GrpcDataManagement{}
	PacketValidating := PacketValidation.ValidatePacketData{}
	MemguardControl := MemguardManipulation.MemgurdControl{}
	ServerInfo := GettingInfo.SeverManage{}
	redisConn := RedisUse.RedisUsing{Connect: ConnectRedis}
	RsaKeyControl := rsaKeyManipulation.RsaKeyManipulation{}
	ConvertData := ConverterData.ConvertingData{}

	Injective1 := DipendsInjective.NewRsaKeyManipulationWithRsaAndMemory(&MemguardControl, &RsaKeyControl)
	AnotherProcessController := GlobalProces.ControllingExchange{
		E: GlobalProces.ProcessController{
			Cryptos:          &Encrypting,
			CryptoGen:        &CryptoGenerate,
			RedisInteracting: &redisConn,
			ServerManagement: &ServerInfo,
			CryptoKey:        &CryptoKey,
		},
	}
	SwapRsaKey(*Injective1)
	ticker := time.NewTicker(InftarctionLevel.TimeForSwapping)
	defer ticker.Stop()

	G := GrcpCmds.GetGrpcConn("tcp", os.Getenv("GRPC_ADDRESS"))

	go G.GrpcHandleRequests(GrpcHandlingData, ServerInfo, Encrypting, PacketValidating, CryptoGenerate, ConvertData, Decrypted)

	for _ = range ticker.C {
		slog.Info("We got the tick")
		SwapRsaKey(*Injective1)
		handling := Cmds.StartHandling(&ServerInfo, &AnotherProcessController)
		if handling {
			return
		}
		slog.Info("Finished the exchange")
	}

}

func SwapRsaKey(RsaKey DipendsInjective.RsaKeyManipulationWithRsaAndMemory) {

	slog.Info("Swaping RSA key in memory START")
	fmt.Println("Our key", os.Getenv("GRPC_ADDRESS"))

	TemporallySaving := memguard.NewBufferFromBytes(RsaKey.Key.GenerateRsaKey())
	defer TemporallySaving.Destroy()

	Dto.Keys.Mu.Lock()
	RsaKey.KeyAndMemory.SwapingOldKey()
	RsaKey.KeyAndMemory.InstallingNewKey(TemporallySaving.Bytes())
	Dto.Keys.Mu.Unlock()
	slog.Info("Swaping RSA key in memory END")
}
