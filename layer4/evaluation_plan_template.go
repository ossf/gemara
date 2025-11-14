package layer4

// MarkdownTemplate is the default template for generating markdown checklist output.
// This template can be customized or replaced to change the output format without
// modifying the data extraction logic.
const MarkdownTemplate = `{{if .PlanId}}# Evaluation Plan: {{.PlanId}}

{{end}}{{if .Author}}**Author:** {{.Author}}{{if .AuthorVersion}} (v{{.AuthorVersion}}){{end}}

{{end}}{{range $index, $section := .Sections}}{{if $index}}
---

{{end}}## {{$section.ControlName}}

{{if $section.ControlReference}}**Control:** {{$section.ControlReference}}

{{end}}{{if eq (len $section.Items) 0}}- [ ] No assessments defined
{{else}}{{range $section.Items}}{{if .IsAdditionalProcedure}}  {{end}}- [ ] {{if .RequirementId}}**{{.RequirementId}}**: {{end}}{{.ProcedureName}}{{if and .Description (ne .Description .ProcedureName)}} - {{.Description}}{{end}}
{{if .Documentation}}    > [Documentation]({{.Documentation}})
{{end}}{{end}}{{end}}{{end}}`
