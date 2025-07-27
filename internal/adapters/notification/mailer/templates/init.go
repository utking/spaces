// Package templates provides access to email templates used in the mailer service.
package templates

import "embed"

//go:embed *.html
var MailerTemplates embed.FS
