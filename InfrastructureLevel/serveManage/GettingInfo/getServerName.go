package GettingInfo

import (
	"log/slog"
)

func (s *SeverManage) GetServerName(i int) string {

	slog.Info("Start searching info about a server")

	switch i {

	case 1:
		slog.Info("Found the server name", "ServerNumber", i)
		return "SERVER_1"

	case 2:
		slog.Info("Found the server name", "ServerNumber", i)
		return "SERVER_2"

	}

	slog.Info("Couldn't find the server name", "ServerNumber", i)
	return ""
}
