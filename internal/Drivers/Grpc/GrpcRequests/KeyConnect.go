package Requests

import (
	"Kaban/internal/DomainLevel"
	"context"
	"log/slog"
	"time"

	pb "Kaban/internal/Drivers/Grpc/protoFiles"

	hub "Kaban/internal/Drivers/Grpc"
)

type SetNewKeyRequest struct {
	ctx        context.Context
	cancelFunc func()
}

func GetNewSetNewKeyRequest() SetNewKeyRequest {
	return SetNewKeyRequest{}
}

func (s *SetNewKeyRequest) SetNewKeyRequest(data []byte) ([]byte, error) {
	defer s.cancelFunc()
	clientRequest := pb.NewSendingGettingClient(hub.GetConn())
	sec := 2
	for range 12 {
		OutputData, err := clientRequest.GetNewKey(s.ctx, &pb.InputSendData{SendData: data})
		if err != nil {
			slog.Error("SetNewKeyRequest: error make a request", "ERROR", slog.Int("Retry info: seconds", sec))
			sec++
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		return OutputData.BytesOutput, nil
	}

	slog.Error("SetNewKeyRequest; attempts are expired")
	return nil, ErrorAttempts
}

func (s *SetNewKeyRequest) SetAdditionalData(bytes []byte) DomainLevel.MakerKeyRequest {
	return s
}

func (s *SetNewKeyRequest) Make(ctx context.Context) (DomainLevel.Requests, error) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), 7*time.Second)
	}
	s.ctx = ctx
	s.cancelFunc = func() {
		if cancel != nil {
			cancel()
		}
		return
	}
	return s, nil
}
