package cute

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFillBaseProps_WhenBasePropsIsNil(t *testing.T) {
	testObj := &Test{}
	cuteObj := &cute{}

	cuteObj.fillBaseProps(testObj)

	require.Nil(t, testObj.httpClient)
	require.Nil(t, testObj.jsonMarshaler)
	require.Nil(t, testObj.Middleware)
}

func TestFillBaseProps_WhenBasePropsIsNotNil(t *testing.T) {
	testObj := &Test{}
	cuteObj := &cute{}

	qtBaseProps := &HTTPTestMaker{
		httpClient:    &http.Client{},
		jsonMarshaler: &jsonMarshaler{},
		middleware: &Middleware{
			After: []AfterExecute{
				func(*http.Response, []error) error {
					return nil
				},
				func(*http.Response, []error) error {
					return nil
				},
			},
			AfterT: []AfterExecuteT{func(T, *http.Response, []error) error { return nil }},
			Before: []BeforeExecute{
				func(*http.Request) error {
					return nil
				},
				func(*http.Request) error {
					return nil
				},
			},
			BeforeT: []BeforeExecuteT{
				func(T, *http.Request) error { return nil },
				func(t T, request *http.Request) error {
					return nil
				},
			},
		},
	}
	cuteObj.baseProps = qtBaseProps

	cuteObj.fillBaseProps(testObj)

	require.Equal(t, qtBaseProps.httpClient, testObj.httpClient)
	require.Equal(t, qtBaseProps.jsonMarshaler, testObj.jsonMarshaler)
	require.Len(t, testObj.Middleware.After, len(qtBaseProps.middleware.After))
	require.Len(t, testObj.Middleware.AfterT, len(qtBaseProps.middleware.AfterT))
	require.Len(t, testObj.Middleware.Before, len(qtBaseProps.middleware.Before))
	require.Len(t, testObj.Middleware.BeforeT, len(qtBaseProps.middleware.BeforeT))
}

func TestFillBaseProps_WhenBasePropsIsNotNil_After(t *testing.T) {
	testObj := &Test{
		Middleware: &Middleware{
			After: []AfterExecute{
				func(response *http.Response, errors []error) error {
					return nil
				},
			},
		},
	}
	cuteObj := &cute{}

	qtBaseProps := &HTTPTestMaker{
		httpClient:    &http.Client{},
		jsonMarshaler: &jsonMarshaler{},
		middleware: &Middleware{
			After: []AfterExecute{
				func(*http.Response, []error) error {
					return nil
				},
				func(*http.Response, []error) error {
					return nil
				},
			},
			BeforeT: []BeforeExecuteT{
				func(T, *http.Request) error { return nil },
				func(t T, request *http.Request) error {
					return nil
				},
			},
		},
	}
	cuteObj.baseProps = qtBaseProps

	cuteObj.fillBaseProps(testObj)

	require.Equal(t, qtBaseProps.httpClient, testObj.httpClient)
	require.Equal(t, qtBaseProps.jsonMarshaler, testObj.jsonMarshaler)
	require.Len(t, testObj.Middleware.After, len(qtBaseProps.middleware.After)+1)
	require.Len(t, testObj.Middleware.AfterT, len(qtBaseProps.middleware.AfterT))
	require.Len(t, testObj.Middleware.Before, len(qtBaseProps.middleware.Before))
	require.Len(t, testObj.Middleware.BeforeT, len(qtBaseProps.middleware.BeforeT))
}

func TestFillBaseProps_WhenBasePropsIsNotNil_Middleware(t *testing.T) {
	testObj := &Test{
		Middleware: &Middleware{
			After: []AfterExecute{
				func(response *http.Response, errors []error) error {
					return nil
				},
				func(response *http.Response, errors []error) error {
					return nil
				},
			},
			AfterT: []AfterExecuteT{
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
			},
			Before: []BeforeExecute{
				func(request *http.Request) error {
					return nil
				},
				func(request *http.Request) error {
					return nil
				},
				func(request *http.Request) error {
					return nil
				},
				func(request *http.Request) error {
					return nil
				},
			},
			BeforeT: []BeforeExecuteT{
				func(t T, request *http.Request) error {
					return nil
				},
			},
		},
	}
	cuteObj := &cute{}

	qtBaseProps := &HTTPTestMaker{
		httpClient:    &http.Client{},
		jsonMarshaler: &jsonMarshaler{},
		middleware: &Middleware{
			After: []AfterExecute{
				func(*http.Response, []error) error {
					return nil
				},
				func(*http.Response, []error) error {
					return nil
				},
			},
			AfterT: []AfterExecuteT{func(T, *http.Response, []error) error { return nil }},
			Before: []BeforeExecute{
				func(*http.Request) error {
					return nil
				},
				func(*http.Request) error {
					return nil
				},
			},
			BeforeT: []BeforeExecuteT{
				func(T, *http.Request) error { return nil },
				func(t T, request *http.Request) error {
					return nil
				},
			},
		},
	}
	cuteObj.baseProps = qtBaseProps

	cuteObj.fillBaseProps(testObj)

	require.Equal(t, qtBaseProps.httpClient, testObj.httpClient)
	require.Equal(t, qtBaseProps.jsonMarshaler, testObj.jsonMarshaler)
	require.Len(t, testObj.Middleware.After, len(qtBaseProps.middleware.After)+2)
	require.Len(t, testObj.Middleware.AfterT, len(qtBaseProps.middleware.AfterT)+3)
	require.Len(t, testObj.Middleware.Before, len(qtBaseProps.middleware.Before)+4)
	require.Len(t, testObj.Middleware.BeforeT, len(qtBaseProps.middleware.BeforeT)+1)
}
