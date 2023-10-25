package recovery

import (
	"context"
	"google.golang.org/grpc"
)

var (
	defaultOptions = &options{
		recoveryHandlerFunc: nil,
	}
)

type options struct {
	recoveryHandlerFunc RecoveryHandlerFuncContext
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithRecoveryHandler customizes the function for recovering from a panic.
func WithRecoveryHandler(f RecoveryHandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = RecoveryHandlerFuncContext(func(ctx context.Context, p any, req any, info *grpc.UnaryServerInfo) error {
			return f(p)
		})
	}
}

// WithRecoveryHandlerContext customizes the function for recovering from a panic.
func WithRecoveryHandlerContext(f RecoveryHandlerFuncContext) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
