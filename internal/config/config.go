package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Organization struct{ Name string }
	Runtime      struct{ Mode string }
	Intelligence struct{ Enabled bool }
}

func Load(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	section := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") && !strings.Contains(strings.TrimSuffix(line, ":"), " ") {
			section = strings.TrimSuffix(line, ":")
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		switch section + "." + key {
		case "organization.name":
			cfg.Organization.Name = value
		case "runtime.mode":
			cfg.Runtime.Mode = value
		case "intelligence.enabled":
			cfg.Intelligence.Enabled = strings.EqualFold(value, "true")
		}
	}
	if err := s.Err(); err != nil {
		return cfg, err
	}
	if cfg.Organization.Name == "" {
		return cfg, errors.New("organization.name is required")
	}
	if cfg.Runtime.Mode == "" {
		cfg.Runtime.Mode = "compose"
	}
	return cfg, nil
}
