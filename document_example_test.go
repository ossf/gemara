package gemara

import (
	"fmt"
	"os"
	"text/template"

	"github.com/goccy/go-yaml"
)

// Adapted from: https://github.com/finos/ai-governance-framework/blob/main/docs/_mitigations/mi-11_human-feedback-loop-for-ai-systems.md

func ExampleGuidanceDocument() {
	tmpl := `# {{ .Title }} ({{ .Metadata.Id }})
---
**Front Matter:** {{ .FrontMatter }}
---
{{- $doc := . }}{{ range .Families }}

### {{ .Title }} ({{ .Id }})
{{ .Description }}
#### Guidelines:
{{- $familyId := .Id }}{{ range $doc.Guidelines }}{{ if eq .Family $familyId }}

##### {{ .Title }} ({{ .Id }})
**Objective:** {{ .Objective }}
{{- if .SeeAlso }}

**See Also:** {{ range $index, $item := .SeeAlso }}{{ if $index }} {{ end }}{{ $item }}{{ end }}
{{- end }}

{{- end }}{{ end -}}
{{ end }}`

	l1Docs, err := goodAIGFExample()
	if err != nil {
		fmt.Printf("error getting testdata: %v\n", err)
		return
	}
	t, err := template.New("guidance").Parse(tmpl)
	if err != nil {
		fmt.Printf("error parsing template: %v\n", err)
		return
	}

	err = t.Execute(os.Stdout, l1Docs)
	if err != nil {
		fmt.Printf("error executing template: %v\n", err)
	}
	// Output:
	//# AI Governance Framework (FINOS-AIR)
	//---
	//**Front Matter:** The following framework has been developed by FINOS (Fintech Open Source Foundation).
	//---
	//
	//### Detective (DET)
	//Detection and Continuous Improvement
	//#### Guidelines:
	//
	//##### Human Feedback Loop for AI Systems (AIR-DET-011)
	//**Objective:** A Human Feedback Loop is a critical detective and continuous improvement mechanism that involves systematically collecting, analyzing, and acting upon feedback provided by human users, subject matter experts (SMEs), or reviewers regarding an AI system's performance, outputs, or behavior.
	//
	//**See Also:** AIR-DET-015 AIR-DET-004 AIR-PREV-005
	//
	//
	//##### Example Detective Control 004 (AIR-DET-004)
	//**Objective:** Placeholder control for testing references.
	//
	//
	//##### Example Detective Control 015 (AIR-DET-015)
	//**Objective:** Placeholder control for testing references.
	//
	//
	//
	//### Preventive (PREV)
	//Prevention and Risk Mitigation
	//#### Guidelines:
	//
	//##### Example Preventive Control 005 (AIR-PREV-005)
	//**Objective:** Placeholder control for testing references.
}

func goodAIGFExample() (GuidanceDocument, error) {
	testdataPath := "./test-data/good-aigf.yaml"
	data, err := os.ReadFile(testdataPath)
	if err != nil {
		return GuidanceDocument{}, err
	}
	var l1Docs GuidanceDocument
	if err := yaml.Unmarshal(data, &l1Docs); err != nil {
		return GuidanceDocument{}, err
	}
	return l1Docs, nil
}
