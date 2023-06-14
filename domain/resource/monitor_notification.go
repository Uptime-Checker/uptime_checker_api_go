package resource

// MonitorNotificationType type
type MonitorNotificationType int

// List of types
const (
	MonitorNotificationTypeMonitorUp MonitorNotificationType = iota + 1
	MonitorNotificationTypeMonitorDown
)

// Valid checks if the MonitorType is valid
func (m MonitorNotificationType) Valid() bool {
	integrations := []MonitorNotificationType{
		MonitorNotificationTypeMonitorUp,
		MonitorNotificationTypeMonitorDown,
	}
	for _, p := range integrations {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorNotificationType) String() string {
	return [...]string{"monitor:up", "monitor:down"}[m-1]
}
