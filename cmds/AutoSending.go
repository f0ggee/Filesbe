package cmds

import (
	"MasterServer_/Dto"
	InftarctionLevel "MasterServer_/InfrastructureLevel"
	"MasterServer_/InfrastructureLevel/GlobalProces"
	"MasterServer_/InfrastructureLevel/serveManage/GettingInfo"
	"log/slog"
)

func StartHandling(serverManagementPack *GettingInfo.SeverManage, Sa *GlobalProces.ControllingExchange) bool {
	for i := 1; i <= InftarctionLevel.ServersCount; i++ {
		ServerKey := serverManagementPack.GetServerKey(i)
		if ServerKey == nil {
			slog.Error("ServerKey is nil")
			continue
		}

		ServerName := serverManagementPack.GetServerName(i)
		if ServerName == "" {
			slog.Error("we can't find the server", "ServerNumber", i)
			continue
		}

		err := Sa.SwapKeys(ServerKey, Dto.Keys.NewPrivateKey.Bytes(), ServerName)
		if err != nil {
			continue
		}

	}
	return false
}
