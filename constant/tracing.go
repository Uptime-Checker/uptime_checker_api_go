package constant

type contextKey string

const (
	TracingKey contextKey = "tracing"
	UserKey    contextKey = "user"
)

func (c contextKey) String() string {
	return string(c)
}
