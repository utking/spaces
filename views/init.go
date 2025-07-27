// Package views provides the template files for the web application.
// It uses the embed package to include the template files at compile time.
// The template files are stored in the TemplateFiles variable.
package views

import "embed"

var (
	// TemplateFiles contains the template files for the web application.
	//go:embed *
	TemplateFiles embed.FS
)
