package verification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

type CommandRunner interface {
	Run(context.Context, string, ...string) error
}

type commandRunner struct{}

func (commandRunner) Run(ctx context.Context, name string, arguments ...string) error {
	command := exec.CommandContext(ctx, name, arguments...)
	command.Env = []string{"PATH=" + os.Getenv("PATH")}
	output, err := command.CombinedOutput()
	if err != nil {
		if len(output) > 4096 {
			output = output[:4096]
		}
		return fmt.Errorf("isolated command failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

type DockerConfig struct {
	Binary         string
	Image          string
	TempRoot       string
	Memory         string
	CPUs           string
	PIDs           string
	MaxOutputBytes int64
}

type DockerRunner struct {
	config DockerConfig
	run    CommandRunner
}

func NewDockerRunner(config DockerConfig, run CommandRunner) (*DockerRunner, error) {
	if config.Binary == "" {
		config.Binary = "docker"
	}
	if run == nil {
		run = commandRunner{}
	}
	if !pinnedVerifierImage(config.Image) || !filepath.IsAbs(config.TempRoot) ||
		config.Memory == "" || config.CPUs == "" || config.PIDs == "" || config.MaxOutputBytes < 1 {
		return nil, errors.New("Docker verifier configuration must pin image and resources")
	}
	return &DockerRunner{config: config, run: run}, nil
}

func pinnedVerifierImage(value string) bool {
	digest := value
	if index := strings.LastIndex(value, "@sha256:"); index >= 0 {
		digest = value[index+1:]
	}
	if !strings.HasPrefix(digest, "sha256:") || len(digest) != 71 {
		return false
	}
	for _, character := range digest[7:] {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func (runner *DockerRunner) Verify(ctx context.Context, plan verificationcontract.AssurancePlan, bundle workercontract.CandidateBundle, workspace string) (verificationcontract.Result, error) {
	if err := bundle.Validate(); err != nil {
		return verificationcontract.Result{}, err
	}
	root, err := os.MkdirTemp(runner.config.TempRoot, "forge-verifier-")
	if err != nil {
		return verificationcontract.Result{}, err
	}
	clean := false
	defer func() {
		if !clean {
			_ = os.RemoveAll(root)
		}
	}()
	input, output := filepath.Join(root, "input"), filepath.Join(root, "output")
	if err := os.Mkdir(input, 0o700); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.Mkdir(output, 0o777); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.Chmod(output, 0o777); err != nil {
		return verificationcontract.Result{}, err
	}
	planBytes, _ := json.Marshal(plan)
	bundleBytes, _ := json.Marshal(bundle)
	if err := os.WriteFile(filepath.Join(input, "plan.json"), planBytes, 0o644); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.WriteFile(filepath.Join(input, "bundle.json"), bundleBytes, 0o644); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.Chmod(filepath.Join(input, "plan.json"), 0o644); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.Chmod(filepath.Join(input, "bundle.json"), 0o644); err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.Chmod(input, 0o755); err != nil {
		return verificationcontract.Result{}, err
	}
	arguments := []string{
		"run", "--rm", "--network", "none", "--read-only", "--cap-drop", "ALL",
		"--security-opt", "no-new-privileges", "--pids-limit", runner.config.PIDs,
		"--memory", runner.config.Memory, "--cpus", runner.config.CPUs, "--user", "65532:65532",
		"--env", "GOMAXPROCS=1",
		"--tmpfs", "/tmp:rw,noexec,nosuid,nodev,size=16m",
		"--mount", "type=bind,src=" + workspace + ",dst=/workspace,readonly",
		"--mount", "type=bind,src=" + input + ",dst=/input,readonly",
		"--mount", "type=bind,src=" + output + ",dst=/output",
		runner.config.Image, "/forge-verifier", "--plan", "/input/plan.json",
		"--bundle", "/input/bundle.json", "--workspace", "/workspace", "--result", "/output/result.json",
	}
	if err := runner.run.Run(ctx, runner.config.Binary, arguments...); err != nil {
		return verificationcontract.Result{}, err
	}
	resultBytes, err := readVerifierOutput(filepath.Join(output, "result.json"), runner.config.MaxOutputBytes)
	if err != nil {
		return verificationcontract.Result{}, errors.New("verifier result output is unavailable or oversized")
	}
	result, err := verificationcontract.DecodeResult(resultBytes)
	if err != nil {
		return verificationcontract.Result{}, err
	}
	if err := os.RemoveAll(root); err != nil {
		return verificationcontract.Result{}, errors.New("verifier cleanup could not be proven")
	}
	clean = true
	return result, nil
}

func readVerifierOutput(name string, maximum int64) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	value, err := io.ReadAll(io.LimitReader(file, maximum+1))
	if err != nil || int64(len(value)) > maximum {
		return nil, errors.New("bounded verifier output limit exceeded")
	}
	return value, nil
}
