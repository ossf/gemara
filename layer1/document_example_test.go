package layer1

import (
	"fmt"
	"os"
	"text/template"
)

// Adapted from: https://github.com/finos/ai-governance-framework/blob/main/docs/_mitigations/mi-11_human-feedback-loop-for-ai-systems.md

func ExampleGuidanceDocument() {
	tmpl := `
# {{ .Metadata.Title }} ({{ .Metadata.Id }})
---
**Front Matter:** {{ .FrontMatter }}
---
{{ range .Categories }}
### {{ .Title }} ({{ .Id }})
{{ .Description }}
#### Guidelines:
{{ range .Guidelines }}
##### {{ .Title }} ({{ .Id }})
**Objective:** {{ .Objective }}
{{ if .SeeAlso }}
**See Also:** {{ range .SeeAlso }}{{ . }} {{ end }}
{{ end }}
{{ end }}
{{ end }}
`

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
	//
	// ### Detective (DET)
	//Detection and Continuous Improvement
	//#### Guidelines:
	//
	//##### Human Feedback Loop for AI Systems (AIR-DET-011)
	//**Objective:** A Human Feedback Loop is a critical detective and continuous improvement mechanism that involves systematically collecting, analyzing, and acting upon feedback provided by human users, subject matter experts (SMEs), or reviewers regarding an AI system’s performance, outputs, or behavior.
	//
	//**See Also:** AIR-DET-015 AIR-DET-004 AIR-PREV-005
}
