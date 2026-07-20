package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Organization struct{ Name string }
	Runtime      struct{ Mode, Profile string }
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
	seen := make(map[string]bool)
	validSections := map[string]bool{
		"organization": true,
		"runtime":      true,
		"intelligence": true,
	}
	s := bufio.NewScanner(f)
	lineNumber := 0
	for s.Scan() {
		lineNumber++
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") && !strings.Contains(strings.TrimSuffix(line, ":"), " ") {
			section = strings.TrimSuffix(line, ":")
			if !validSections[section] {
				return cfg, fmt.Errorf("line %d: unknown section %q", lineNumber, section)
			}
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return cfg, fmt.Errorf("line %d: expected key: value", lineNumber)
		}
		if section == "" {
			return cfg, fmt.Errorf("line %d: key outside a section", lineNumber)
		}
		key, value := strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		path := section + "." + key
		if seen[path] {
			return cfg, fmt.Errorf("line %d: duplicate key %q", lineNumber, path)
		}
		seen[path] = true
		switch path {
		case "organization.name":
			cfg.Organization.Name = value
		case "runtime.mode":
			cfg.Runtime.Mode = value
		case "runtime.profile":
			cfg.Runtime.Profile = value
		case "intelligence.enabled":
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				return cfg, fmt.Errorf("line %d: intelligence.enabled must be true or false", lineNumber)
			}
			cfg.Intelligence.Enabled = enabled
		default:
			return cfg, fmt.Errorf("line %d: unknown key %q", lineNumber, path)
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
	if cfg.Runtime.Profile == "" {
		return cfg, errors.New("runtime.profile is required")
	}
	return cfg, nil
}
