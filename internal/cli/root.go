package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/3corners/forge/internal/config"
	"github.com/3corners/forge/internal/doctor"
	"github.com/3corners/forge/internal/platform"
)

func Execute() error {
	return execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, newDeveloperTrialJourney)
}

func execute(args []string, input io.Reader, output, errorOutput io.Writer, journeys journeyFactory) error {
	if len(args) < 1 {
		printUsageTo(output)
		return nil
	}
	switch args[0] {
	case "version":
		fmt.Fprintln(output, "forge 0.2.0")
		return nil
	case "doctor":
		checks := doctor.Run()
		if len(args) > 2 || (len(args) == 2 && args[1] != "--operational") {
			return errors.New("usage: forge doctor [--operational]")
		}
		if len(args) == 2 {
			checks = doctor.RunOperational()
		}
		if doctor.PrintTo(output, checks) {
			return nil
		}
		return errors.New("one or more prerequisites are missing")
	case "init":
		return runInit(args[1:])
	case "delivery":
		return runDelivery(args[1:], bufio.NewReader(input), output, errorOutput, journeys)
	case "workbench":
		return runWorkbench(args[1:], input, output, errorOutput, journeys)
	case "up":
		return runCompose("up")
	case "down":
		return runCompose("down")
	case "help", "-h", "--help":
		printUsageTo(output)
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage() {
	printUsageTo(os.Stdout)
}

func printUsageTo(output io.Writer) {
	fmt.Fprintln(output, "3Corners Forge engineering platform CLI")
	fmt.Fprintln(output, "Usage: forge <command>")
	fmt.Fprintln(output, "Commands: workbench, delivery, init, doctor, up, down, version")
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

func runCompose(action string) error {
	path, err := filepath.Abs("forge.yaml")
	if err != nil {
		return err
	}
	return runComposeWithConfig(path, action, platform.Compose)
}

func runComposeWithConfig(path, action string, invoke func(string) error) error {
	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	if err := cfg.ValidateDeveloperBootstrap(); err != nil {
		return err
	}
	return invoke(action)
}
