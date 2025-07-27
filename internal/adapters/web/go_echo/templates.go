package goecho

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/utking/extemplate"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/handlers"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/config"
	"gogs.utking.net/utking/spaces/internal/infra/session"
	"gogs.utking.net/utking/spaces/internal/ports"
	"gogs.utking.net/utking/spaces/views"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	kbSize = 1024
	mbSize = kbSize * kbSize
	gbSize = kbSize * mbSize
)

// InitTemplates initializes the templates for the Echo framework.
// It uses several functions to format data for rendering.
func InitTemplates(e *echo.Echo, appLoger ports.LoggingService, cfg *config.Config) error {
	xt := extemplate.New().Funcs(template.FuncMap{
		"isHomeDir": func(fileName string) bool {
			return fileName == "" || fileName == "/"
		},
		"isImage":         helpers.FileIsImage,
		"isViewable":      helpers.FileIsViewable,
		"fileIconFromExt": helpers.FileIconNameFromExt,
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("Jan 2 2006 15:04:05")
		},
		"cmpToInt64Ref": func(j int64, i *int64) bool {
			if i == nil {
				return false
			}

			return *i == j
		},
		"formatNumber": func(num interface{}) string {
			switch v := num.(type) {
			case int, int8, int16, int32, uint, uint8, uint16, uint32, int64, uint64:
				p := message.NewPrinter(language.English)
				return p.Sprintf("%d", v)
			default:
				return "NaN"
			}
		},
		"bytesToHuman": func(bytes int64) string {
			if bytes < kbSize {
				return fmt.Sprintf("%d B", bytes)
			}

			if bytes < mbSize {
				return fmt.Sprintf("%.2f KB", float64(bytes)/kbSize)
			}

			if bytes < gbSize {
				return fmt.Sprintf("%.2f MB", float64(bytes)/mbSize)
			}

			return fmt.Sprintf("%.2f GB", float64(bytes)/gbSize)
		},
		"commaSeparated": func(items []string) string {
			if len(items) == 0 {
				return ""
			}

			return strings.Join(items, ", ")
		},
		"cleanPath": func(path string) string {
			return filepath.Clean("/" + path) // Ensure the path starts with a slash
		},
	})

	if err := xt.ParseFS(
		views.TemplateFiles,
		[]string{".tmpl", ".html"},
	); err != nil {
		return err
	}

	t := &Template{
		worker: xt,
		logger: appLoger,
		cfg:    cfg,
	}
	e.Renderer = t

	return nil
}

// Template is a struct that implements the echo.Renderer interface.
type Template struct {
	worker *extemplate.Extemplate
	logger ports.LoggingService
	cfg    *config.Config
}

// Render is a method that renders the template with the given name and data.
func (t *Template) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	menu := handlers.GetMenu()
	username, _ := session.GetSessionUsername(c)

	isAdmin, _ := session.IsAdminSession(c)
	if !isAdmin {
		menu.AdminItems = nil
	}

	err := t.worker.ExecuteTemplate(
		w,
		name,
		map[string]interface{}{
			"data":     data,
			"title":    t.cfg.GetAppName(),
			"year":     time.Now().Year(),
			"menu":     menu,
			"username": strings.TrimSpace(username),
			"version":  helpers.GetReleaseVersion(),
		})
	if err != nil {
		t.logger.Error(
			c.Request().Context(),
			fmt.Sprintf("Error rendering template: %v", err),
		)
	}

	return err
}
