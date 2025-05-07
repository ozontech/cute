package cute

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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

// RequestSanitizerHook is a function used to modify the request URL
// before it is logged or attached to test reports (e.g., for hiding secrets).
type RequestSanitizerHook func(req *http.Request)

// ResponseSanitizerHook is a function used to modify the response
// before it is logged or attached to test reports (e.g., for hiding secrets).
type ResponseSanitizerHook func(resp *http.Response)

// Test is a main struct of test.
// You may field Request and Expect for create simple test
// Parallel can be used to control the parallelism of a Test
type Test struct {
	httpClient     *http.Client
	jsonMarshaler  JSONMarshaler
	lastRequestURL string

	Name     string
	Parallel bool
	Retry    *Retry

	AllureStep *AllureStep
	Middleware *Middleware
	Request    *Request
	Expect     *Expect

	RequestSanitizer  RequestSanitizerHook
	ResponseSanitizer ResponseSanitizerHook
}

// Retry is a struct to control the retry of a whole single test (not only the request)
// The test will be retried up to MaxAttempts times
// The retries will only be executed if the test is having errors
// If the test is successful at any iteration between attempt 1 and MaxAttempts, the loop will break and return the result as successful
// The status of the test (success or fail) will be based on either the first attempt that is successful, or, if no attempt
// is successful, it will be based on the latest execution
// Delay is the number of seconds to wait before attempting to run the test again. It will only wait if Delay is set.
type Retry struct {
	currentCount int
	MaxAttempts  int
	Delay        time.Duration
}

// Request is struct with HTTP request.
// You may use your *http.Request or create new with help Builders
type Request struct {
	Base     *http.Request
	Builders []RequestBuilder
	Retry    *RequestRetryPolitic
}

// RequestRetryPolitic is struct for repeat politic
// if Optional is true and request is failed, than test step allure will be skip, and t.Fail() will not execute.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
// If Optional and Broken is false, than test step will be failed, and t.Fail() will execute.
// If response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
type RequestRetryPolitic struct {
	Count    int
	Delay    time.Duration
	Optional bool
	Broken   bool
}

// RequestRepeatPolitic is struct for repeat politic
// Deprecated: use RequestRetryPolitic
type RequestRepeatPolitic struct {
	Count    int
	Delay    time.Duration
	Optional bool
	Broken   bool
}

// Middleware is struct for executeInsideAllure something before or after test
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

// Execute is common method for run test from builder
func (it *Test) Execute(ctx context.Context, t tProvider) ResultsHTTPBuilder {
	var (
		internalT allureProvider
		res       ResultsHTTPBuilder
	)

	if t == nil {
		panic("could not start test without testing.T")
	}

	stepCtx, isStepCtx := t.(provider.StepCtx)
	if isStepCtx {
		return it.executeInsideStep(ctx, stepCtx)
	}

	tOriginal, ok := t.(*testing.T)
	if ok {
		tOriginal.Helper()
		internalT = createAllureT(tOriginal)
	}

	allureT, ok := t.(provider.T)
	if ok {
		internalT = allureT
	}

	internalT.Run(it.Name, func(inT provider.T) {
		if it.Parallel {
			inT.Parallel()
		}
		res = it.executeInsideAllure(ctx, inT)
	})

	return res
}

func (it *Test) clearFields() {
	it.AllureStep = new(AllureStep)
	it.Middleware = new(Middleware)
	it.Expect = new(Expect)
	it.Request = new(Request)
	it.Request.Retry = new(RequestRetryPolitic)
	it.Expect.JSONSchema = new(ExpectJSONSchema)
}

func (it *Test) initEmptyFields() {
	if it.httpClient == nil {
		it.httpClient = http.DefaultClient
	}

	if it.jsonMarshaler == nil {
		it.jsonMarshaler = jsonMarshaler{}
	}

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

	if it.Request.Retry == nil {
		it.Request.Retry = new(RequestRetryPolitic)
	}

	if it.Expect.JSONSchema == nil {
		it.Expect.JSONSchema = new(ExpectJSONSchema)
	}

	if it.Retry == nil {
		it.Retry = &Retry{
			// we set the default value to 1, because we count the first attempt as 1
			MaxAttempts:  1,
			currentCount: 1,
		}
	}

	// We need to set the current count to 1 here, because we count the first attempt as 1
	it.Retry.currentCount = 1
}

// executeInsideStep is method for start test with provider.StepCtx
// It's test inside the step
func (it *Test) executeInsideStep(ctx context.Context, t internalT) ResultsHTTPBuilder {
	// Set empty fields in test
	it.initEmptyFields()

	// we don't want to defer the finish message, because it will be logged in processTestErrors
	it.Info(t, "Start test")

	return it.startRepeatableTest(ctx, t)
}

func (it *Test) executeInsideAllure(ctx context.Context, allureProvider allureProvider) ResultsHTTPBuilder {
	// Set empty fields in test
	it.initEmptyFields()

	// we don't want to defer the finish message, because it will be logged in processTestErrors
	it.Info(allureProvider, "Start test")

	if it.AllureStep.Name != "" {
		// Execute test inside step
		return it.startTestInsideStep(ctx, allureProvider)
	} else {
		// Execute Test
		return it.startRepeatableTest(ctx, allureProvider)
	}
}

// startRepeatableTest is method for start test with repeatable execution
func (it *Test) startRepeatableTest(ctx context.Context, t internalT) ResultsHTTPBuilder {
	var (
		resp        *http.Response
		errs        []error
		resultState ResultState
	)

	for ; it.Retry.currentCount <= it.Retry.MaxAttempts; it.Retry.currentCount++ {
		resp, errs = it.startTest(ctx, t)

		resultState = it.processTestErrors(t, errs)

		// we don't want to keep errors if we will retry test
		// we have to return to user only errors from last try

		// if the test is successful, we break the loop
		if resultState == ResultStateSuccess {
			break
		}

		// if we have a delay, we wait before the next attempt
		// and we only wait if we are not at the last attempt
		if it.Retry.currentCount != it.Retry.MaxAttempts && it.Retry.Delay != 0 {
			it.Info(t, "The test had errors, retrying...")
			time.Sleep(it.Retry.Delay)
		}
	}

	switch resultState {
	case ResultStateBroken:
		t.BrokenNow()
		it.Info(t, "Test broken")
	case ResultStateFail:
		t.Fail()
		it.Error(t, "Test failed")
	case resultStateFailNow:
		t.FailNow()
		it.Error(t, "Test failed")
	case ResultStateSuccess:
		it.Info(t, "Test finished successfully")
	}

	return newTestResult(it.Name, resp, resultState, errs)
}

func (it *Test) startTestInsideStep(ctx context.Context, t internalT) ResultsHTTPBuilder {
	var (
		result ResultsHTTPBuilder
	)

	t.WithNewStep(it.AllureStep.Name, func(stepCtx provider.StepCtx) {
		it.Info(t, "Start step %v", it.AllureStep.Name)
		defer it.Info(t, "Finish step %v", it.AllureStep.Name)

		result = it.startRepeatableTest(ctx, stepCtx)

		if result.GetResultState() == ResultStateFail {
			stepCtx.Fail()
		}
	})

	return result
}

// processTestErrors returns flag, which mean finish test or not.
// If test has only optional errors, than test will be success
// If test has broken errors, than test will be broken on allure
// If test has require errors, than test will be failed on allure
func (it *Test) processTestErrors(t internalT, errs []error) ResultState {
	if len(errs) == 0 {
		it.Info(t, "Test finished successfully")

		return ResultStateSuccess
	}

	var (
		countNotOptionalErrors = 0
		state                  ResultState
	)

	for _, err := range errs {
		message := fmt.Sprintf("error %v", err.Error())

		if tErr, ok := err.(cuteErrors.OptionalError); ok {
			if tErr.IsOptional() {
				it.Info(t, "[OPTIONAL ERROR] %v", err.Error())

				state = ResultStateSuccess

				continue
			}
		}

		if tErr, ok := err.(cuteErrors.BrokenError); ok {
			if tErr.IsBroken() {
				it.Error(t, "[BROKEN ERROR], error %v", err.Error())

				state = ResultStateBroken

				continue
			}
		}

		if tErr, ok := err.(cuteErrors.RequireError); ok {
			if tErr.IsRequire() {
				state = resultStateFailNow
			}
		}

		if tErr, ok := err.(cuteErrors.WithFields); ok {
			actual := tErr.GetFields()[cuteErrors.ActualField]
			expected := tErr.GetFields()[cuteErrors.ExpectedField]

			if actual != nil || expected != nil {
				message = fmt.Sprintf("%s\nActual %v\nExpected %v", message, actual, expected)
			}
		}

		it.Error(t, message)

		countNotOptionalErrors++
	}

	if countNotOptionalErrors != 0 {
		state = ResultStateFail

		it.Error(t, "Test finished with %v errors", countNotOptionalErrors)
	}

	return state
}

func (it *Test) startTest(ctx context.Context, t internalT) (*http.Response, []error) {
	var (
		resp *http.Response
		err  error
	)

	// CreateWithStep executeInsideAllure timer
	if it.Expect.ExecuteTime == 0 {
		it.Expect.ExecuteTime = defaultExecuteTestTime
	}

	ctx, cancel := context.WithTimeout(ctx, it.Expect.ExecuteTime)
	defer cancel()

	// CreateWithStep request
	req, err := it.createRequest(ctx)
	if err != nil {
		return nil, []error{err}
	}

	// Execute Before
	if errs := it.beforeTest(t, req); len(errs) > 0 {
		return nil, errs
	}

	it.Info(t, "Start make request")

	// Make request
	resp, errs := it.makeRequest(t, req)
	if len(errs) > 0 {
		return resp, errs
	}

	it.Info(t, "Finish make request")

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

	return it.executeWithStep(t, "After", func(t T) []error {
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

	return it.executeWithStep(t, "Before", func(t T) []error {
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

// createRequest builds the final *http.Request to be executed by the test.
// If the Test.RequestSanitizer hook is defined, it will be called after validation
// to allow safe modification of the request before logging or execution.
func (it *Test) createRequest(ctx context.Context) (*http.Request, error) {
	var (
		req = it.Request.Base
		err error
	)

	if req == nil {
		req, err = it.buildRequest(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Validate Request
	if err := it.validateRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

// buildRequest
// Priority for create body:
// 1. requestOptions.body <- low priority
// 2. requestOptions.bodyMarshal
// 3. requestOptions.forms and requestOptions.fileForms <- high priority.
func (it *Test) buildRequest(ctx context.Context) (*http.Request, error) {
	var (
		req *http.Request
		err error

		o = newRequestOptions()
	)

	// Set builder parameters
	for _, builder := range it.Request.Builders {
		builder(o)
	}

	reqURL := o.url

	if reqURL == nil {
		reqURL, err = url.Parse(o.uri)
		if err != nil {
			return nil, err
		}
	}

	// Set query parameters
	query := reqURL.Query()
	for key, values := range o.query {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	reqURL.RawQuery = query.Encode()

	// Set body
	body := o.body
	if o.bodyMarshal != nil {
		body, err = it.jsonMarshaler.Marshal(o.bodyMarshal)

		if err != nil {
			return nil, err
		}
	}

	// Set multipart
	if len(o.fileForms) != 0 || len(o.forms) != 0 {
		var (
			buffer = new(bytes.Buffer)
			mp     = multipart.NewWriter(buffer)
		)

		// set file forms
		for fieldName, file := range o.fileForms {
			err = createFormFile(mp, fieldName, file)
			if err != nil {
				return nil, err
			}
		}

		// set forms
		for fieldName, fieldBody := range o.forms {
			err = createFormField(mp, fieldName, fieldBody)
			if err != nil {
				return nil, err
			}
		}

		if err = mp.Close(); err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, o.method, reqURL.String(), buffer)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", mp.FormDataContentType())
	} else {
		req, err = http.NewRequestWithContext(ctx, o.method, reqURL.String(), io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			return nil, err
		}
	}

	// Set headers
	for nameHeader, valuesHeader := range o.headers {
		req.Header[nameHeader] = valuesHeader
	}

	return req, nil
}

func createFormFile(mp *multipart.Writer, fieldName string, file *File) error {
	var (
		data = file.Body
		name = file.Name
	)

	// read file, if path is not empty
	if len(file.Path) != 0 {
		f, err := os.Open(file.Path)
		if err != nil {
			return err
		}

		data, err = io.ReadAll(f)
		if err != nil {
			return err
		}

		name = f.Name()
	}

	field, err := mp.CreateFormFile(fieldName, name)
	if err != nil {
		return fmt.Errorf("error when creating %v file form field, %w", fieldName, err)
	}

	_, err = field.Write(data)
	if err != nil {
		return fmt.Errorf("error when writing %v file form field, %w", fieldName, err)
	}

	return nil
}

func createFormField(mp *multipart.Writer, fieldName string, body []byte) error {
	field, err := mp.CreateFormField(fieldName)
	if err != nil {
		return fmt.Errorf("error when creating %v form field, %w", fieldName, err)
	}

	_, err = field.Write(body)
	if err != nil {
		return fmt.Errorf("error when writing %v form field, %w", fieldName, err)
	}

	return nil
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
		return append(scope, fmt.Errorf("could not drain response body. error %w", err))
	}

	body, err := utils.GetBody(saveBody)
	if err != nil {
		return append(scope, fmt.Errorf("could not get response body. error %w", err))
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
