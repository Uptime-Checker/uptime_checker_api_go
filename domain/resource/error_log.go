package resource

// ErrorLogType type
type ErrorLogType int

// list of error log types
const (
	ErrorLogTypeTimeout ErrorLogType = iota + 1
	ErrorLogTypeAddr
	ErrorLogTypeDNS
	ErrorLogTypeInvalidAddr
	ErrorLogTypeParse
	ErrorLogTypeUnknownNetwork
	ErrorLogTypeSyscall
	ErrorLogTypeResponseMalformed
	ErrorLogTypeUnknown ErrorLogType = 99
)

// Valid checks if the ErrorLogType is valid
func (e ErrorLogType) Valid() bool {
	errorLogTypes := []ErrorLogType{
		ErrorLogTypeTimeout,
		ErrorLogTypeAddr,
		ErrorLogTypeDNS,
		ErrorLogTypeInvalidAddr,
		ErrorLogTypeParse,
		ErrorLogTypeUnknownNetwork,
		ErrorLogTypeSyscall,
		ErrorLogTypeResponseMalformed,
		ErrorLogTypeUnknown,
	}
	for _, p := range errorLogTypes {
		if p == e {
			return true
		}
	}
	return false
}

func (e ErrorLogType) String() string {
	return [...]string{
		"Timeout", "Addr", "DNS", "Invalid Addr", "Parse", "Unknown Network", "Syscall", "Response Malformed", "Unknown",
	}[e-1]
}
