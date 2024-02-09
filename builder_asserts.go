package cute

func (qt *cute) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, asserts...)

	return qt
}

func (qt *cute) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, optionalAssertBody(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, brokenAssertBody(assert))
	}

	return qt
}

func (qt *cute) RequireBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, requireAssertBody(assert))
	}

	return qt
}

func (qt *cute) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, asserts...)

	return qt
}

func (qt *cute) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, optionalAssertHeaders(assert))
	}

	return qt
}

func (qt *cute) RequireHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, requireAssertHeaders(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, brokenAssertHeaders(assert))
	}

	return qt
}

func (qt *cute) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, asserts...)

	return qt
}

func (qt *cute) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, optionalAssertResponse(assert))
	}

	return qt
}

func (qt *cute) RequireResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, requireAssertResponse(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, brokenAssertResponse(assert))
	}

	return qt
}

func (qt *cute) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, optionalAssertBodyT(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, brokenAssertBodyT(assert))
	}

	return qt
}

func (qt *cute) RequireBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, requireAssertBodyT(assert))
	}

	return qt
}

func (qt *cute) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, optionalAssertHeadersT(assert))
	}

	return qt
}

func (qt *cute) RequireHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, requireAssertHeadersT(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, brokenAssertHeadersT(assert))
	}

	return qt
}

func (qt *cute) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, optionalAssertResponseT(assert))
	}

	return qt
}

func (qt *cute) BrokenAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, brokenAssertResponseT(assert))
	}

	return qt
}

func (qt *cute) RequireResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, requireAssertResponseT(assert))
	}

	return qt
}
