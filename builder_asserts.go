package cute

import "time"

func (qt *cute) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, assertBodyWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, assertBodyWithTrace(optionalAssertBody(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, assertBodyWithTrace(brokenAssertBody(assert), trace))
	}

	return qt
}

func (qt *cute) RequireBody(asserts ...AssertBody) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, assertBodyWithTrace(requireAssertBody(assert), trace))
	}

	return qt
}

func (qt *cute) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, assertHeadersWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, assertHeadersWithTrace(optionalAssertHeaders(assert), trace))
	}

	return qt
}

func (qt *cute) RequireHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, assertHeadersWithTrace(requireAssertHeaders(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, assertHeadersWithTrace(brokenAssertHeaders(assert), trace))
	}

	return qt
}

func (qt *cute) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, assertResponseWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, assertResponseWithTrace(optionalAssertResponse(assert), trace))
	}

	return qt
}

func (qt *cute) RequireResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, assertResponseWithTrace(requireAssertResponse(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, assertResponseWithTrace(brokenAssertResponse(assert), trace))
	}

	return qt
}

func (qt *cute) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, assertBodyTWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, assertBodyTWithTrace(optionalAssertBodyT(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, assertBodyTWithTrace(brokenAssertBodyT(assert), trace))
	}

	return qt
}

func (qt *cute) RequireBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, assertBodyTWithTrace(requireAssertBodyT(assert), trace))
	}

	return qt
}

func (qt *cute) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, assertHeadersTWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, assertHeadersTWithTrace(optionalAssertHeadersT(assert), trace))
	}

	return qt
}

func (qt *cute) RequireHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, assertHeadersTWithTrace(requireAssertHeadersT(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, assertHeadersTWithTrace(brokenAssertHeadersT(assert), trace))
	}

	return qt
}

func (qt *cute) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, assertResponseTWithTrace(assert, trace))
	}

	return qt
}

func (qt *cute) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, assertResponseTWithTrace(optionalAssertResponseT(assert), trace))
	}

	return qt
}

func (qt *cute) BrokenAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, assertResponseTWithTrace(brokenAssertResponseT(assert), trace))
	}

	return qt
}

func (qt *cute) RequireResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	trace := getTrace()

	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, assertResponseTWithTrace(requireAssertResponseT(assert), trace))
	}

	return qt
}

func (qt *cute) ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.ExecuteTime = t

	return qt
}

func (qt *cute) ExpectStatus(code int) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.Code = code

	return qt
}

func (qt *cute) ExpectJSONSchemaString(schema string) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.String = schema

	return qt
}

func (qt *cute) ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.Byte = schema

	return qt
}

func (qt *cute) ExpectJSONSchemaFile(filePath string) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.File = filePath

	return qt
}

func (qt *cute) HideBody() ExpectHTTPBuilder {
	qt.tests[qt.countTests].HideBody = true

	return qt
}

func (qt *cute) HideHeaders() ExpectHTTPBuilder {
	qt.tests[qt.countTests].HideHeaders = true

	return qt
}

func (qt *cute) HideResponse() ExpectHTTPBuilder {
	qt.tests[qt.countTests].HideResponse = true

	return qt
}

func (qt *cute) HideResponseHeaders() ExpectHTTPBuilder {
	qt.tests[qt.countTests].HideResponseHeaders = true

	return qt
}
