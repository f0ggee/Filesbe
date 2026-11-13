package httpController

import (
	"Kaban/internal/Deliver/httpController/httpData"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterController_Register(t *testing.T) {

	type registerNet struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type fields struct {
		net     registerNet
		Session Session
		Decoder DomainLevel.Decoder
		App     Application.RegisterApplication
	}
	tests := []struct {
		name      string
		fields    fields
		wantError bool
	}{
		{
			name: "test1",
			fields: fields{
				net: registerNet{
					w: httptest.NewRecorder(),
					r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/test", bytes.NewReader(func() []byte {
						data := httpData.RegisterData{
							Name:     "1234",
							Email:    "12345@gmail.com",
							Password: []byte("12345512331"),
						}
						d, _ := json.Marshal(data)
						return d
					}())),
				},
				Session: GetNewMocks(),
				Decoder: Parsers.GetNewParsing(),
				App:     testFunctions.GetTestRegister(),
			},
			wantError: false,
		},
		{
			name: "test2",
			fields: fields{
				net: registerNet{
					w: httptest.NewRecorder(),
					r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/test", bytes.NewReader(func() []byte {
						data := httpData.RegisterData{
							Name:     "1234",
							Email:    "123.com",
							Password: []byte("12345512331"),
						}
						d, _ := json.Marshal(data)
						return d

					}())),
				},
				Session: GetNewMocks(),
				Decoder: Parsers.GetNewParsing(),
				App:     testFunctions.GetRegisterError(),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			D := RegisterController{
				Session: tt.fields.Session,
				Decoder: tt.fields.Decoder,
				App:     tt.fields.App,
			}
			D.RegisterNet.W = tt.fields.net.w
			D.RegisterNet.R = tt.fields.net.r
			D.Register()

			var answer AnswerRegister

			if err := json.NewDecoder(tt.fields.net.w.Body).Decode(&answer); err != nil {
				t.Errorf("TestRegister() got error = %v", err)
				return
			}

			if (tt.wantError == false) != (len(answer.Error) == 0) {
				t.Errorf("TestRegister() wanted = %v but got = %v ", tt.wantError, answer.Error)
				return
			}
		})
	}
}
