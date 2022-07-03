package util

import "encoding/json"

// JSONEr interface allows return JSON representation
type JSONEr interface {
	ToJSON() *string
}

// ToJSON returns json representation of the given value, or nil in case of marshalling error
func ToJSON[T any](val T) *string {
	bytes, err := json.Marshal(val)
	if err != nil {
		return nil
	}
	value := string(bytes)
	return &value
}
