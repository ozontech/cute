package cute

import (
	"net/url"
)

type requestBuilder func(o *requestOptions)

type requestOptions struct {
	method      string
	url         *url.URL
	uri         string
	headers     map[string][]string
	body        []byte
	bodyMarshal interface{}
}

// WithMethod is a function for set method (GET, POST ...) in request
func WithMethod(method string) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.method = method
	}
}

// WithURL is a function for set url in request
func WithURL(url *url.URL) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.url = url
	}
}

// WithURI is a function for set url in request
func WithURI(uri string) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.uri = uri
	}
}

// WithHeaders is a function for set headers in request
func WithHeaders(headers map[string][]string) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.headers = headers
	}
}

// WithBody is a function for set body in request
func WithBody(body []byte) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.body = body
	}
}

// WithMarshalBody is a function for marshal body and set body in request
func WithMarshalBody(body interface{}) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.bodyMarshal = body
	}
}
