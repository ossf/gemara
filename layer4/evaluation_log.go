package layer4

import (
	"encoding/json"
	"fmt"
)

// ToSARIF converts the evaluation results into a SARIF document (v2.1.0).
// Each AssessmentLog is emitted as a SARIF result. The rule id is derived from
// the control id and requirement id.
//
// PhysicalLocation: Uses Metadata.Author.Uri as the artifact URI (typically a repository URL).
// IMPORTANT: GitHub Code Scanning may require file paths (e.g., "README.md") instead of repository URLs.
// If testing shows repository URLs don't work with GitHub Code Scanning, consider implementing
// an optional artifactURI parameter (see Option A implementation below).
func (e EvaluationLog) ToSARIF() ([]byte, error) {
	report := &SarifReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/123e95847b13fbdd4cbe2120fa5e33355d4a042b/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
	}
	driver := ToolComponent{
		Name:           e.Metadata.Author.Name,
		InformationURI: e.Metadata.Author.Uri,
		Version:        e.Metadata.Author.Version,
	}
	run := Run{Tool: Tool{Driver: driver}}

	// Build a simple in-memory set of rules to avoid duplicates
	ruleIdSeen := map[string]bool{}
	rules := []ReportingDescriptor{}

	for _, evaluation := range e.Evaluations {
		for _, log := range evaluation.AssessmentLogs {
			if log == nil {
				continue
			}

			ruleID := fmt.Sprintf("%s/%s", evaluation.Control.EntryId, log.Requirement.EntryId)
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

			// Build PhysicalLocation using repository URI from metadata
			// For repository-level assessments, use Metadata.Author.Uri as the artifact URI
			// Region is left nil as we don't have file-specific line/column data
			//
			// NOTE: This uses Metadata.Author.Uri (repository URL like "https://github.com/ossf/gemera")
			// GitHub Code Scanning may require file paths (e.g., "README.md") instead.
			// TODO: Test with GitHub Code Scanning. If repository URLs don't work, implement
			// Option A: Add optional artifactURI parameter to ToSARIF() method signature.
			var physicalLocation *PhysicalLocation
			if e.Metadata.Author.Uri != "" {
				physicalLocation = &PhysicalLocation{
					ArtifactLocation: ArtifactLocation{
						URI: e.Metadata.Author.Uri,
					},
					// Region left nil - no line/column data available for repository-level assessments
				}
			}

			location := Location{
				PhysicalLocation: physicalLocation,
				LogicalLocations: []LogicalLocation{
					{FullyQualifiedName: ruleID},
				},
			}

			result := ResultEntry{
				RuleID:  ruleID,
				Level:   level,
				Message: Message{Text: msg},
				Locations: []Location{
					location,
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
	PhysicalLocation *PhysicalLocation `json:"physicalLocation,omitempty"`
	LogicalLocations []LogicalLocation `json:"logicalLocations,omitempty"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           *Region          `json:"region,omitempty"`
}

type ArtifactLocation struct {
	URI       string `json:"uri"`
	URIBaseID string `json:"uriBaseId,omitempty"`
	Index     int    `json:"index,omitempty"`
}

type Region struct {
	StartLine   int      `json:"startLine,omitempty"`
	StartColumn int      `json:"startColumn,omitempty"`
	EndLine     int      `json:"endLine,omitempty"`
	EndColumn   int      `json:"endColumn,omitempty"`
	Snippet     *Snippet `json:"snippet,omitempty"`
}

type Snippet struct {
	Text string `json:"text"`
}

type LogicalLocation struct {
	FullyQualifiedName string `json:"fullyQualifiedName,omitempty"`
}
