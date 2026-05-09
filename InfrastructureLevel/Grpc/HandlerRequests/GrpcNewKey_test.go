package HandlerRequests

import (
	"MasterServer_/DomainLevel"
	"MasterServer_/InfrastructureLevel/Grpc/Proto/protoFiles"
	"crypto/rand"
	"testing"
	"time"
)

func TestGrpcHandlerGettingNewKey_CalculateSwapingTime1(t *testing.T) {
	type fields struct {
		UnimplementedSendingGettingServer protoFiles.UnimplementedSendingGettingServer
		S                                 *HandlingRequestsForNewKey
	}
	S := DomainLevel.Get()

	NewCollect := &HandlingRequestsForNewKey{
		Time: S,
	}
	S.NewPreviousTime(rand.Text(), time.Now())
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "test_1",
			fields: fields{
				S: NewCollect,
			},

			want: func() time.Duration {

				Sa := time.Now().Add(12 * time.Hour)

				Xz := time.Until(Sa)
				return Xz
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := GrpcHandlerGettingNewKey{
				UnimplementedSendingGettingServer: tt.fields.UnimplementedSendingGettingServer,
				S:                                 tt.fields.S,
			}
			if got := S.CalculateSwapingTime(); got != tt.want {
				t.Errorf("CalculateSwapingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcHandlerGettingNewKey_CalculateSwapingTimeTestFunc(t *testing.T) {
	type fields struct {
		UnimplementedSendingGettingServer protoFiles.UnimplementedSendingGettingServer
		S                                 *HandlingRequestsForNewKey
	}
	S := DomainLevel.Get()

	NewCollect := &HandlingRequestsForNewKey{
		Time: S,
	}

	TimeForTests := time.Now()
	S.NewPreviousTime(rand.Text(), TimeForTests)

	tests := []struct {
		name   string
		fields fields
		args   time.Duration
		want   time.Duration
	}{
		{
			name: "test_1",
			fields: fields{
				S: NewCollect,
			},

			args: 10 * time.Second,

			want: func() time.Duration {

				sax := TimeForTests.Add(10 * time.Second)
				x := time.Until(sax)
				return x
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := GrpcHandlerGettingNewKey{
				UnimplementedSendingGettingServer: tt.fields.UnimplementedSendingGettingServer,
				S:                                 tt.fields.S,
			}
			if got := S.CalculateSwapingTimeTestFunc(tt.args); got.Round(time.Second) != tt.want.Round(time.Second) {
				t.Errorf("CalculateSwapingTimeTestFunc() = %v, want %v", got.String(), tt.want.Round(time.Second).String())
			}
		})
	}
}
