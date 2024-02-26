package cute

import "net/http"

func (qt *cute) CreateTableTest() MiddlewareTable {
	qt.isTableTest = true

	return qt
}

func (qt *cute) PutNewTest(name string, r *http.Request, expect *Expect) TableTest {
	// Validate, that first step is empty
	if qt.countTests == 0 {
		if qt.tests[0].Request.Base == nil &&
			len(qt.tests[0].Request.Builders) == 0 {
			qt.tests[0].Expect = expect
			qt.tests[0].Name = name
			qt.tests[0].Request.Base = r

			return qt
		}
	}

	newTest := createDefaultTest(qt.baseProps)
	newTest.Expect = expect
	newTest.Name = name
	newTest.Request.Base = r
	qt.tests = append(qt.tests, newTest)
	qt.countTests++ // async?

	return qt
}

func (qt *cute) PutTests(tests ...*Test) TableTest {
	for _, test := range tests {
		// Fill common fields
		qt.fillBaseProps(test)

		// Validate, that first step is empty
		if qt.countTests == 0 {
			if qt.tests[0].Request.Base == nil &&
				len(qt.tests[0].Request.Builders) == 0 {
				qt.tests[0] = test

				continue
			}
		}

		qt.tests = append(qt.tests, test)
		qt.countTests++
	}

	return qt
}

func (qt *cute) fillBaseProps(t *Test) {
	if qt.baseProps == nil {
		return
	}

	if qt.baseProps.httpClient != nil {
		t.httpClient = qt.baseProps.httpClient
	}

	if qt.baseProps.jsonMarshaler != nil {
		t.jsonMarshaler = qt.baseProps.jsonMarshaler
	}

	if t.Middleware == nil {
		t.Middleware = createMiddlewareFromTemplate(qt.baseProps.middleware)
	} else {
		t.Middleware.After = append(t.Middleware.After, qt.baseProps.middleware.After...)
		t.Middleware.AfterT = append(t.Middleware.AfterT, qt.baseProps.middleware.AfterT...)
		t.Middleware.Before = append(t.Middleware.Before, qt.baseProps.middleware.Before...)
		t.Middleware.BeforeT = append(t.Middleware.BeforeT, qt.baseProps.middleware.BeforeT...)
	}
}

func (qt *cute) NextTest() NextTestBuilder {
	qt.countTests++ // async?

	qt.tests = append(qt.tests, createDefaultTest(qt.baseProps))

	return qt
}
