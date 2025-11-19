package layer4

// markdownTemplate is the default template for generating markdown checklist output.
// This template is used internally by ToMarkdownChecklist().
const markdownTemplate = `{{if .PlanId}}# Evaluation Plan: {{.PlanId}}

{{end}}{{if .Author}}**Author:** {{.Author}}{{if .AuthorVersion}} (v{{.AuthorVersion}}){{end}}

{{end}}{{range $index, $section := .Sections}}{{if $index}}
---

{{end}}## {{$section.ControlName}}

{{if $section.ControlReference}}**Control:** {{$section.ControlReference}}

{{end}}{{if eq (len $section.Items) 0}}- [ ] No assessments defined
{{else}}{{range $section.Items}}{{if .IsAdditionalProcedure}}  {{end}}- [ ] {{if and .RequirementId (eq false .IsAdditionalProcedure)}}**{{.RequirementId}}**: {{end}}{{.ProcedureName}}{{if and .Description (ne .Description .ProcedureName)}} - {{.Description}}{{end}}
{{if .Documentation}}    > [Documentation]({{.Documentation}})
{{end}}{{end}}{{end}}{{end}}`
