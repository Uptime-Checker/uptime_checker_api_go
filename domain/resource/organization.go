package resource

// RoleType type
type RoleType int

// list of role types
const (
	RoleTypeSuperAdmin RoleType = iota + 1
	RoleTypeAdmin
	RoleTypeEditor
	RoleTypeMember
)

// Valid checks if the RoleType is valid
func (r RoleType) Valid() bool {
	roleTypes := []RoleType{
		RoleTypeSuperAdmin,
		RoleTypeAdmin,
		RoleTypeEditor,
		RoleTypeMember,
	}
	for _, c := range roleTypes {
		if c == r {
			return true
		}
	}
	return false
}

func (r RoleType) String() string {
	return [...]string{"SuperAdmin", "Admin", "Editor", "Member"}[r-1]
}
