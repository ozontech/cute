package utils

import (
	"bytes"
	"io"
	"net/http"
)

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

func DrainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
