package httpController

import (
	"Kaban/internal/Drivers/Tokens"
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadEncryptController_FileUploaderEncrypt(t *testing.T) {
	type fields struct {
		NewFileUploaderEncryptSession NewFileUploaderEncryptSession
		UploaderEncryptApplication    Application.UploaderEncryptApplication
	}
	type args struct {
		W http.ResponseWriter
		R *http.Request
	}

	defualtFields := fields{
		NewFileUploaderEncryptSession: NewFileUploaderEncryptSession{
			ReadSession: newTestSession(),
			Jwt:         Tokens.NewJwtToken(),
			Rf:          Tokens.NewRfToken(),
		},
		UploaderEncryptApplication: testFunctions.TestEncryptUploading(),
	}
	defaultArgs := args{
		W: httptest.NewRecorder(),
		R: newFileUploadRequest("/data", "File", "fileName.txt", []byte("DatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatatDatat")),
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError bool
	}{
		{
			name:      "test1",
			fields:    defualtFields,
			args:      defaultArgs,
			wantError: false,
		},
		{
			name: "test2",
			fields: fields{
				NewFileUploaderEncryptSession: NewFileUploaderEncryptSession{
					ReadSession: newTestBadSession(),
					Jwt:         defualtFields.NewFileUploaderEncryptSession.Jwt,
					Rf:          defualtFields.NewFileUploaderEncryptSession.Rf,
				},
				UploaderEncryptApplication: defualtFields.UploaderEncryptApplication,
			},
			args:      defaultArgs,
			wantError: true,
		},
		{
			name: "test3",
			fields: fields{
				NewFileUploaderEncryptSession: NewFileUploaderEncryptSession{
					ReadSession: defualtFields.NewFileUploaderEncryptSession.ReadSession,
					Jwt:         defualtFields.NewFileUploaderEncryptSession.Jwt,
					Rf:          defualtFields.NewFileUploaderEncryptSession.Rf,
				},
				UploaderEncryptApplication: defualtFields.UploaderEncryptApplication,
			},
			args: args{
				W: defaultArgs.W,
				R: newFileUploadRequest("/data", "file", "f", []byte("data")),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &UploadEncryptController{
				NewFileUploaderEncryptSession: tt.fields.NewFileUploaderEncryptSession,
				UploaderEncryptApplication:    tt.fields.UploaderEncryptApplication,
			}
			S.FileUploaderEncrypt(tt.args.W, tt.args.R)

			var answer AnswerUploadEncrypt
			recorder := tt.args.W.(*httptest.ResponseRecorder)
			if err := json.NewDecoder(recorder.Body).Decode(&answer); err != nil {
				t.Errorf("FileUploaderEncrypt() got error = %v", err)
				return
			}
			if (answer.Error != "") != tt.wantError {
				t.Errorf("FileUploaderEncrypt() wantError = %v but got error = %v", tt.wantError, answer.Error)
				return
			}
		})
	}
}
