package errors

// NewErrorWithTrace is a function for create error with trace
func NewErrorWithTrace(err, trace string) error {
	return &CuteError{
		Trace:   trace,
		Message: err,
	}
}

// WrapErrorWithTrace is a function for wrap error with trace
func WrapErrorWithTrace(err error, trace string) error {
	return &CuteError{
		Trace: trace,
		Err:   err,
	}
}
