package resource

// AlarmPolicyName type
type AlarmPolicyName string

// list of alarm policies
const (
	AlarmPolicyErrorThreshold    AlarmPolicyName = "MAX_ERROR_THRESHOLD"
	AlarmPolicyDurationThreshold AlarmPolicyName = "MAX_DURATION_THRESHOLD"
	AlarmPolicyRegionThreshold   AlarmPolicyName = "REGION_THRESHOLD"
)

// Valid checks if the AlarmPolicyName is valid
func (a AlarmPolicyName) Valid() bool {
	alarmPolicies := []AlarmPolicyName{
		AlarmPolicyErrorThreshold,
		AlarmPolicyDurationThreshold,
		AlarmPolicyRegionThreshold,
	}
	for _, p := range alarmPolicies {
		if p == a {
			return true
		}
	}
	return false
}
