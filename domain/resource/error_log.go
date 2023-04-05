package resource

// ErrorLogType type
type ErrorLogType int

// List of error log types
const (
	ErrorLogTypeTimeout ErrorLogType = iota + 1
	ErrorLogTypeAddr
	ErrorLogTypeURL
	ErrorLogTypeDNS
	ErrorLogTypeInvalidAddr
	ErrorLogTypeParse
	ErrorLogTypeUnknownNetwork
	ErrorLogTypeSyscall
	ErrorLogTypeResponseMalformed
	ErrorLogTypeAssertionFailure
	ErrorLogTypeUnknown ErrorLogType = 99
)

// Valid checks if the ErrorLogType is valid
func (e ErrorLogType) Valid() bool {
	errorLogTypes := []ErrorLogType{
		ErrorLogTypeTimeout,
		ErrorLogTypeAddr,
		ErrorLogTypeURL,
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
	if e == ErrorLogTypeUnknown {
		return "Unknown"
	}
	all := [...]string{
		"Timeout", "Addr", "URL", "DNS", "Invalid Addr", "Parse", "Unknown Network", "Syscall", "Response Malformed",
	}
	index := e - 1
	return all[index]
}
