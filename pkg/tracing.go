package pkg

import (
	"context"

	"github.com/segmentio/ksuid"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
)

// GetTracingID returns tracing id of the call
func GetTracingID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val := ctx.Value(constant.TracingKey.String())
	if val == nil {
		return ""
	}
	s, _ := val.(string)
	return s
}

// NewTracingID sets tracing id
func NewTracingID(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, constant.TracingKey.String(), GetUniqueString())
	return ctx
}

// GetUniqueString returns unique string
func GetUniqueString() string {
	return ksuid.New().String()
}
