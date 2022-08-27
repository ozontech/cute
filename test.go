package cute

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/ozontech/cute/internal/utils"
)

const (
	defaultExecuteTestTime = 10 * time.Second
	defaultDelayRepeat     = 1 * time.Second
)

var (
	errorRequestMethodEmpty = errors.New("request method must be not empty")
	errorRequestURLEmpty    = errors.New("url request must be not empty")
)

type test struct {
	httpClient *http.Client

	name string

	allureStep *allureStep
	middleware *middleware
	request    *request
	expect     *Expect
}

type request struct {
	base     *http.Request
	builders []requestBuilder
	repeat   *requestRepeatPolitic
}

type requestRepeatPolitic struct {
	count int
	delay time.Duration
}

type middleware struct {
	after   []AfterExecute
	afterT  []AfterExecuteT
	before  []BeforeExecute
	beforeT []BeforeExecuteT
}

type allureStep struct {
	name string
}

type Expect struct {
	ExecuteTime time.Duration

	Code             int
	JSONSchemaString string
	JSONSchemaByte   []byte
	JSONSchemaFile   string

	AssertBody     []AssertBody
	AssertHeaders  []AssertHeaders
	AssertResponse []AssertResponse

	AssertBodyT     []AssertBodyT
	AssertHeadersT  []AssertHeadersT
	AssertResponseT []AssertResponseT
}

func (it *test) execute(ctx context.Context, allureProvider allureProvider) []ResultsHTTPBuilder {
	var (
		res  = make([]ResultsHTTPBuilder, 0)
		resp *http.Response
		errs []error
		name = allureProvider.Name() + "_" + it.name
	)

	if it.allureStep.name != "" {
		// Test with step
		name = it.allureStep.name

		// Execute test
		allureProvider.Logf("Start step %v", name)
		resp, errs = it.startTestWithStep(ctx, allureProvider)
		allureProvider.Logf("Finish step %v", name)
	} else {
		// Execute test
		// Test without step
		resp, errs = it.startTest(ctx, allureProvider)
	}

	processTestErrors(allureProvider, errs)

	res = append(res, newTestResult(name, resp, errs))

	return res
}

func processTestErrors(t internalT, errs []error) {
	if len(errs) == 0 {
		return
	}

	resErrors := make([]error, 0)
	for _, err := range errs {
		if tErr, ok := err.(cuteErrors.OptionalError); ok {
			if tErr.IsOptional() {
				t.Logf("[OPTIONAL ERROR] %v", tErr.(error).Error())
				continue
			}
		}

		resErrors = append(resErrors, err)
	}

	if len(resErrors) == 0 {
		return
	}

	if len(resErrors) == 1 {
		t.Errorf("[ERROR] %v", resErrors[0])

		return
	}

	for _, err := range resErrors {
		t.Errorf("[ERROR] %v", err)
	}

	t.Errorf("Test finished with %v errors", len(resErrors))
}

func (it *test) startTestWithStep(ctx context.Context, t internalT) (*http.Response, []error) {
	var (
		resp *http.Response
		errs []error
	)

	t.WithNewStep(it.allureStep.name, func(stepCtx provider.StepCtx) {
		resp, errs = it.startTest(ctx, stepCtx)

		if len(errs) != 0 {
			stepCtx.Fail()
		}
	})

	return resp, errs
}

func (it *test) startTest(ctx context.Context, t internalT) (*http.Response, []error) {
	var (
		resp *http.Response
		err  error
	)

	// CreateWithStep execute timer
	if it.expect.ExecuteTime == 0 {
		it.expect.ExecuteTime = defaultExecuteTestTime
	}

	ctx, cancel := context.WithTimeout(ctx, it.expect.ExecuteTime)
	defer cancel()

	// CreateWithStep request
	req, err := it.createRequest(ctx)
	if err != nil {
		return resp, []error{err}
	}

	// Execute before
	if errs := it.beforeTest(t, req); len(errs) > 0 {
		return nil, errs
	}

	t.Logf("Start make request")

	// Make request
	resp, errs := it.makeRequest(t, req)
	if len(errs) > 0 {
		return resp, errs
	}

	t.Logf("Finish make request")

	// Validate response body
	errs = it.validateResponse(t, resp)

	// Execute after
	afterTestErrs := it.afterTest(t, resp, errs)

	// Return results
	errs = append(errs, afterTestErrs...)
	if len(errs) > 0 {
		return resp, errs
	}

	return resp, nil
}

func (it *test) afterTest(t internalT, resp *http.Response, errs []error) []error {
	if len(it.middleware.after) == 0 && len(it.middleware.afterT) == 0 {
		return nil
	}

	return executeWithStep(t, "After", func(t T) []error {
		scope := make([]error, 0)

		for _, execute := range it.middleware.after {
			if err := execute(resp, errs); err != nil {
				scope = append(scope, err)
			}
		}

		for _, executeSuite := range it.middleware.afterT {
			if err := executeSuite(t, resp, errs); err != nil {
				scope = append(scope, err)
			}
		}

		return scope
	}, false)
}

func (it *test) beforeTest(t internalT, req *http.Request) []error {
	if len(it.middleware.before) == 0 && len(it.middleware.beforeT) == 0 {
		return nil
	}

	return executeWithStep(t, "Before", func(t T) []error {
		scope := make([]error, 0)

		for _, execute := range it.middleware.before {
			if err := execute(req); err != nil {
				scope = append(scope, err)
			}
		}

		for _, executeT := range it.middleware.beforeT {
			if err := executeT(t, req); err != nil {
				scope = append(scope, err)
			}
		}

		return scope
	}, false)
}

func (it *test) createRequest(ctx context.Context) (*http.Request, error) {
	var (
		req = it.request.base
		err error
	)

	// CreateWithStep request
	if req == nil {
		o := new(requestOptions)

		for _, builder := range it.request.builders {
			builder(o)
		}

		url := o.uri
		if o.url != nil {
			url = o.url.String()
		}

		body := o.body
		if o.bodyMarshal != nil {
			body, err = json.Marshal(o.bodyMarshal) // TODO move marshaler to it struct

			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequestWithContext(ctx, o.method, url, ioutil.NopCloser(bytes.NewReader(body)))
		if err != nil {
			return nil, err
		}

		if len(o.headers) != 0 {
			req.Header = o.headers
		}
	}

	// Validate request
	if err := it.validateRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (it *test) validateRequest(req *http.Request) error {
	if req.URL == nil {
		return errorRequestURLEmpty
	}

	if req.Method == "" {
		return errorRequestMethodEmpty
	}

	return nil
}

func (it *test) validateResponse(t internalT, resp *http.Response) []error {
	var (
		err      error
		saveBody io.ReadCloser
		scope    = make([]error, 0)
	)

	// Execute asserts for headers
	if errs := it.assertHeaders(t, resp.Header); len(errs) > 0 {
		scope = append(scope, errs...)
	}

	// Prepare body for validate
	if resp.Body == nil {
		// todo create errors if body is empty, but assert is not empty
		return scope
	}

	saveBody, resp.Body, err = utils.DrainBody(resp.Body)
	if err != nil {
		return append(scope, fmt.Errorf("could not drain response body. error %v", err))
	}

	body, err := utils.GetBody(saveBody)
	if err != nil {
		return append(scope, fmt.Errorf("could not get response body. error %v", err))
	}

	// Execute asserts for body
	if errs := it.assertBody(t, body); len(errs) > 0 {
		// add assert
		scope = append(scope, errs...)
	}

	// Validate response by json schema
	if errs := it.validateJSONSchema(t, body); len(errs) > 0 {
		scope = append(scope, errs...)
	}

	// Execute asserts for response body
	if errs := it.assertResponse(t, resp); len(errs) > 0 {
		scope = append(scope, errs...)
	}

	return scope
}
