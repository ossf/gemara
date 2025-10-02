package layer4

// AddAssessment creates a new AssessmentLog object and adds it to the ControlEvaluation.
func (c *ControlEvaluation) AddAssessment(requirementId string, description string, applicability []string, steps []AssessmentStep) (assessment *AssessmentLog) {
	assessment, err := NewAssessment(requirementId, description, applicability, steps)
	if err != nil {
		c.Result = Failed
		c.Message = err.Error()
	}
	c.AssessmentLogs = append(c.AssessmentLogs, assessment)
	return
}

// Evaluate runs each step in each assessment, updating the relevant fields on the control evaluation.
// It will halt if a step returns a failed result. The targetData is the data that the assessment will be run against.
// The userApplicability is a slice of strings that determine when the assessment is applicable. The changesAllowed
// determines whether the assessment is allowed to execute its changes.
func (c *ControlEvaluation) Evaluate(targetData interface{}, userApplicability []string) {
	if len(c.AssessmentLogs) == 0 {
		c.Result = NeedsReview
		return
	}
	for _, assessment := range c.AssessmentLogs {
		var applicable bool
		for _, aa := range assessment.Applicability {
			for _, ua := range userApplicability {
				if aa == ua {
					applicable = true
					break
				}
			}
		}
		if applicable {
			result := assessment.Run(targetData)
			c.Result = UpdateAggregateResult(c.Result, result)
			c.Message = assessment.Message
			if c.Result == Failed {
				break
			}
		}
	}
}
