package config

import "testing"

func TestGetEnvValue(t *testing.T) {
	tests := []struct {
		name       string
		envVar     string
		defaultVal string
		want       string
	}{
		{"ExistingEnvVar", "EXISTING_VAR", "", "value"},
		{"NonExistingEnvVarWithDefault", "NON_EXISTING_VAR", "default_value", "default_value"},
		{"NonExistingEnvVarWithoutDefault", "NON_EXISTING_VAR", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVar == "EXISTING_VAR" {
				t.Setenv("EXISTING_VAR", "value")
			}
			got := getEnvValue(tt.envVar, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
