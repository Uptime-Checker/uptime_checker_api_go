package pkg

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
)

// GetTracingID returns tracing id of the call
func GetTracingID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val := ctx.Value(constant.TracingKey)
	if val == nil {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

// NewTracingID sets tracing id
func NewTracingID(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, constant.TracingKey, GetUniqueString())
	return ctx
}
