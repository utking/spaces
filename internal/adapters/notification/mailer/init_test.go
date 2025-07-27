package mailer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gogs.utking.net/utking/spaces/internal/adapters/notification/mailer"
	"gogs.utking.net/utking/spaces/internal/application/domain"
)

func TestRenderTemplate(t *testing.T) {
	_, err := mailer.RenderTemplate(t.Context(), "test_template", map[string]interface{}{})
	assert.Error(t, err, "should return an error for non-existing template")

	// Test with a valid template - welcome.html
	rendered, renderErr := mailer.RenderTemplate(t.Context(), "welcome.html", map[string]interface{}{})
	if assert.NoError(t, renderErr) {
		assert.NotEmpty(t, rendered, "rendered template should not be empty")
	}
}

func TestSend(t *testing.T) {
	// testing the Send function only with data that will cause validation errors
	testMailer := mailer.New("smtp.example.com", 587, "user", "password", "fake-from@localhost", true)

	// empty title
	assert.Error(t,
		testMailer.Send(
			t.Context(),
			&domain.Notification{
				Title:   "",
				Message: "Test message",
				To:      "fake-to@localhost",
			},
		),
		"should return an error for empty title",
	)

	// empty message
	assert.Error(t,
		testMailer.Send(
			t.Context(),
			&domain.Notification{
				Title:   "Sample Title",
				Message: "",
				To:      "fake-to@localhost",
			},
		),
		"should return an error for empty title",
	)

	// empty recipient
	assert.Error(t,
		testMailer.Send(
			t.Context(),
			&domain.Notification{
				Title:   "Sample Title",
				Message: "Test message",
				To:      "",
			},
		),
		"should return an error for empty title",
	)
}
