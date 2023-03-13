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
	JobStatuses := []JobStatus{JobStatusScheduled, JobStatusRunning, JobStatusComplete}
	for _, p := range JobStatuses {
		if p == j {
			return true
		}
	}
	return false
}

func (j JobStatus) String() string {
	return [...]string{"scheduled", "running", "complete"}[j-1]
}
