package RepoParsers

import (
	"Kaban/internal/Dto"
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func TestParsing_JsonParsers(t *testing.T) {
	type args struct {
		a    Dto.UserDataRegister
		body io.ReadCloser
	}

	originalData := Dto.UserDataRegister{
		Name: "Pavel",
	}

	convertedData, _ := json.Marshal(originalData)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				body: io.NopCloser(bytes.NewReader(convertedData)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserData := &Dto.UserDataRegister{}
			p := Parsing{}
			err := p.JsonDecode(UserData, tt.args.body)
			if err != nil {
				t.Errorf("ERROR to pass the test %v\n", err)
			}
			if UserData.Name != originalData.Name {
				t.Errorf("incoming data and parsed data aren't alike\n")
			}

		})
	}
}
