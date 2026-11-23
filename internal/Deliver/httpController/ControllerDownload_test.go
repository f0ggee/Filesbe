package httpController

import (
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestDownloadController_Download(t *testing.T) {

	type network struct {
		r *http.Request
		w io.Writer
	}
	type fields struct {
		DownloadNetwork     network
		DownloadApplication Application.DownloadApplication
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
				DownloadNetwork: network{
					w: httptest.NewRecorder(),
					r: setTestRequest(),
				},
				DownloadApplication: testFunctions.GetDownloadTest(bytes.NewReader([]byte("Hello!"))),
			},
			wantError: false,
			inputData: []byte("Hello!"),
		},
		{
			name: "test2",
			fields: fields{
				DownloadNetwork: network{
					r: setTestRequest(),
					w: httptest.NewRecorder(),
				},
				DownloadApplication: testFunctions.GetDownloadErrorTest(),
			},
			wantError: true,
			inputData: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DownloadController{
				DownloadApplication: tt.fields.DownloadApplication,
			}
			d.Download(tt.fields.DownloadNetwork.w.(http.ResponseWriter), tt.fields.DownloadNetwork.r)
			if tt.wantError {
				var answer AnswerDownloadNoEncrypt
				recorder := tt.fields.DownloadNetwork.w.(*httptest.ResponseRecorder)
				if err := json.NewDecoder(recorder.Body).Decode(&answer); err != nil {
					t.Errorf("TestDownload() got error = %v", err)
					return
				}
				if answer.Error != "" {
					t.Errorf("TestDownload() wanted = %v but got = %v", tt.wantError, answer.Error)
					return
				}
			} else {
				recorder := tt.fields.DownloadNetwork.w.(*httptest.ResponseRecorder)
				if recorder.Body.String() != string(tt.inputData) {
					t.Errorf("TestDownload() wanted = %v but got = %v", string(tt.inputData), recorder.Body.String())
					return
				}
			}
		})

	}
}

func setTestRequest() *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/d/example/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "example.txt",
	})

	return req
}
