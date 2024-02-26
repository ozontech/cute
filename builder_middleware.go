package cute

func (qt *cute) StepName(name string) MiddlewareRequest {
	qt.tests[qt.countTests].AllureStep.Name = name

	return qt
}

func (qt *cute) BeforeExecute(fs ...BeforeExecute) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.Before = append(qt.tests[qt.countTests].Middleware.Before, fs...)

	return qt
}

func (qt *cute) BeforeExecuteT(fs ...BeforeExecuteT) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.BeforeT = append(qt.tests[qt.countTests].Middleware.BeforeT, fs...)

	return qt
}

func (qt *cute) After(fs ...AfterExecute) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Middleware.After = append(qt.tests[qt.countTests].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterT(fs ...AfterExecuteT) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Middleware.AfterT = append(qt.tests[qt.countTests].Middleware.AfterT, fs...)

	return qt
}

func (qt *cute) AfterExecute(fs ...AfterExecute) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.After = append(qt.tests[qt.countTests].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterExecuteT(fs ...AfterExecuteT) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.AfterT = append(qt.tests[qt.countTests].Middleware.AfterT, fs...)

	return qt
}

func (qt *cute) AfterTestExecute(fs ...AfterExecute) NextTestBuilder {
	previousTest := 0
	if qt.countTests != 0 {
		previousTest = qt.countTests - 1
	}

	qt.tests[previousTest].Middleware.After = append(qt.tests[previousTest].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterTestExecuteT(fs ...AfterExecuteT) NextTestBuilder {
	previousTest := 0
	if qt.countTests != 0 {
		previousTest = qt.countTests - 1
	}

	qt.tests[previousTest].Middleware.AfterT = append(qt.tests[previousTest].Middleware.AfterT, fs...)

	return qt
}
