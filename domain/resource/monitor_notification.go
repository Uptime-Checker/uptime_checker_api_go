package resource

// MonitorNotificationType type
type MonitorNotificationType int

// List of types
const (
	MonitorIntegrationTypeRaiseAlarm MonitorNotificationType = iota + 1
	MonitorIntegrationTypeResolveAlarm
)

// Valid checks if the MonitorType is valid
func (m MonitorNotificationType) Valid() bool {
	integrations := []MonitorNotificationType{
		MonitorIntegrationTypeRaiseAlarm,
		MonitorIntegrationTypeResolveAlarm,
	}
	for _, p := range integrations {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorNotificationType) String() string {
	return [...]string{"alarm:raise", "alarm:resolve"}[m-1]
}
