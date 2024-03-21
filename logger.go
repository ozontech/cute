package cute

import (
	"fmt"
)

type tlogger interface {
	Name() string
	Logf(format string, args ...any)
	Errorf(format string, args ...interface{})
}

func (it *Test) Info(t tlogger, format string, args ...interface{}) {
	it.logf(t, "INFO", format, args...)
}

func (it *Test) Error(t tlogger, format string, args ...interface{}) {
	it.errorf(t, "ERROR", format, args...)
}

func (it *Test) Debug(t tlogger, format string, args ...interface{}) {
	it.logf(t, "DEBUG", format, args...)
}

func (it *Test) logf(t tlogger, level, format string, args ...interface{}) {
	name := it.Name

	if it.Name == "" {
		name = t.Name()
	}

	t.Logf("[%s][%s] %v\n", name, level, fmt.Sprintf(format, args...))
}

func (it *Test) errorf(t tlogger, level, format string, args ...interface{}) {
	name := it.Name

	if it.Name == "" {
		name = t.Name()
	}

	t.Errorf("[%s][%s] %v\n", name, level, fmt.Sprintf(format, args...))
}
