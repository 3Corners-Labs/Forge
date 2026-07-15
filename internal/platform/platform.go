package platform

import (
	"fmt"
	"os"
	"os/exec"
)

func Compose(action string) error {
	args := []string{"compose", "-f", "deploy/compose/docker-compose.yml", action}
	if action == "up" {
		args = append(args, "-d", "--remove-orphans")
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose %s failed: %w", action, err)
	}
	return nil
}
