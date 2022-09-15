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
	"testing"
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

// Test is a main struct of test.
// You may field Request and Expect for create simple test
type Test struct {
	httpClient *http.Client

	Name string

	AllureStep *AllureStep
	Middleware *Middleware
	Request    *Request
	Expect     *Expect
}

// Request is struct with HTTP request.
// You may use your *http.Request or create new with help Builders
type Request struct {
	Base     *http.Request
	Builders []RequestBuilder
	Repeat   *RequestRepeatPolitic
}

// RequestRepeatPolitic is struct for repeat politic
// If response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
type RequestRepeatPolitic struct {
	Count int
	Delay time.Duration
}

// Middleware is struct for execute something before or after test
type Middleware struct {
	After   []AfterExecute
	AfterT  []AfterExecuteT
	Before  []BeforeExecute
	BeforeT []BeforeExecuteT
}

// AllureStep is struct with test name
type AllureStep struct {
	Name string
}

// Expect is structs with validate politics for response
type Expect struct {
	ExecuteTime time.Duration

	Code       int
	JSONSchema *ExpectJSONSchema

	AssertBody     []AssertBody
	AssertHeaders  []AssertHeaders
	AssertResponse []AssertResponse

	AssertBodyT     []AssertBodyT
	AssertHeadersT  []AssertHeadersT
	AssertResponseT []AssertResponseT
}

// ExpectJSONSchema is structs with JSON politics for response
type ExpectJSONSchema struct {
	String string
	Byte   []byte
	File   string
}

func (it *Test) Execute(ctx context.Context, t testing.TB) ResultsHTTPBuilder {
	var (
		internalT allureProvider
		res       ResultsHTTPBuilder
	)

	tOriginal, ok := t.(*testing.T)
	if ok {
		internalT = createAllureT(tOriginal)
	}

	allureT, ok := t.(provider.T)
	if ok {
		internalT = allureT
	}

	internalT.Run(it.Name, func(inT provider.T) {
		res = it.execute(ctx, inT)
	})

	return res
}

func (it *Test) initEmptyFields() {
	it.httpClient = http.DefaultClient

	if it.AllureStep == nil {
		it.AllureStep = new(AllureStep)
	}
	if it.Middleware == nil {
		it.Middleware = new(Middleware)
	}
	if it.Expect == nil {
		it.Expect = new(Expect)
	}
	if it.Request == nil {
		it.Request = new(Request)
	}
	if it.Request.Repeat == nil {
		it.Request.Repeat = new(RequestRepeatPolitic)
	}
	if it.Expect.JSONSchema == nil {
		it.Expect.JSONSchema = new(ExpectJSONSchema)
	}
}

func (it *Test) execute(ctx context.Context, allureProvider allureProvider) ResultsHTTPBuilder {
	var (
		resp *http.Response
		errs []error
		name = allureProvider.Name() + "_" + it.Name
	)

	// Set empty fields in test
	it.initEmptyFields()

	if it.AllureStep.Name != "" {
		// Test with step
		name = it.AllureStep.Name

		// Execute Test
		allureProvider.Logf("Start step %v", name)
		resp, errs = it.startTestWithStep(ctx, allureProvider)
		allureProvider.Logf("Finish step %v", name)
	} else {
		// Execute Test
		// Test without step
		resp, errs = it.startTest(ctx, allureProvider)
	}

	processTestErrors(allureProvider, errs)

	return newTestResult(name, resp, errs)
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

func (it *Test) startTestWithStep(ctx context.Context, t internalT) (*http.Response, []error) {
	var (
		resp *http.Response
		errs []error
	)

	t.WithNewStep(it.AllureStep.Name, func(stepCtx provider.StepCtx) {
		resp, errs = it.startTest(ctx, stepCtx)

		if len(errs) != 0 {
			stepCtx.Fail()
		}
	})

	return resp, errs
}

func (it *Test) startTest(ctx context.Context, t internalT) (*http.Response, []error) {
	var (
		resp *http.Response
		err  error
	)

	// CreateWithStep execute timer
	if it.Expect.ExecuteTime == 0 {
		it.Expect.ExecuteTime = defaultExecuteTestTime
	}

	ctx, cancel := context.WithTimeout(ctx, it.Expect.ExecuteTime)
	defer cancel()

	// CreateWithStep request
	req, err := it.createRequest(ctx)
	if err != nil {
		return resp, []error{err}
	}

	// Execute Before
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

	// Execute After
	afterTestErrs := it.afterTest(t, resp, errs)

	// Return results
	errs = append(errs, afterTestErrs...)
	if len(errs) > 0 {
		return resp, errs
	}

	return resp, nil
}

func (it *Test) afterTest(t internalT, resp *http.Response, errs []error) []error {
	if len(it.Middleware.After) == 0 && len(it.Middleware.AfterT) == 0 {
		return nil
	}

	return executeWithStep(t, "After", func(t T) []error {
		scope := make([]error, 0)

		for _, execute := range it.Middleware.After {
			if err := execute(resp, errs); err != nil {
				scope = append(scope, err)
			}
		}

		for _, executeSuite := range it.Middleware.AfterT {
			if err := executeSuite(t, resp, errs); err != nil {
				scope = append(scope, err)
			}
		}

		return scope
	})
}

func (it *Test) beforeTest(t internalT, req *http.Request) []error {
	if len(it.Middleware.Before) == 0 && len(it.Middleware.BeforeT) == 0 {
		return nil
	}

	return executeWithStep(t, "Before", func(t T) []error {
		scope := make([]error, 0)

		for _, execute := range it.Middleware.Before {
			if err := execute(req); err != nil {
				scope = append(scope, err)
			}
		}

		for _, executeT := range it.Middleware.BeforeT {
			if err := executeT(t, req); err != nil {
				scope = append(scope, err)
			}
		}

		return scope
	})
}

func (it *Test) createRequest(ctx context.Context) (*http.Request, error) {
	var (
		req = it.Request.Base
		err error
	)

	// CreateWithStep Request
	if req == nil {
		o := new(requestOptions)

		for _, builder := range it.Request.Builders {
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

	// Validate Request
	if err := it.validateRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (it *Test) validateRequest(req *http.Request) error {
	if req.URL == nil {
		return errorRequestURLEmpty
	}

	if req.Method == "" {
		return errorRequestMethodEmpty
	}

	return nil
}

func (it *Test) validateResponse(t internalT, resp *http.Response) []error {
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
