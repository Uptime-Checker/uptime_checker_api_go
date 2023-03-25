package resource

import "net/http"

func GetMonitorMethod(method string) int32 {
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

// Valid checks if the MonitorType is valid
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

// MonitorStatus type
type MonitorStatus int

// list of types
const (
	MonitorStatusPending MonitorStatus = iota + 1
	MonitorStatusPassing
	MonitorStatusDegraded
	MonitorStatusFailing
)

// Valid checks if the MonitorStatus is valid
func (m MonitorStatus) Valid() bool {
	statuses := []MonitorStatus{MonitorStatusPending, MonitorStatusPassing, MonitorStatusDegraded, MonitorStatusFailing}
	for _, p := range statuses {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorStatus) String() string {
	return [...]string{"pending", "passing", "degraded", "failing"}[m-1]
}

// MonitorBodyFormat type
type MonitorBodyFormat int

// list of types
const (
	MonitorBodyFormatNoBody MonitorBodyFormat = iota + 1
	MonitorBodyFormatJSON
	MonitorBodyFormatGraphQL
	MonitorBodyFormatFormParam
	MonitorBodyFormatRAW
)

// Valid checks if the MonitorBodyFormat is valid
func (m MonitorBodyFormat) Valid() bool {
	formats := []MonitorBodyFormat{
		MonitorBodyFormatNoBody,
		MonitorBodyFormatJSON,
		MonitorBodyFormatGraphQL,
		MonitorBodyFormatFormParam,
		MonitorBodyFormatRAW,
	}
	for _, p := range formats {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorBodyFormat) String() string {
	return [...]string{"no-body", "json", "graphql", "form-parameters", "raw"}[m-1]
}
