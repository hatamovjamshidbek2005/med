package serialize

import (
	"encoding/json"
)

func MarshalUnMarshal(input any, output any) error {
	rawBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(rawBytes, output)
}

func StructToMapViaJson(payload interface{}) (map[string]interface{}, error) {
	mapBasedStruct := make(map[string]interface{})

	j, _ := json.Marshal(payload)

	err := json.Unmarshal(j, &mapBasedStruct)
	if err != nil {
		return nil, err
	}

	return mapBasedStruct, nil
}
