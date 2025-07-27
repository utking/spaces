// Package static contains the static files for the web server.
// It uses the embed package to include the static files at compile time.
// The static files are stored in the StaticFiles variable.
package static

import "embed"

var (
	// StaticFiles contains the static files for the web server.
	//go:embed *
	StaticFiles embed.FS
)
