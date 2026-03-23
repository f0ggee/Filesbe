package DataConverting

import (
	"encoding/json"
)

func (k Converter) JsonConverter(a any) ([]byte, error) {

	JsonDataType, err := json.Marshal(&a)

	if err != nil {
		return nil, err
	}

	return JsonDataType, err

}
