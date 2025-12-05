package checklist

// markdownTemplate is the default template for generating checklist checklist output.
// This template is used internally by ToMarkdownChecklist().
const markdownTemplate = `{{if .PolicyId}}# Evaluation Plan: {{.PolicyId}}

{{end}}{{if .Author}}**Author:** {{.Author}}{{if .AuthorVersion}} (v{{.AuthorVersion}}){{end}}

{{end}}{{range $index, $section := .Sections}}{{if $index}}
---

{{end}}## {{$section.ControlRef}}

{{range $reqIndex, $requirement := $section.Requirements}}{{if $reqIndex}}

{{end}}### {{$requirement.RequirementId}}

{{if $requirement.RequirementRecommendation}}{{$requirement.RequirementRecommendation}}

{{end}}{{if eq (len $requirement.Items) 0}}- [ ] No evaluators defined
{{else}}{{range $requirement.Items}}{{if .IsAdditionalEvaluator}}  {{end}}- [ ] {{if and .EvaluatorName (eq false .IsAdditionalEvaluator)}}**{{.EvaluatorName}}**{{end}}{{if and .EvaluatorName .Description}} - {{.Description}}{{else if .Description}}{{.Description}}{{else if .EvaluatorName}}{{.EvaluatorName}}{{end}}
{{if .Documentation}}    > [Documentation]({{.Documentation}})
{{end}}{{end}}{{end}}{{end}}{{end}}`
