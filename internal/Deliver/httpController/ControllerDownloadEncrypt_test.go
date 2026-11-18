package httpController

import (
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestDownloadEncryptController_DownloadEncrypt(t *testing.T) {
	type net struct {
		W io.Writer
		R *http.Request
	}
	type fields struct {
		Net                        net
		DownloadEncryptApplication func([]byte) Application.DownloadEncryptApplication
	}
	tests := []struct {
		name      string
		fields    fields
		wantError bool
		inputData []byte
	}{
		{
			name: "test1",
			fields: fields{
				Net: net{
					W: httptest.NewRecorder(),
					R: setTestRequest(),
				},
				DownloadEncryptApplication: testFunctions.DownloadEncryptTest,
			},
			wantError: false,
			inputData: []byte("it's example! !!"),
		},
		{
			name: "test2",
			fields: fields{
				Net: net{
					W: httptest.NewRecorder(),
					R: setTestRequest(),
				},
				DownloadEncryptApplication: testFunctions.DownloadEncryptError1,
			},
			wantError: true,
			inputData: []byte("data"),
		},
		{
			name: "test3",
			fields: fields{
				Net: net{
					W: httptest.NewRecorder(),
					R: setContextRequest(),
				},
				DownloadEncryptApplication: testFunctions.DownloadEncryptTest,
			},
			wantError: true,
			inputData: []byte("data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.fields.DownloadEncryptApplication(tt.inputData)
			d := DownloadEncryptController{
				Net: NetworkDownloadEncrypt{
					W: tt.fields.Net.W.(http.ResponseWriter),
					R: tt.fields.Net.R,
				},
				DownloadEncryptApplication: app,
			}

			d.DownloadEncrypt()
			if tt.wantError {
				var answer AnswerFileDownloadEncrypt
				recorder := tt.fields.Net.W.(*httptest.ResponseRecorder)
				if err := json.NewDecoder(recorder.Body).Decode(&answer); err != nil {
					t.Errorf("TestDownload() got error = %v", err)
					return
				}
				if answer.Error == "" {
					t.Errorf("TestDownload() wanted = %v but got = %v", tt.wantError, answer.Error)
					return
				}
				t.Logf("TestDownload()  %v  %v", tt.wantError, answer.Error)
			} else {
				recorder := tt.fields.Net.W.(*httptest.ResponseRecorder)
				if recorder.Body.String() != string(tt.inputData) {
					t.Errorf("TestDownload() wanted = %v but got = %v", string(tt.inputData), recorder.Body.String())
					return
				}
			}
		})
	}
}
func setContextRequest() *http.Request {

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Millisecond)
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/d/example/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "example.txt",
	})

	return req
}
