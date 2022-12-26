package cute

import (
	"net/url"
)

type RequestBuilder func(o *requestOptions)

// File is struct for upload file in form field
// If you set Path, file will read from file system
// If you set Name and Body, file will set from this fields
type File struct {
	Path string
	Name string
	Body []byte
}

type requestOptions struct {
	method      string
	url         *url.URL
	uri         string
	headers     map[string][]string
	body        []byte
	bodyMarshal interface{}
	fileForms   map[string]*File
	forms       map[string][]byte
}

func newRequestOptions() *requestOptions {
	return &requestOptions{
		headers:   make(map[string][]string),
		fileForms: make(map[string]*File),
		forms:     make(map[string][]byte),
	}
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

// WithHeadersKV is a function for set headers in request
func WithHeadersKV(name string, value string) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.headers[name] = []string{value}
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

// WithFileFormKV is a function for set file form in request
func WithFileFormKV(name string, file *File) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.fileForms[name] = file
	}
}

// WithFileForm is a function for set file form in request
func WithFileForm(fileForms map[string]*File) func(o *requestOptions) {
	return func(o *requestOptions) {
		for name, file := range fileForms {
			o.fileForms[name] = file
		}
	}
}

// WithFormKV is a function for set body in form request
func WithFormKV(name string, body []byte) func(o *requestOptions) {
	return func(o *requestOptions) {
		o.forms[name] = body
	}
}

// WithForm is a function for set body in form request
func WithForm(forms map[string][]byte) func(o *requestOptions) {
	return func(o *requestOptions) {
		for name, body := range forms {
			o.forms[name] = body
		}
	}
}
