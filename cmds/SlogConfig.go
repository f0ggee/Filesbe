package cmds

import (
	"log/slog"
	"os"
	"time"
)

func SettingSlog() {
	handler := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			if a.Key == "Time" || a.Key == "time" {
				a.Value = slog.StringValue(time.Now().Format("2006-01-02 15:04"))
			}
			return a
		},
	}))

	slog.SetDefault(handler)
}
