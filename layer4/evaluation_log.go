package layer4

import (
	"encoding/json"
	"fmt"
)

// EvaluationLog contains the results of evaluating a set of Layer 4 controls.
type EvaluationLog struct {
	Evaluations []*ControlEvaluation `yaml:"evaluations"`
	Metadata    Metadata             `yaml:"metadata"`
}

// ToSARIF converts the evaluation results into a SARIF document (v2.1.0).
// Each AssessmentLog is emitted as a SARIF result. The rule id is derived from
// the control id and requirement id.
func (e EvaluationLog) ToSARIF() ([]byte, error) {
	report := &SarifReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/123e95847b13fbdd4cbe2120fa5e33355d4a042b/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
	}
	driver := ToolComponent{
		Name:           e.Metadata.Evaluator.Name,
		InformationURI: e.Metadata.Evaluator.URI,
		Version:        e.Metadata.Evaluator.Version,
	}
	run := Run{Tool: Tool{Driver: driver}}

	// Build a simple in-memory set of rules to avoid duplicates
	ruleIdSeen := map[string]bool{}
	rules := []ReportingDescriptor{}

	for _, evaluation := range e.Evaluations {
		if evaluation == nil {
			continue
		}

		for _, log := range evaluation.AssessmentLogs {
			if log == nil {
				continue
			}

			ruleID := fmt.Sprintf("%s/%s", evaluation.ControlID, log.RequirementId)
			if !ruleIdSeen[ruleID] {
				rule := ReportingDescriptor{ID: ruleID}
				if log.Description != "" {
					rule.Name = log.Description
				} else if evaluation.Name != "" {
					rule.Name = evaluation.Name
				}
				rules = append(rules, rule)
				ruleIdSeen[ruleID] = true
			}

			level := mapResultToSarifLevel(log.Result)

			// Message: prefer specific message, fallback to description
			msg := log.Message
			if msg == "" {
				msg = log.Description
			}

			result := ResultEntry{
				RuleID:  ruleID,
				Level:   level,
				Message: Message{Text: msg},
				Locations: []Location{
					{
						LogicalLocations: []LogicalLocation{
							{FullyQualifiedName: ruleID},
						},
					},
				},
			}
			run.Results = append(run.Results, result)
		}
	}

	// attach rules if any
	if len(rules) > 0 {
		run.Tool.Driver.Rules = rules
	}

	report.Runs = append(report.Runs, run)
	return json.Marshal(report)
}

func mapResultToSarifLevel(r Result) string {
	switch r {
	case Failed:
		return "error"
	case NeedsReview, Unknown:
		return "warning"
	case Passed, NotApplicable, NotRun:
		fallthrough
	default:
		return "note"
	}
}

// Minimal SARIF v2.1.0 model we need for export without external deps
type SarifReport struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool    Tool          `json:"tool"`
	Results []ResultEntry `json:"results,omitempty"`
}

type Tool struct {
	Driver ToolComponent `json:"driver"`
}

type ToolComponent struct {
	Name                  string                `json:"name"`
	InformationURI        string                `json:"informationUri,omitempty"`
	Version               string                `json:"version,omitempty"`
	SemanticVersion       string                `json:"semanticVersion,omitempty"`
	DottedQuadFileVersion string                `json:"dottedQuadFileVersion,omitempty"`
	Rules                 []ReportingDescriptor `json:"rules,omitempty"`
}

type ReportingDescriptor struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type ResultEntry struct {
	RuleID    string     `json:"ruleId"`
	Level     string     `json:"level,omitempty"`
	Message   Message    `json:"message"`
	Locations []Location `json:"locations,omitempty"`
}

type Message struct {
	Text string `json:"text"`
}

type Location struct {
	LogicalLocations []LogicalLocation `json:"logicalLocations,omitempty"`
}

type LogicalLocation struct {
	FullyQualifiedName string `json:"fullyQualifiedName,omitempty"`
}
