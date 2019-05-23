package jsonHttpHandler

import (
	"context"
	"fmt"
	"github.com/renra/go-errtrace/errtrace"
)

const (
	PathParamsKey = "path_params"
	PayloadKey    = "payload"
)

func GetPathParam(ctx context.Context, key string) (*string, *errtrace.Error) {
	urlParams, ok := ctx.Value(PathParamsKey).(map[string]string)

	if !ok {
		return nil, errtrace.New(fmt.Sprintf("Could not find url params", key))
	}

	value, ok := urlParams[key]

	if !ok {
		return nil, errtrace.New(fmt.Sprintf("Could not find key %s in url params", key))
	}

	return &value, nil
}

func GetPathParamP(ctx context.Context, key string) string {
	value, err := GetPathParam(ctx, key)

	if err != nil {
		panic(err)
	}

	return *value
}
