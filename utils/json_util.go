package utils

import jsoniter "github.com/json-iterator/go"

func MarshalAnyToString(param any) string {
	s, err := jsoniter.MarshalToString(param)
	if err != nil {
		return "{}"
	}
	return s
}
