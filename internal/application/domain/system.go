package domain

// SystemStats represents the system statistics.
type SystemStats struct {
	ActiveUsers   int64  `json:"active_users"`
	InactiveUsers int64  `json:"inactive_users"`
	NoteTags      int64  `json:"note_tags"`
	SecretTags    int64  `json:"secret_tags"`
	Notes         int64  `json:"notes"`
	Secrets       int64  `json:"secrets"`
	Bookmarks     int64  `json:"bookmarks"`
	BookmarkTags  int64  `json:"bookmark_tags"`
	MemoryAlloc   int64  `json:"memory_alloc"`
	Allocations   uint64 `json:"mallocs"`
	TotalCPU      int    `json:"total_cpu"`
}
