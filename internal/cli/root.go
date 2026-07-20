package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/3corners/forge/internal/config"
	"github.com/3corners/forge/internal/doctor"
	"github.com/3corners/forge/internal/platform"
)

func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}
	switch os.Args[1] {
	case "version":
		fmt.Println("forge 0.1.0")
		return nil
	case "doctor":
		if doctor.Print(doctor.Run()) {
			return nil
		}
		return errors.New("one or more prerequisites are missing")
	case "init":
		return runInit(os.Args[2:])
	case "up":
		return runUp()
	case "down":
		return platform.Compose("down")
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}
}

func printUsage() {
	fmt.Println("3Corners Forge engineering platform CLI")
	fmt.Println("Usage: forge <command>")
	fmt.Println("Commands: init, doctor, up, down, version")
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	path := fs.String("file", "forge.yaml", "configuration file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if _, err := os.Stat(*path); err == nil {
		return fmt.Errorf("%s already exists", *path)
	}
	content := "organization:\n  name: 3Corners\nruntime:\n  mode: compose\n  profile: developer\nintelligence:\n  enabled: false\n"
	return os.WriteFile(*path, []byte(content), 0o644)
}

func runUp() error {
	path, err := filepath.Abs("forge.yaml")
	if err != nil {
		return err
	}
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	if cfg.Runtime.Mode != "compose" {
		return fmt.Errorf("runtime mode %q is not implemented", cfg.Runtime.Mode)
	}
	if cfg.Runtime.Profile != "developer" {
		return fmt.Errorf("compose bootstrap is developer-only; refusing runtime profile %q", cfg.Runtime.Profile)
	}
	return platform.Compose("up")
}
