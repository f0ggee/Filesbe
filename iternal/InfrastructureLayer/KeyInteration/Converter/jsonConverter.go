package Converter

import "encoding/json"

func (k KeyConverter) JsonConverter(a any) ([]byte, error) {

	JsonDataType, err := json.Marshal(&a)

	if err != nil {
		return nil, err
	}

	return JsonDataType, err

}
