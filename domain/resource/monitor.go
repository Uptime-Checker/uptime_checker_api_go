package resource

import "net/http"

func GetMonitorMethod(method string) int {
	switch method {
	case http.MethodGet:
		return 1
	case http.MethodPost:
		return 2
	case http.MethodPut:
		return 3
	case http.MethodPatch:
		return 4
	case http.MethodDelete:
		return 5
	}
	return 0
}

// MonitorType type
type MonitorType int

// list of types
const (
	MonitorTypeAPI MonitorType = iota + 1
	MonitorTypeBrowser
	MonitorTypeAPISnapshot
)

// Valid checks if the UserLoginProvider is valid
func (m MonitorType) Valid() bool {
	providers := []MonitorType{MonitorTypeAPI, MonitorTypeBrowser, MonitorTypeAPISnapshot}
	for _, p := range providers {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorType) String() string {
	return [...]string{"api", "browser", "apiSnapshot"}[m-1]
}
