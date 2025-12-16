package common

import (
	"fmt"
	"time"
)

func MessageWithRequestId(message string, id string) string {
	return fmt.Sprintf("%s (request id: %s)", message, id)
}
func GetPointer[T any](v T) *T {
	return &v
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}
