package pkg

import "github.com/segmentio/ksuid"

func IsEmpty(str string) bool {
	return str == ""
}

// GetUniqueString returns unique string
func GetUniqueString() string {
	return ksuid.New().String()
}
