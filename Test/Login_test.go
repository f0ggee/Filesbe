package Test

import (
	"Kaban/iternal/Controller"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoging(t *testing.T) {

	type LoginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type Test struct {
		ExpectStatus int
		ExpectData   string
		LoginDatas   []byte
	}

	p2 := LoginData{
		Email:    "FERA.com",
		Password: "129221121",
	}

	jsonBody2, err := json.Marshal(p2)
	if err != nil {
		t.Fatalf("Err got %v", err)
	}
	test := []Test{

		{ExpectStatus: 400,
			ExpectData: "breakInExecution",
			LoginDatas: jsonBody2,
		},
	}

	type AnswerLogin struct {
		StatusOfOperation   string `json:"StatusOperation"`
		UrlToRedict         string `json:"UrlToRedict"`
		StatusCodeExucation int    `json:"StatusCodeExucation"`
	}
	req := httptest.NewRequest(http.MethodPost, "/login/api", bytes.NewBuffer(jsonBody2))

	w := httptest.NewRecorder()

	Controller.Login(w, req)

	for _, tt := range test {
		var z AnswerLogin
		if err := json.NewDecoder(w.Body).Decode(&z); err != nil {
			t.Fatalf("Got err when try decode %v", err)
		}

		if z.StatusCodeExucation != tt.ExpectStatus {
			t.Fatalf("Expect but %v got %v", z.StatusCodeExucation, tt.ExpectStatus)
		}
		if tt.ExpectData == "processed" {
			if z.StatusOfOperation != "precessed" {
				t.Fatalf("Expect %v but got %v", tt.ExpectData, z.StatusOfOperation)
			}

		}
		if tt.ExpectData == "breakInExecution" {
			if z.StatusOfOperation != "breakInExecution" {
				t.Fatalf("Expect %v but got %v", tt.ExpectData, z.StatusOfOperation)
			}

		}

	}
}
