package HandlerRequests

import (
	"MasterServer_/DomainLevel"
	"crypto/rand"
	"log/slog"
	"testing"
	"time"
)

func TestGrpcHandlerGettingNewKey_CalculateSwapingTime(t *testing.T) {
	type fields struct {
		S HandlingRequestsForNewKey
	}

	sa := DomainLevel.Get()
	sa.NewPreviousTime(rand.Text(), time.Now())
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Get New Key",
			fields: fields{
				S: HandlingRequestsForNewKey{
					Time: *sa,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := GrpcHandlerGettingNewKey{
				S: tt.fields.S,
			}
			got := S.CalculateSwapingTime()
			slog.Info("Func CalculateSwapingTime", "Time for next swaping", got)
		})
	}
}
