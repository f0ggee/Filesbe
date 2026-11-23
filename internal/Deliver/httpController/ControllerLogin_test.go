package httpController

import (
	"Kaban/internal/Deliver/httpController/httpData"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Drivers/Parsers"
	"Kaban/internal/Service/Application"
	"Kaban/internal/Service/Application/testFunctions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginController_LoginController(t *testing.T) {
	type fields struct {
		Sess   Session
		Parses DomainLevel.Decoder
		App    Application.LoginApplication
	}

	type args struct {
		W http.ResponseWriter
		R *http.Request
	}
	defaultFields := fields{
		Sess:   newTestSession(),
		Parses: Parsers.GetNewParsing(),
		App:    testFunctions.NewTestLogin(),
	}
	defaultArgs := args{
		W: httptest.NewRecorder(),
		R: httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(func() []byte {
			a := httpData.LoginData{
				Email:    "33123@gmail.com",
				Password: []byte("123331232131231"),
			}

			d, _ := json.Marshal(a)
			return d
		}())),
	}
	tests := []struct {
		name      string
		fields    fields
		arg       args
		wantError bool
	}{

		{
			name:      "test1",
			fields:    defaultFields,
			arg:       defaultArgs,
			wantError: false,
		},
		{
			name:   "test2",
			fields: defaultFields,
			arg: args{
				W: defaultArgs.W,
				R: httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(func() []byte {
					a := httpData.LoginData{
						Email:    "21312321313123131232132323j21jdsjfjjffjfjfjjjjjjjjfffffffffffffffffffffffffffffffffffffffff@gmail.com",
						Password: []byte("123331232131231"),
					}
					marshal, _ := json.Marshal(a)
					return marshal

				}())),
			},
			wantError: true,
		},
		{
			name:   "test3",
			fields: defaultFields,
			arg: args{
				W: httptest.NewRecorder(),
				R: httptest.NewRequest(http.MethodPost, "/data", bytes.NewReader(func() []byte {

					h := httpData.LoginData{
						Email:    "",
						Password: []byte("123331232131231"),
					}
					marshal, _ := json.Marshal(h)
					return marshal
				}())),
			},
			wantError: true,
		},
		{
			name:   "test4",
			fields: defaultFields,
			arg: args{
				W: httptest.NewRecorder(),
				R: httptest.NewRequest(http.MethodPost, "/data", bytes.NewReader(func() []byte {
					h := httpData.LoginData{
						Email:    "",
						Password: []byte{},
					}
					marshal, _ := json.Marshal(h)
					return marshal
				}())),
			},
			wantError: true,
		},
		{
			name:   "test5",
			fields: defaultFields,
			arg: args{
				W: httptest.NewRecorder(),
				R: httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(func() []byte {
					h := httpData.LoginData{
						Email:    "12312321@gmail.com",
						Password: []byte("123331232131231123331232131231123331232131231123331232131231123331232131231123331232131231123331232131231"),
					}
					marshal, _ := json.Marshal(h)
					return marshal
				}())),
			},
			wantError: true,
		},
		{
			name:   "test5",
			fields: defaultFields,
			arg: args{
				W: httptest.NewRecorder(),
				R: httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(func() []byte {

					data := httpData.RegisterData{
						Name:     "123",
						Email:    "123",
						Password: []byte("1231321"),
					}
					marshal, _ := json.Marshal(data)
					return marshal
				}())),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			D := &LoginController{
				Sess:   tt.fields.Sess,
				Parses: tt.fields.Parses,
				App:    tt.fields.App,
			}
			D.LoginController(tt.arg.W, tt.arg.R)

			var answer AnswerLogin
			recorder := tt.arg.W.(*httptest.ResponseRecorder)
			if err := json.NewDecoder(recorder.Body).Decode(&answer); err != nil {
				t.Errorf("LoginController() got error = %v", err)
				return
			}
			if (answer.ErrorMessage != "") != tt.wantError {
				t.Errorf("LoginController() wantError = %v but got error = %v", tt.wantError, answer.ErrorMessage)
				return
			}
		})
	}
}
