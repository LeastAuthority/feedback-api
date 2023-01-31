package main

const fullFeedbackTemplate = `{{if .Content.Title}}Title: {{.Content.Title}}

From: {{.Channel}}

Rate: {{if .Content.Rate.Value}}{{.Content.Rate.Value}} ({{.Content.Rate.Type}}){{else}}not rated{{end}}
{{range .Content.Questions}}
Q: {{.Question}}
A: {{.Answer}}
{{end}}{{end}}`
