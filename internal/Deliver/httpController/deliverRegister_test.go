package httpController

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSession"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RegisterAppTest(ctx context.Context, de *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData {

	return DomainLevel.RegisterApplicationOutComingData{
		Err: nil,
	}
}

func TestNewRegister_Register(t *testing.T) {

	type fields struct {
		dataJson           []byte
		RegisterNet        RegisterNet
		NewRegisterDetails NewRegisterDetails
		RegisterService    func(ctx context.Context, de *Dto.UserDataRegister) DomainLevel.RegisterApplicationOutComingData
	}

	details := NewRegisterDetails{
		Session: RepoSession.Mocks{},
		D:       RepoParsers.GetNewParsing(),
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1",
			fields: fields{
				dataJson: func() []byte {
					data := &Dto.UserDataRegister{
						Name:     "Pavel",
						Email:    "sa2asa@gmail.com",
						Password: "1234567",
					}

					dataJson, _ := json.Marshal(data)

					return dataJson
				}(),
				RegisterNet: RegisterNet{
					W: nil,
					R: nil,
				},
				NewRegisterDetails: details,
				RegisterService:    RegisterAppTest,
			},
		},
		{
			name: "test2",
			fields: fields{
				dataJson: func() []byte {
					data := &Dto.UserDataRegister{
						Name:     "Pavlo",
						Email:    "sa2asa@gmail.com",
						Password: "123",
					}
					dataJson, _ := json.Marshal(data)
					return dataJson

				}(),
				RegisterNet:        RegisterNet{},
				NewRegisterDetails: details,
				RegisterService:    RegisterAppTest,
			},
		},
		{
			name: "test3",
			fields: fields{
				dataJson: func() []byte {
					data := &Dto.UserDataRegister{
						Name:     "P",
						Email:    "",
						Password: "123456",
					}
					dataJson, _ := json.Marshal(data)
					return dataJson

				}(),
				RegisterNet:        RegisterNet{},
				NewRegisterDetails: details,
				RegisterService:    RegisterAppTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			D := NewRegister{
				RegisterNet:        tt.fields.RegisterNet,
				NewRegisterDetails: tt.fields.NewRegisterDetails,
				RegisterService:    tt.fields.RegisterService,
			}

			wri := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(tt.fields.dataJson))
			D.W = wri
			D.R = req
			D.Register()
			t.Logf("Out data is %v\n", wri.Code)
			t.Logf("Out data is %v\n", wri.Body.String())
		})
	}
}
