package domain

// RequestPageMeta is a struct that represents the pagination parameters for a request.
type RequestPageMeta struct {
	Limit uint8 `query:"limit"`
	Page  uint  `query:"page"`
}
