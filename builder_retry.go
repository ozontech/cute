package cute

import "time"

// Retry is a function for configure test repeat
// if response.Code != Expect.Code or any of asserts are failed/broken than test will repeat counts with delay.
// Default delay is 1 second.
func (qt *cute) Retry(count int) MiddlewareRequest {
	if count < 1 {
		panic("count must be greater than 0")
	}

	qt.tests[qt.countTests].Retry.MaxAttempts = count

	return qt
}

// RetryDelay set delay for test repeat.
// if response.Code != Expect.Code or any of asserts are failed/broken than test will repeat counts with delay.
// Default delay is 1 second.
func (qt *cute) RetryDelay(delay time.Duration) MiddlewareRequest {
	if delay < 0 {
		panic("delay must be greater than or equal to 0")
	}

	qt.tests[qt.countTests].Retry.Delay = delay

	return qt
}
