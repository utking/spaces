package domain

// Status represents the status of an entity in the system.
type Status struct {
	Title string
	ID    int64
}

// String returns the string representation of the status.
func (s Status) String() string {
	return s.Title
}
