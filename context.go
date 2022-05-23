package cute

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type testCtxKey struct{}

func withProviderT(parentCtx context.Context, t internalT) context.Context {
	return context.WithValue(parentCtx, testCtxKey{}, t)
}

func getProviderT(ctx context.Context) (internalT, error) {
	testCtx := ctx.Value(testCtxKey{})
	if testCtx == nil {
		errMsg := "TestContext did not passed with context.getProviderT method"
		return nil, errors.New(errMsg)
	}
	if reflect.ValueOf(testCtx).Kind() != reflect.Ptr {
		errMsg := fmt.Sprintf("Value is not a pointer (%v)", testCtx)
		return nil, errors.New(errMsg)
	}
	if p, ok := testCtx.(internalT); ok {
		return p, nil
	}
	errMsg := fmt.Sprintf("Wrong pointer type. Expected: (provider.T). Actual: (%v)", testCtx)
	return nil, errors.New(errMsg)
}
