package cute

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/moul/http2curl"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/cute/internal/utils"
)

type allureRoundTripper struct {
	roundTripper http.RoundTripper
}

func newAllureRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &allureRoundTripper{
		roundTripper: next,
	}
}

func (a *allureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t, err := getProviderT(req.Context())
	if err != nil {
		return nil, err
	}

	step, err := prepareStep(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		t.Step(step)
	}()

	response, err := a.roundTripper.RoundTrip(req)

	// if err add allure.Attachment with text fail
	if err != nil {
		step.Status = allure.Failed
		step.WithAttachments(allure.NewAttachment("Error", allure.Text, []byte(err.Error())))

		return response, err
	}

	updateStepWithResponse(step, response)

	return response, nil
}

// PrepareStep returns step based on http request.
func prepareStep(req *http.Request) (*allure.Step, error) {
	var (
		saveBody io.ReadCloser
		err      error
	)

	step := allure.NewSimpleStep(req.Method + " " + req.URL.String())

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	headers, err := utils.ToJSON(req.Header)
	if err != nil {
		return nil, err
	}

	step.Parameters = allure.NewParameters(
		"method", req.Method,
		"host", req.Host,
		"headers", headers,
		"curl", curl.String(),
	)

	if req.Body != nil {
		saveBody, req.Body, err = utils.DrainBody(req.Body)
		if err != nil {
			return nil, err
		}

		body, err := utils.GetBody(saveBody)
		if err != nil {
			return nil, err
		}

		step.Parameters = append(step.Parameters, allure.NewParameter("body", string(body)))
	}

	return step, nil
}

// UpdateStepWithResponse returns step based on already created step and grpc response.
func updateStepWithResponse(step *allure.Step, response *http.Response) *allure.Step {
	var (
		saveBody io.ReadCloser
		err      error
	)

	step.Stop = time.Now().UnixNano() / int64(time.Millisecond)

	if response == nil {
		return step
	}

	headers, _ := utils.ToJSON(response.Header)
	if headers != "" {
		step.WithNewParameters("response_headers", headers)
	}

	step.WithNewParameters("response_code", fmt.Sprint(response.StatusCode))

	if response.Body == nil {
		return step
	}

	saveBody, response.Body, err = utils.DrainBody(response.Body)
	// if could not get body from response, no add to allure
	if err != nil {
		return step
	}
	body, err := utils.GetBody(saveBody)
	// if could not get body from response, no add to allure
	if err != nil {
		return step
	}

	responseType := allure.Text
	if _, ok := response.Header["Content-Type"]; ok {
		if len(response.Header["Content-Type"]) > 0 {
			if strings.Contains(response.Header["Content-Type"][0], "application/json") {
				responseType = allure.JSON
			} else {
				responseType = allure.MimeType(response.Header["Content-Type"][0])
			}
		}
	}
	if responseType == allure.JSON {
		body, _ = utils.PrettyJSON(body)
	}

	step.WithAttachments(allure.NewAttachment("response", responseType, body))
	return step
}
