package Controller

import (
	"context"
	"log/slog"
	"os"
)

type LoggerCustom struct {
	slog.Handler
}

var logger = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource:   true,
	Level:       slog.LevelError,
	ReplaceAttr: nil,
})

func (l *LoggerCustom) Handle(ctx context.Context, Attr slog.Record) error {

	value, ok := ctx.Value(RequestId).(int)
	switch ok {

	case false:
		Attr.Add(slog.Any("REQUEST ID", "Nil"))
		return l.Handler.Handle(ctx, Attr)

	case true:
		Attr.Add("REQUEST ID", value)
		return l.Handler.Handle(ctx, Attr)

	}

	return nil
}

var ErrorLogger = &LoggerCustom{logger}
var ControllerErrorLogger = slog.New(ErrorLogger)
