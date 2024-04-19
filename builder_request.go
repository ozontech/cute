package cute

import (
	"net/http"
	"time"
)

// RequestRepeat is a function for set options in request
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
func (qt *cute) RequestRepeat(count int) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Count = count

	return qt
}

// RequestRepeatDelay set delay for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// Default delay is 1 second.
func (qt *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Delay = delay

	return qt
}

// RequestRepeatPolitic set politic for request repeat.
// if response.Code != Expect.Code, than request will repeat Count counts with Delay delay.
// if Optional is true and request is failed, than test step allure will be option, and t.Fail() will not execute.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
func (qt *cute) RequestRepeatPolitic(politic *RequestRepeatPolitic) RequestHTTPBuilder {
	if politic == nil {
		panic("politic is nil in RequestRepeatPolitic")
	}

	qt.tests[qt.countTests].Request.Repeat = politic

	return qt
}

// RequestRepeatOption set option politic for request repeat.
// if Optional is true and request is failed, than test step allure will be option, and t.Fail() will not execute.
func (qt *cute) RequestRepeatOptional(option bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Optional = option

	return qt
}

// RequestRepeatBroken set broken politic for request repeat.
// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will not execute.
func (qt *cute) RequestRepeatBroken(broken bool) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Broken = broken

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
