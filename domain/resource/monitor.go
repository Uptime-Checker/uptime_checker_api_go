package resource

import "net/http"

func GetMonitorHTTPMethod(method string) int32 {
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

func GetMonitorMethod(method int32) string {
	switch method {
	case 1:
		return http.MethodGet
	case 2:
		return http.MethodPost
	case 3:
		return http.MethodPut
	case 4:
		return http.MethodPatch
	case 5:
		return http.MethodDelete
	}
	return http.MethodHead
}

// MonitorType type
type MonitorType int

// List of types
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

// List of types
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

// List of types
const (
	MonitorBodyFormatNoBody MonitorBodyFormat = iota + 1
	MonitorBodyFormatXML
	MonitorBodyFormatJSON
	MonitorBodyFormatHTML
	MonitorBodyFormatGraphQL
	MonitorBodyFormatFormParam
	MonitorBodyFormatRAW
)

// Valid checks if the MonitorBodyFormat is valid
func (m MonitorBodyFormat) Valid() bool {
	formats := []MonitorBodyFormat{
		MonitorBodyFormatNoBody,
		MonitorBodyFormatXML,
		MonitorBodyFormatJSON,
		MonitorBodyFormatHTML,
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
	return [...]string{
		"", "application/xml", "application/json", "text/html", "application/x-www-form-urlencoded", "",
	}[m-1]
}
