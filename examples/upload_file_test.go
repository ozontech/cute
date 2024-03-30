//go:build example_upload_file
// +build example_upload_file

package examples

import (
	"context"
	"net/http"
	"testing"

	"github.com/ozontech/cute"
)

func TestUploadFile(t *testing.T) {
	cute.NewTestBuilder().
		Title("Upload file").
		Create().
		RequestBuilder(
			cute.WithURI("http://localhost:7000/v1/banner"),
			cute.WithMethod("POST"),
			cute.WithFormKV("body", []byte("{\"name\": \"Vasya\"}")), // Fill the form with the body
			cute.WithFileFormKV("image", &cute.File{ // Fill the form with the file
				Path: "/vasya/thebestmypicture.png",
			}),
		).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}
