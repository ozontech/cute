package cute

import (
	"fmt"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute/errors"
)

func (it *test) executeWithStep(t internalT, stepName string, execute func(t T) []error) []error {
	var (
		errs []error
	)

	t.WithNewStep(stepName, func(stepCtx provider.StepCtx) {
		errs = execute(stepCtx)
		it.processStepErrors(stepCtx.CurrentStep(), errs)
	})

	return errs
}

func (it *test) processStepErrors(step *allure.Step, errs []error) {
	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		currentStep := step

		if tErr, ok := err.(errors.WithNameError); ok {
			if tErr.GetName() != "" {
				currentStep = allure.NewSimpleStep(tErr.GetName())
				currentStep.Status = allure.Failed
				currentStep.WithParent(step)
			}
		}

		if tErr, ok := err.(errors.ExpectedError); ok {
			if tErr.GetActual() != nil {
				currentStep.WithNewParameters("Actual", fmt.Sprintf("%v", tErr.GetActual()))
			}
			if tErr.GetExpected() != nil {
				currentStep.WithNewParameters("Expected", fmt.Sprintf("%v", tErr.GetExpected()))
			}
		}

		currentStep.WithAttachments(allure.NewAttachment("Error", allure.Text, []byte(err.Error())))
	}

	step.Status = allure.Failed
}
