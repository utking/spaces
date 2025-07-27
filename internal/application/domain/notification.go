package domain

import (
	"errors"
	"strings"
)

// Notification represents a notification message that can be sent to users.
type Notification struct {
	Title   string
	Message string
	To      string
}

// Trim trims the notification fields' whitespace.
func (n *Notification) Trim() {
	n.Title = strings.TrimSpace(n.Title)
	n.Message = strings.TrimSpace(n.Message)
	n.To = strings.TrimSpace(n.To)
}

// Validate checks if the notification has valid data.
func (n *Notification) Validate() error {
	var err error

	if n.Title == "" {
		err = errors.Join(err, errors.New("notification title cannot be empty"))
	}

	if n.Message == "" {
		err = errors.Join(err, errors.New("notification message cannot be empty"))
	}

	if n.To == "" {
		err = errors.Join(err, errors.New("notification recipient cannot be empty"))
	}

	return err
}
