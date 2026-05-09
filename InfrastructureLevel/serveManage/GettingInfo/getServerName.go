package GettingInfo

import (
	"log/slog"
)

func (s *SeverManage) GetServerName(i int) string {

	switch i {

	case 1:
		return "SERVER_1"

	case 2:
		return "SERVER_2"

	}

	slog.Info("Couldn't find the server name", "ServerNumber", i)
	return ""
}
