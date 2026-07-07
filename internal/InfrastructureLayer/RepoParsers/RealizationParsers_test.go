package RepoParsers

import (
	"encoding/json"
	"testing"
)

func TestParsing_JsonParseNotThread(t *testing.T) {
	type ExampleStruct struct {
		Name  string
		Email string
	}
	d := &ExampleStruct{Email: "example@gmail.com", Name: "Pe"}
	convertedData, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	NewStruct := ExampleStruct{}
	type args struct {
		a     any
		bytes []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				a:     &NewStruct,
				bytes: convertedData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parsing{}
			err = p.JsonDecodeMarshall(tt.args.a, tt.args.bytes)
			if err != nil {
				panic(err)
			}
			if NewStruct.Email != d.Email {
				t.Errorf("Data isn't similliar %v,%v", NewStruct.Email, d.Email)
			}

		})
	}
}
