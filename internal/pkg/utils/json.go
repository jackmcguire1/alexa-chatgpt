package utils

import "encoding/json"

func ToJSON(i any) string {
	res, _ := json.Marshal(i)
	return string(res)
}
