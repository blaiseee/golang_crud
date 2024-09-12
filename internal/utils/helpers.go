package utils

import "encoding/json"

func ParseJSONResponse(body []byte, out interface{}) error {
	return json.Unmarshal(body, out)
}
