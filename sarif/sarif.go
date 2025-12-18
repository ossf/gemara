package sarif

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara"
)

var emptyArtifactURIMessage = "no file associated with this alert"

// FromEvaluationLog converts the evaluation results into a SARIF document (v2.1.0).
// Each AssessmentLog is emitted as a SARIF result. The rule id is derived from
// the control id and requirement id.
//
// Parameters:
//   - evaluationLog: The evaluation log to convert
//   - artifactURI: File path or URI for PhysicalLocation.artifactLocation.uri.
//     If empty, PhysicalLocation will be nil (no resource URI available).
//     For GitHub Code Scanning, typically use a file path like "README.md".
//   - catalog: Optional catalog data to enrich SARIF output with requirement text
//     and recommendations. If nil, only basic information is included.
//
// PhysicalLocation identifies the artifact (file/repository) where the result was found.
// LogicalLocation identifies the logical component (assessment step) that produced the result.
// Region is left nil as we don't have file-specific line/column data.
func FromEvaluationLog(evaluationLog gemara.EvaluationLog, artifactURI string, catalog *gemara.Catalog) ([]byte, error) {
	report := &SarifReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/123e95847b13fbdd4cbe2120fa5e33355d4a042b/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
	}
	driver := ToolComponent{
		Name:           evaluationLog.Metadata.Author.Name,
		InformationURI: evaluationLog.Metadata.Author.Uri,
		Version:        evaluationLog.Metadata.Author.Version,
	}
	run := Run{Tool: Tool{Driver: driver}}

	// Build a simple in-memory set of rules to avoid duplicates
	ruleIdSeen := map[string]bool{}
	rules := []ReportingDescriptor{}

	for _, evaluation := range evaluationLog.Evaluations {
		for _, log := range evaluation.AssessmentLogs {
			if log == nil {
				continue
			}

			// Skip NotRun and NotApplicable results - only include Passed, Failed, NeedsReview, Unknown
			if log.Result == gemara.NotRun || log.Result == gemara.NotApplicable {
				continue
			}

			ruleID := log.Requirement.EntryId
			if !ruleIdSeen[ruleID] {
				rule := ReportingDescriptor{ID: ruleID}
				if log.Description != "" {
					rule.Name = log.Description
				} else if evaluation.Name != "" {
					rule.Name = evaluation.Name
				}

				// Enrich rule with catalog data if available
				if catalog != nil {
					control, requirement := findControlAndRequirement(catalog, evaluation.Control.EntryId, log.Requirement.EntryId)
					if control != nil {
						// Use control title if name is still empty
						if rule.Name == "" {
							rule.Name = control.Title
						}

						// Add requirement text as short description
						if requirement != nil && requirement.Text != "" {
							rule.ShortDescription = &Message{Text: requirement.Text}
						}

						// Add full description: control objective + requirement text
						if control.Objective != "" {
							fullDesc := control.Objective
							if requirement != nil && requirement.Text != "" {
								fullDesc = fmt.Sprintf("%s\n\nRequirement: %s", control.Objective, requirement.Text)
							}
							rule.FullDescription = &Message{Text: fullDesc}
						}

						// Add recommendation as help text
						if log.Recommendation != "" {
							rule.Help = &Message{Text: log.Recommendation}
						} else if requirement != nil && requirement.Recommendation != "" {
							rule.Help = &Message{Text: requirement.Recommendation}
						}

						// HelpUri is left empty - catalog-specific URI generation should be handled by the caller
					}
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

			var physicalLocation *PhysicalLocation
			if artifactURI == "" {
				artifactURI = emptyArtifactURIMessage
			}
			physicalLocation = &PhysicalLocation{
				ArtifactLocation: ArtifactLocation{
					URI: artifactURI,
				},
			}

			// Use the last AssessmentStep for LogicalLocation (the location is for the entire evaluation)
			logicalLocationName := ruleID
			if len(log.Steps) > 0 {
				lastStep := log.Steps[len(log.Steps)-1]
				if lastStep != nil {
					logicalLocationName = lastStep.String()
				}
			}

			location := Location{
				PhysicalLocation: physicalLocation,
				LogicalLocations: []LogicalLocation{
					{FullyQualifiedName: logicalLocationName},
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

func mapResultToSarifLevel(r gemara.Result) string {
	switch r {
	case gemara.Failed:
		return "error"
	case gemara.NeedsReview, gemara.Unknown:
		return "warning"
	case gemara.Passed, gemara.NotApplicable, gemara.NotRun:
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
	ID               string   `json:"id"`
	Name             string   `json:"name,omitempty"`
	ShortDescription *Message `json:"shortDescription,omitempty"`
	FullDescription  *Message `json:"fullDescription,omitempty"`
	Help             *Message `json:"help,omitempty"`
	HelpUri          string   `json:"helpUri,omitempty"`
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

// findControlAndRequirement searches the catalog for a control and requirement by their IDs.
// Returns the control and requirement if found, nil otherwise.
func findControlAndRequirement(catalog *gemara.Catalog, controlID, requirementID string) (*gemara.Control, *gemara.AssessmentRequirement) {
	if catalog == nil {
		return nil, nil
	}

	for i := range catalog.Controls {
		control := &catalog.Controls[i]
		if control.Id == controlID {
			// Found the control, now find the requirement
			for j := range control.AssessmentRequirements {
				requirement := &control.AssessmentRequirements[j]
				if requirement.Id == requirementID {
					return control, requirement
				}
			}
			// Control found but requirement not found
			return control, nil
		}
	}

	return nil, nil
}
