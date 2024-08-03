package cute

import (
	"net/http"
	"time"
)

// RequestRepeat is a function for set options in request
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
// Deprecated: use RequestRetry instead
func (qt *cute) RequestRepeat(count int) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Count = count

	return qt
}

// RequestRepeatDelay set delay for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
// Deprecated: use RequestRetryDelay instead
func (qt *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Delay = delay

	return qt
}

// RequestRepeatPolitic set politic for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
// Deprecated: use RequestRetryPolitic instead
func (qt *cute) RequestRepeatPolitic(politic *RequestRepeatPolitic) RequestHTTPBuilder {
	if politic == nil {
		panic("politic is nil in RequestRetryPolitic")
	}

	qt.tests[qt.countTests].Request.Retry = &RequestRetryPolitic{
		Count:    politic.Count,
		Delay:    politic.Delay,
		Optional: politic.Optional,
		Broken:   politic.Broken,
	}

	return qt
}

// RequestRepeatOptional set option politic for request repeat.
// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
// Deprecated: use RequestRetryOptional instead
func (qt *cute) RequestRepeatOptional(option bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Optional = option

	return qt
}

// RequestRepeatBroken set broken politic for request repeat.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
// Deprecated: use RequestRetryBroken instead
func (qt *cute) RequestRepeatBroken(broken bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Broken = broken

	return qt
}

// RequestRetry is a function for set options in request
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
func (qt *cute) RequestRetry(count int) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Count = count

	return qt
}

// RequestRetryDelay set delay for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
func (qt *cute) RequestRetryDelay(delay time.Duration) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Delay = delay

	return qt
}

// RequestRetryPolitic set politic for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
func (qt *cute) RequestRetryPolitic(politic *RequestRetryPolitic) RequestHTTPBuilder {
	if politic == nil {
		panic("politic is nil in RequestRetryPolitic")
	}

	qt.tests[qt.countTests].Request.Retry = politic

	return qt
}

// RequestRetryOptional set option politic for request repeat.
// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
func (qt *cute) RequestRetryOptional(option bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Optional = option

	return qt
}

// RequestRetryBroken set broken politic for request repeat.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
func (qt *cute) RequestRetryBroken(broken bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Retry.Broken = broken

	return qt
}

func (qt *cute) Request(r *http.Request) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Request.Base = r

	return qt
}

func (qt *cute) RequestBuilder(r ...RequestBuilder) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Request.Builders = append(qt.tests[qt.countTests].Request.Builders, r...)

	return qt
}
