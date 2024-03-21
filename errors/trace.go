package errors

// NewErrorWithTrace ...
func NewErrorWithTrace(err, trace string) error {
	return &CuteError{
		Trace:   trace,
		Message: err,
	}
}

// WrapErrorWithTrace ...
func WrapErrorWithTrace(err error, trace string) error {
	return &CuteError{
		Trace: trace,
		Err:   err,
	}
}
