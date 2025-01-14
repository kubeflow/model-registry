package grpc

import (
	"fmt"
	"reflect"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// maxRetryAttempts is the maximum number of times to retry a request.
	maxRetryAttempts = 25
)

var retryableStatusCodes = map[codes.Code]bool{
	codes.Unavailable: true,
}

func RetryOnGRPCError[T any](funcToRetry any, funcParams ...any) (T, error) {
	var outErr error
	var res T

	fnVal := reflect.ValueOf(funcToRetry)
	if fnVal.Kind() != reflect.Func {
		return res, fmt.Errorf("grpc retry error: function parameter is not a function")
	}

	if len(funcParams) != fnVal.Type().NumIn() {
		return res, fmt.Errorf("grpc retry error: function parameters count mismatch")
	}

	inputs := make([]reflect.Value, len(funcParams))
	for i, param := range funcParams {
		inputs[i] = reflect.ValueOf(param)
	}

	for i := 0; i < maxRetryAttempts; i++ {
		outs := fnVal.Call(inputs)

		if len(outs) != 2 {
			return res, fmt.Errorf("grpc retry error: function from input does not return 2 values")
		}

		res = outs[0].Interface().(T)
		errI := outs[1].Interface()
		if errI == nil {
			outErr = nil
		} else {
			outErr = errI.(error)
		}

		if status, ok := status.FromError(outErr); ok {
			if !retryableStatusCodes[status.Code()] {
				break
			}
		} else {
			break
		}

		backoff := time.Duration(i+1) * time.Second
		time.Sleep(backoff)
	}

	return res, outErr
}
