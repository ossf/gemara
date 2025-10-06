package layer4

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// AssessmentStep is a function type that inspects the provided targetData and returns a Result with a message.
// The message may be an error string or other descriptive text.
type AssessmentStep func(payload interface{}) (Result, string)

func (as AssessmentStep) String() string {
	// Get the function pointer correctly
	fn := runtime.FuncForPC(reflect.ValueOf(as).Pointer())
	if fn == nil {
		return "<unknown function>"
	}
	return fn.Name()
}

func (as AssessmentStep) MarshalJSON() ([]byte, error) {
	return json.Marshal(as.String())
}

func (as AssessmentStep) MarshalYAML() (interface{}, error) {
	return as.String(), nil
}

// NewAssessment creates a new AssessmentLog object and returns a pointer to it.
func NewAssessment(requirementId string, description string, applicability []string, steps []AssessmentStep) (*AssessmentLog, error) {
	a := &AssessmentLog{
		Requirement: Mapping{
			EntryId: requirementId,
		},
		Description:   description,
		Applicability: applicability,
		Result:        NotRun,
		Steps:         steps,
	}
	err := a.precheck()
	return a, err
}

// AddStep queues a new step in the AssessmentLog
func (a *AssessmentLog) AddStep(step AssessmentStep) {
	a.Steps = append(a.Steps, step)
}

func (a *AssessmentLog) runStep(targetData interface{}, step AssessmentStep) Result {
	a.StepsExecuted++
	result, message := step(targetData)
	a.Result = UpdateAggregateResult(a.Result, result)
	a.Message = message
	return result
}

// Run will execute all steps, halting if any step does not return layer4.Passed.
func (a *AssessmentLog) Run(targetData interface{}) Result {
	a.Result = NotRun
	if a.Result != NotRun {
		return a.Result
	}

	a.Start = Datetime(time.Now().Format(time.RFC3339))
	err := a.precheck()
	if err != nil {
		a.Result = Unknown
		return a.Result
	}
	for _, step := range a.Steps {
		if a.runStep(targetData, step) == Failed {
			return Failed
		}
	}
	a.End = Datetime(time.Now().Format(time.RFC3339))
	return a.Result
}

// precheck verifies that the assessment has all the required fields.
// It returns an error if the assessment is not valid.
func (a *AssessmentLog) precheck() error {
	if a.Requirement.EntryId == "" || a.Description == "" || a.Applicability == nil || a.Steps == nil || len(a.Applicability) == 0 || len(a.Steps) == 0 {
		message := fmt.Sprintf(
			"expected all AssessmentLog fields to have a value, but got: requirementId=len(%v), description=len=(%v), applicability=len(%v), steps=len(%v)",
			len(a.Requirement.EntryId), len(a.Description), len(a.Applicability), len(a.Steps),
		)
		a.Result = Unknown
		a.Message = message
		return errors.New(message)
	}

	return nil
}
