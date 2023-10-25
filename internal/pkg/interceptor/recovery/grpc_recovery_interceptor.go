package recovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"runtime"
)

type RecoveryHandlerFunc func(p any) (err error)

type RecoveryHandlerFuncContext func(ctx context.Context, p any, req any, info *grpc.UnaryServerInfo) (err error)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, r, req, info, o.recoveryHandlerFunc)
			}
		}()

		return handler(ctx, req)
	}
}

func recoverFrom(ctx context.Context, p any, req any, info *grpc.UnaryServerInfo, r RecoveryHandlerFuncContext) error {
	if r != nil {
		return r(ctx, p, req, info)
	}
	stack := make([]byte, 64<<10)
	stack = stack[:runtime.Stack(stack, false)]
	return &PanicError{Panic: p, Stack: stack}
}

type PanicError struct {
	Panic any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught: %v\n\n%s", e.Panic, e.Stack)
}
