package cute

import (
	"net/http"
	"time"
)

func (qt *cute) RequestRepeat(count int) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Count = count

	return qt
}

func (qt *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Delay = delay

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
