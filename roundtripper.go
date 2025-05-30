package cute

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"moul.io/http2curl/v2"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/ozontech/cute/internal/utils"
)

func (it *Test) makeRequest(t internalT, req *http.Request) (*http.Response, []error) {
	var (
		delay       = defaultDelayRepeat
		countRepeat = 1

		resp  *http.Response
		err   error
		scope = make([]error, 0)
	)

	if it.Request.Retry.Delay != 0 {
		delay = it.Request.Retry.Delay
	}

	if it.Request.Retry.Count != 0 {
		countRepeat = it.Request.Retry.Count
	}

	for i := 1; i <= countRepeat; i++ {
		it.executeWithStep(t, it.createTitle(i, countRepeat, req), func(t T) []error {
			resp, err = it.doRequest(t, req)
			if err != nil {
				if it.Request.Retry.Broken {
					err = wrapBrokenError(err)
				}

				if it.Request.Retry.Optional {
					err = wrapOptionalError(err)
				}

				return []error{err}
			}

			return nil
		})

		if err == nil {
			break
		}

		scope = append(scope, err)

		if i != countRepeat {
			time.Sleep(delay)
		}
	}

	return resp, scope
}

func (it *Test) doRequest(t T, baseReq *http.Request) (*http.Response, error) {
	// copy request, because body can be read once
	req, err := copyRequest(baseReq.Context(), baseReq)
	if err != nil {
		return nil, cuteErrors.NewCuteError("[Internal] Could not copy request", err)
	}

	resp, httpErr := it.httpClient.Do(req)

	// if the timeout is triggered, we properly log the timeout error on allure and in traces
	if errors.Is(httpErr, context.DeadlineExceeded) {
		// Add information (method, host, curl) about request to Allure step
		// should be after httpClient.Do and from resp.Request, because in roundTripper request may be changed
		if addErr := it.addInformationRequest(t, req); addErr != nil {
			// Ignore err return, because it's connected with test logic
			it.Error(t, "Could not log information about request. error %v", addErr)
		}

		return nil, cuteErrors.NewEmptyAssertError(
			"Request timeout",
			fmt.Sprintf("expected request to be completed in %v, but was not", it.Expect.ExecuteTime))
	}

	// http client has case wheh it return response and error in one time
	// we have to check this case
	if resp == nil {
		if httpErr != nil {
			return nil, cuteErrors.NewCuteError("[HTTP] Could not do request", httpErr)
		}

		// if response is nil, we can't get information about request and response
		return nil, cuteErrors.NewCuteError("[HTTP] Response is nil", httpErr)
	}

	// BAD CODE. Need to copy body, because we can't read body again from resp.Request.Body. Problem is io.Reader
	resp.Request.Body, baseReq.Body, err = utils.DrainBody(baseReq.Body)
	if err != nil {
		it.Error(t, "Could not drain body from baseReq.Body. error %v", err)
		// Ignore err return, because it's connected with test logic
	}

	// Add information (method, host, curl) about request to Allure step
	// should be after httpClient.Do and from resp.Request, because in roundTripper request may be changed
	if addErr := it.addInformationRequest(t, resp.Request); addErr != nil {
		it.Error(t, "Could not log information about request. error %v", addErr)
		// Ignore err return, because it's connected with test logic
	}

	if httpErr != nil {
		return nil, cuteErrors.NewCuteError("[HTTP] Could not do request", httpErr)
	}

	// Add information (code, body, headers) about response to Allure step
	if addErr := it.addInformationResponse(t, resp); addErr != nil {
		// Ignore err return, because it's connected with test logic
		it.Error(t, "Could not log information about response. error %v", addErr)
	}

	if validErr := it.validateResponseCode(resp); validErr != nil {
		return resp, validErr
	}

	return resp, nil
}

func (it *Test) validateResponseCode(resp *http.Response) error {
	if it.Expect.Code != 0 && it.Expect.Code != resp.StatusCode {
		return cuteErrors.NewAssertError(
			"Assert response code",
			fmt.Sprintf("Response code expect %v, but was %v", it.Expect.Code, resp.StatusCode),
			resp.StatusCode,
			it.Expect.Code)
	}

	return nil
}

func (it *Test) addInformationRequest(t T, req *http.Request) error {
	var (
		saveBody io.ReadCloser
		err      error
	)

	if it.RequestSanitizer != nil {
		it.RequestSanitizer(req)
	}

	it.lastRequestURL = req.URL.String()

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return err
	}

	if c := curl.String(); len(c) <= 2048 {
		it.Info(t, "[Request] "+c)
	} else {
		it.Info(t, "[Request] Do request")
	}

	// Do not change to JSONMarshaler
	// In this case we can keep default for keep JSON, independence from JSONMarshaler
	headers, err := utils.ToJSON(req.Header)
	if err != nil {
		return err
	}

	t.WithParameters(
		allure.NewParameters(
			"method", req.Method,
			"host", req.Host,
			"headers", headers,
			"curl", curl.String(),
		)...,
	)

	if req.Body != nil {
		saveBody, req.Body, err = utils.DrainBody(req.Body)
		if err != nil {
			return err
		}

		body, err := utils.GetBody(saveBody)
		if err != nil {
			return err
		}

		if len(body) != 0 {
			t.WithNewParameters("body", string(body))
		}
	}

	return nil
}

func copyRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	var (
		err error

		clone = req.Clone(ctx)
	)

	req.Body, clone.Body, err = utils.DrainBody(req.Body)
	if err != nil {
		return nil, err
	}

	return clone, nil
}

func (it *Test) addInformationResponse(t T, response *http.Response) error {
	var (
		saveBody io.ReadCloser
		err      error
	)

	if it.ResponseSanitizer != nil {
		it.ResponseSanitizer(response)
	}

	headers, _ := utils.ToJSON(response.Header)
	if headers != "" {
		t.WithNewParameters("response_headers", headers)
	}

	t.WithNewParameters("response_code", fmt.Sprint(response.StatusCode))
	it.Info(t, "[Response] Status: "+response.Status)

	if response.Body == nil {
		return nil
	}

	saveBody, response.Body, err = utils.DrainBody(response.Body)
	// if could not get body from response, no add to allure
	if err != nil {
		return err
	}

	body, err := utils.GetBody(saveBody)
	// if could not get body from response, no add to allure
	if err != nil {
		return err
	}

	// if body is empty - skip
	if len(body) == 0 {
		return nil
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

	t.WithAttachments(allure.NewAttachment("response", responseType, body))

	return nil
}

func (it *Test) createTitle(try, countRepeat int, req *http.Request) string {
	toProcess := req

	// We have to execute sanitizer hook because
	// we need to log it and it can contain sensitive data
	if it.RequestSanitizer != nil {
		clone, err := copyRequest(req.Context(), req)

		// ignore error, because we want to log request
		// and it does not matter if we can copy request
		if err == nil {
			it.RequestSanitizer(clone)

			toProcess = clone
		}
	}

	title := toProcess.Method + " " + toProcess.URL.String()

	if countRepeat == 1 {
		return title
	}

	return fmt.Sprintf("[%v/%v] %v", try, countRepeat, title)
}
