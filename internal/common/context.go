package common

import (
	"context"
	"fmt"

	"reflect"
	"unsafe"

	"net/http"
	"net/url"
)

func GetCtxValue(r *http.Request, key string) (interface{}, error) {
	val := r.Context().Value(key)
	if val == nil {
		return nil, fmt.Errorf("no value exists in the context for key %q", key)
	}
	return val, nil
}

// The type is VERY important. If you just use a string in the below functions - it won't actually match!
type CtxKey string

const (
	CSRFCtxKey             CtxKey = "csrf"
	ToastsCtxKey           CtxKey = "toastsCtx"
	CurrentUrlCtxKey       CtxKey = "currentUrlCtx"
	CurrentUrlParamsCtxKey CtxKey = "currentUrlParamsCtx"
	HotReloadUrlCtxKey     CtxKey = "hotreloadUrlCtx"
)

func SetCtxValueOnRequest(r *http.Request, key CtxKey, val interface{}) *http.Request {
	ctx := r.Context()
	newCtx := context.WithValue(ctx, key, val)
	return r.WithContext(newCtx)
}

// func GetCtxToasts(ctx context.Context) []ignite.FlashMessage {
// 	if toasts, ok := ctx.Value(ToastsCtxKey).([]ignite.FlashMessage); ok {
// 		return toasts
// 	}
// 	return make([]ignite.FlashMessage, 0)
// }

func GetCtxCSRF(ctx context.Context) string {
	if csrf, ok := ctx.Value(CSRFCtxKey).(string); ok {
		return csrf
	}

	return ""
}

func GetCtxCurrentUrl(ctx context.Context) string {
	if url, ok := ctx.Value(CurrentUrlCtxKey).(string); ok {
		return url
	}

	return ""
}

func GetCtxHotReloadURL(ctx context.Context) string {
	if url, ok := ctx.Value(HotReloadUrlCtxKey).(string); ok {
		return url
	}
	// random whatever
	return "http://localhost:8082"
}

func GetCtxCurrentUrlParams(ctx context.Context) url.Values {
	if url, ok := ctx.Value(CurrentUrlParamsCtxKey).(url.Values); ok {
		return url
	}

	return url.Values{}
}

// Debug only helper functi0on
func printContextInternals(ctx interface{}, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				printContextInternals(reflectValue.Interface(), true)
			} else {
				fmt.Printf("field name: %+v\n", reflectField.Name)
				fmt.Printf("value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("context is empty (int)\n")
	}
}
