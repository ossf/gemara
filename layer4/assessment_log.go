package layer4

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// AssessmentLog is a struct that contains the results of a single step within a ControlEvaluation.
type AssessmentLog struct {
	// RequirementID is the unique identifier for the requirement being tested
	RequirementId string `yaml:"requirement-id"`
	// ProcedureId uniquely identifies the assessment procedure associated with the log
	ProcedureId string `json:"procedure-id,omitempty"`
	// Applicability is a slice of identifier strings to determine when this test is applicable
	Applicability []string `yaml:"applicability"`
	// Description is a human-readable description of the test
	Description string `yaml:"description"`
	// Result is true if the test passed
	Result Result `yaml:"result"`
	// Message is the human-readable result of the test
	Message string `yaml:"message"`
	// Steps is a slice of steps that were executed during the test
	Steps []AssessmentStep `yaml:"steps"`
	// StepsExecuted is the number of steps that were executed during the test
	StepsExecuted int `yaml:"steps-executed,omitempty"`
	// Start is the time the assessment run began.
	Start string `yaml:"start"`
	// End is the time the assessment run finished.
	// This is omitted if the assessment was interrupted or did not complete.
	End string `yaml:"end,omitempty"`
	// Value is the object that was returned during the test
	Value interface{} `yaml:"value,omitempty"`
	// Changes is a slice of changes that were made during the test
	Changes map[string]*Change `yaml:"changes,omitempty"`
	// Recommendation is a string to aid users in remediation, such as the text from a layer 2 assessment requirement
	Recommendation string `yaml:"recommendation,omitempty"`
}

// AssessmentStep is a function type that inspects the provided targetData and returns a Result with a message.
// The message may be an error string or other descriptive text.
type AssessmentStep func(payload interface{}, c map[string]*Change) (Result, string)

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
		RequirementId: requirementId,
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
	result, message := step(targetData, a.Changes)
	a.Result = UpdateAggregateResult(a.Result, result)
	a.Message = message
	return result
}

// Run will execute all steps, halting if any step does not return layer4.Passed.
func (a *AssessmentLog) Run(targetData interface{}, changesAllowed bool) Result {
	if a.Result != NotRun {
		return a.Result
	}

	a.Start = time.Now().Format(time.RFC3339)
	err := a.precheck()
	if err != nil {
		a.Result = Unknown
		return a.Result
	}
	for _, change := range a.Changes {
		if changesAllowed {
			change.Allow()
		}
	}
	for _, step := range a.Steps {
		if a.runStep(targetData, step) == Failed {
			return Failed
		}
	}
	a.End = time.Now().Format(time.RFC3339)
	return a.Result
}

// NewChange creates a new Change object and adds it to the AssessmentLog.
func (a *AssessmentLog) NewChange(
	changeName,
	targetName,
	description string,
	targetObject interface{},
	applyFunc ApplyFunc,
	revertFunc RevertFunc,
) *Change {
	change := NewChange(targetName, description, targetObject, applyFunc, revertFunc)
	if a.Changes == nil {
		a.Changes = make(map[string]*Change)
	}
	a.Changes[changeName] = &change
	return &change
}

// RevertChanges reverts all changes made by the assessment.
// It will not revert changes that have not been applied.
func (a *AssessmentLog) RevertChanges() (corrupted bool) {
	for _, change := range a.Changes {
		if !corrupted && (change.Applied || change.Error != nil) {
			if !change.Reverted {
				change.Revert(nil)
			}
			if change.Error != nil || !change.Reverted {
				corrupted = true // do not break loop here; continue attempting to revert all changes
			}
		}
	}
	return
}

// precheck verifies that the assessment has all the required fields.
// It returns an error if the assessment is not valid.
func (a *AssessmentLog) precheck() error {
	if a.RequirementId == "" || a.Description == "" || a.Applicability == nil || a.Steps == nil || len(a.Applicability) == 0 || len(a.Steps) == 0 {
		message := fmt.Sprintf(
			"expected all AssessmentLog fields to have a value, but got: requirementId=len(%v), description=len=(%v), applicability=len(%v), steps=len(%v)",
			len(a.RequirementId), len(a.Description), len(a.Applicability), len(a.Steps),
		)
		a.Result = Unknown
		a.Message = message
		return errors.New(message)
	}

	return nil
}
