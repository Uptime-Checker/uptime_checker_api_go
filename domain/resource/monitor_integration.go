package resource

// MonitorIntegrationType type
type MonitorIntegrationType int

// List of types
const (
	MonitorIntegrationTypeEmail MonitorIntegrationType = iota + 1
	MonitorIntegrationTypeWebhook
	MonitorIntegrationTypeSlack
	MonitorIntegrationTypeTeams
	MonitorIntegrationTypeDiscord
)

// Valid checks if the MonitorType is valid
func (m MonitorIntegrationType) Valid() bool {
	integrations := []MonitorIntegrationType{
		MonitorIntegrationTypeEmail,
		MonitorIntegrationTypeWebhook,
		MonitorIntegrationTypeSlack,
		MonitorIntegrationTypeTeams,
		MonitorIntegrationTypeDiscord,
	}
	for _, p := range integrations {
		if p == m {
			return true
		}
	}
	return false
}

func (m MonitorIntegrationType) String() string {
	return [...]string{"email", "webhook", "slack", "teams", "discord"}[m-1]
}
