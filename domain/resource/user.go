package resource

// UserLoginProvider type
type UserLoginProvider int

// list of providers
const (
	UserLoginProviderEmail UserLoginProvider = iota + 1
	UserLoginProviderGoogle
	UserLoginProviderGithub
)

// Valid checks if the MediaType is valid
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
