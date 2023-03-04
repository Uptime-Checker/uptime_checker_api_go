package validate

// UserLoginProvider type
type UserLoginProvider string

// list of providers
const (
	UserLoginProviderEmail  UserLoginProvider = "email"
	UserLoginProviderGoogle UserLoginProvider = "google"
	UserLoginProviderGithub UserLoginProvider = "github"
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
