package utils

import (
	"bytes"
	"io"
	"net/http"
)

// GetBody get body from IO
func GetBody(body io.ReadCloser) ([]byte, error) {
	var (
		err error
		buf = new(bytes.Buffer)
	)

	_, err = io.Copy(buf, body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DrainBody ...
func DrainBody(body io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if body == nil || body == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(body); err != nil {
		return nil, body, err
	}

	if err = body.Close(); err != nil {
		return nil, body, err
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
