package cute

import (
	"fmt"

	"github.com/ozontech/testo"

	allure "github.com/ozontech/testo-allure"

	"github.com/ozontech/cute/errors"
)

func (it *Test) executeWithStep(t internalT, stepName string, execute func(t T) []error) []error {
	var (
		errs []error
	)

	// Add attempt indication in Allure if more than 1 attempt
	if it.Retry.MaxAttempts != 1 {
		stepName = fmt.Sprintf("[Attempt #%d] %v", it.Retry.currentCount, stepName)
	}

	allure.Step(t, stepName, func(stepT T) {
		errs = execute(stepT)
		processStepErrors(stepT, errs)
	})

	return errs
}

func processStepErrors(stepT T, errs []error) {
	var (
		statuses = make([]allure.Status, 0)
	)

	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		currentStatus := allure.StatusFailed

		if tErr, ok := err.(errors.OptionalError); ok {
			if tErr.IsOptional() {
				currentStatus = allure.StatusSkipped
			}
		}

		if tErr, ok := err.(errors.BrokenError); ok {
			if tErr.IsBroken() {
				currentStatus = allure.StatusBroken
			}
		}

		if tErr, ok := err.(errors.WithNameError); ok {
			// Record the error as a named sub-step with its own status
			status, stepErr := currentStatus, err

			testo.Run(stepT, tErr.GetName(), func(subT T) {
				subT.Status(status)
				addErrorDetails(subT, stepErr)
			})
		} else {
			addErrorDetails(stepT, err)
		}

		statuses = append(statuses, currentStatus)
	}

	// If one error was not optional, parent step should be failed
	for _, status := range statuses {
		stepT.Status(status)

		if status == allure.StatusFailed {
			break
		}
	}
}

func addErrorDetails(t T, err error) {
	if tErr, ok := err.(errors.WithFields); ok {
		for k, v := range tErr.GetFields() {
			if v == nil {
				continue
			}

			t.Parameters(allure.NewParameter(k, v))
		}
	}

	if tErr, ok := err.(errors.WithAttachments); ok {
		for _, v := range tErr.GetAttachments() {
			if v == nil {
				continue
			}

			t.Attach(v.Name, allure.Bytes(v.Content).As(allure.MediaType(v.MimeType)))
		}
	}

	t.Attach("Error", allure.Bytes(err.Error()).As(allure.TextPlain))
}
