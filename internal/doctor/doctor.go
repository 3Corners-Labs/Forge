package doctor

import (
	"fmt"
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

func Print(checks []Check) bool {
	allOK := true
	for _, c := range checks {
		status := "OK"
		if !c.OK {
			status = "MISSING"
			allOK = false
		}
		fmt.Printf("%-12s %s\n", c.Name, status)
	}
	return allOK
}
