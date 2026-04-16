package cmds

import (
	"log/slog"
	"os"
	"time"
)

func SettingSlog() {
	handler := slog.New(slog.NewTextHandler(os.Stdout, nil))
	child := handler.With(
		"Time", time.Now().Format("2006-01-02 15:04:05"),
	)

	slog.SetDefault(child)
}
