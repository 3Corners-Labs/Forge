package doctor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Check struct {
	Name    string
	Command string
	OK      bool
}

func Run() []Check {
	commands := []struct{ name, command string }{
		{"Git", "git"},
		{"Go", "go"},
		{"Docker", "docker"},
		{"Task", "task"},
	}
	checks := make([]Check, 0, len(commands))
	for _, c := range commands {
		_, err := exec.LookPath(c.command)
		checks = append(checks, Check{Name: c.name, Command: c.command, OK: err == nil})
	}
	return checks
}

func RunOperational() []Check {
	checks := Run()
	docker := Check{Name: "Docker daemon", Command: "docker info"}
	command := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	docker.OK = command.Run() == nil
	checks = append(checks, docker)
	for _, variable := range []string{
		"FORGE_M1_POSTGRES_DSN",
		"FORGE_M1_MINIO_ENDPOINT",
		"FORGE_M1_MINIO_ACCESS_KEY",
		"FORGE_M1_MINIO_SECRET_KEY",
	} {
		checks = append(checks, Check{
			Name:    variable,
			Command: variable,
			OK:      os.Getenv(variable) != "",
		})
	}
	return checks
}

func Print(checks []Check) bool {
	return PrintTo(os.Stdout, checks)
}

func PrintTo(output io.Writer, checks []Check) bool {
	allOK := true
	for _, c := range checks {
		status := "OK"
		if !c.OK {
			status = "MISSING"
			allOK = false
		}
		fmt.Fprintf(output, "%-30s %s\n", c.Name, status)
	}
	return allOK
}
