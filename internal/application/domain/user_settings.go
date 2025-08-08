package domain

import "encoding/json"

type UserSettings struct {
	DarkModeEnabled  bool `json:"dark_mode_enabled"`
	FileBrowserTiles bool `json:"file_browser_tiles"`
}

// ToJSON converts the UserSettings to a JSON string.
func (s *UserSettings) ToJSON() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// FromJSON populates the UserSettings from a JSON string.
func (UserSettings) FromJSON(jsonStr string) (*UserSettings, error) {
	s := new(UserSettings)
	err := json.Unmarshal([]byte(jsonStr), s)
	return s, err
}
