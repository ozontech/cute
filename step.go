package cute

import (
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute/errors"
)

func executeWithStep(it *Test, t internalT, stepName string, execute func(t T) []error) []error {
	var (
		errs []error
	)
	// Add attempt indication in Allure if more than 1 attempt
	if it.Retry.MaxAttempts > 0 {
		stepName = fmt.Sprintf("[Attempt #%d] %v", it.Retry.currentCount+1, stepName)
	}
	t.WithNewStep(stepName, func(stepCtx provider.StepCtx) {
		errs = execute(stepCtx)
		processStepErrors(stepCtx, errs)
	})

	return errs
}

func processStepErrors(stepCtx provider.StepCtx, errs []error) {
	var (
		step     = stepCtx.CurrentStep()
		statuses = make([]allure.Status, 0)
	)

	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		currentStatus := allure.Failed
		currentStep := step

		if tErr, ok := err.(errors.OptionalError); ok {
			if tErr.IsOptional() {
				currentStatus = allure.Skipped
			}
		}

		if tErr, ok := err.(errors.BrokenError); ok {
			if tErr.IsBroken() {
				currentStatus = allure.Broken
			}
		}

		if tErr, ok := err.(errors.WithNameError); ok {
			currentStep = allure.NewSimpleStep(tErr.GetName())
			currentStep.Status = currentStatus
			currentStep.WithParent(step)
		}

		if tErr, ok := err.(errors.WithFields); ok {
			for k, v := range tErr.GetFields() {
				if v == nil {
					continue
				}

				currentStep.WithNewParameters(k, v)
			}
		}

		if tErr, ok := err.(errors.WithAttachments); ok {
			for _, v := range tErr.GetAttachments() {
				if v == nil {
					continue
				}

				currentStep.WithAttachments(allure.NewAttachment(v.Name, allure.MimeType(v.MimeType), v.Content))
			}
		}

		statuses = append(statuses, currentStatus)

		currentStep.WithAttachments(allure.NewAttachment("Error", allure.Text, []byte(err.Error())))
	}

	// If one error was not optional, parent step should be failed
	for _, status := range statuses {
		step.Status = status

		if status == allure.Failed {
			break
		}
	}
}
