package resource

// JobStatus type
type JobStatus int

// list of providers
const (
	JobStatusScheduled JobStatus = iota + 1
	JobStatusRunning
	JobStatusComplete
)

// Valid checks if the UserLoginProvider is valid
func (j JobStatus) Valid() bool {
	jobStatuses := []JobStatus{JobStatusScheduled, JobStatusRunning, JobStatusComplete}
	for _, p := range jobStatuses {
		if p == j {
			return true
		}
	}
	return false
}

func (j JobStatus) String() string {
	return [...]string{"scheduled", "running", "complete"}[j-1]
}
