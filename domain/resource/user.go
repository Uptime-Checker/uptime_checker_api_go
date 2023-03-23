package resource

// UserLoginProvider type
type UserLoginProvider int

// list of providers
const (
	UserLoginProviderEmail UserLoginProvider = iota + 1
	UserLoginProviderGoogle
	UserLoginProviderGithub
)

// Valid checks if the UserLoginProvider is valid
func (u UserLoginProvider) Valid() bool {
	providers := []UserLoginProvider{UserLoginProviderEmail, UserLoginProviderGoogle, UserLoginProviderGithub}
	for _, p := range providers {
		if p == u {
			return true
		}
	}
	return false
}

func (u UserLoginProvider) String() string {
	return [...]string{"email", "google", "github"}[u-1]
}

// UserContactMode type
type UserContactMode int

// list of contact modes
const (
	UserContactModeEmail UserContactMode = iota + 1
	UserContactModeSMS
	UserContactModePhone
	UserContactModeSMSAndPhone
)

// Valid checks if the UserContactMode is valid
func (u UserContactMode) Valid() bool {
	contactModes := []UserContactMode{
		UserContactModeEmail,
		UserContactModeSMS,
		UserContactModePhone,
		UserContactModeSMSAndPhone,
	}
	for _, c := range contactModes {
		if c == u {
			return true
		}
	}
	return false
}

func (u UserContactMode) String() string {
	return [...]string{"email", "sms", "phone", "sms-and-phone"}[u-1]
}
