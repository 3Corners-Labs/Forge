package execution

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

type CommandRunner interface {
	Run(context.Context, string, ...string) error
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, name string, arguments ...string) error {
	command := exec.CommandContext(ctx, name, arguments...)
	command.Env = []string{"PATH=" + os.Getenv("PATH")}
	return command.Run()
}

type DockerConfig struct {
	Binary   string
	Image    string
	TempRoot string
	Memory   string
	CPUs     string
}

type DockerRunner struct {
	config DockerConfig
	run    CommandRunner
}

func NewDockerRunner(config DockerConfig, runner CommandRunner) (*DockerRunner, error) {
	if config.Binary == "" {
		config.Binary = "docker"
	}
	if runner == nil {
		runner = execRunner{}
	}
	if !strings.Contains(config.Image, "@sha256:") || config.TempRoot == "" ||
		config.Memory == "" || config.CPUs == "" {
		return nil, errors.New("Docker worker configuration must pin image and resources")
	}
	return &DockerRunner{config: config, run: runner}, nil
}

func (runner *DockerRunner) Run(ctx context.Context, request workercontract.Request, workspace string) (workercontract.CandidateBundle, workercontract.Result, error) {
	if err := request.Validate(time.Now().UTC()); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	root, err := os.MkdirTemp(runner.config.TempRoot, "forge-worker-")
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	cleanupComplete := false
	defer func() {
		if !cleanupComplete {
			_ = os.RemoveAll(root)
		}
	}()
	input := filepath.Join(root, "input")
	output := filepath.Join(root, "output")
	if err := os.Mkdir(input, 0o700); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	if err := os.Mkdir(output, 0o777); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	requestBytes, _ := json.Marshal(request)
	if err := os.WriteFile(filepath.Join(input, "request.json"), requestBytes, 0o600); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	arguments := []string{
		"run", "--rm", "--network", "none", "--read-only", "--cap-drop", "ALL",
		"--security-opt", "no-new-privileges", "--pids-limit", strconv.FormatInt(request.Limits.MaxProcesses, 10),
		"--memory", runner.config.Memory, "--cpus", runner.config.CPUs, "--user", "65532:65532",
		"--tmpfs", "/tmp:rw,noexec,nosuid,nodev,size=16m",
		"--mount", "type=bind,src=" + workspace + ",dst=/workspace,rw",
		"--mount", "type=bind,src=" + input + ",dst=/input,readonly",
		"--mount", "type=bind,src=" + output + ",dst=/output,rw",
		runner.config.Image, "/forge-worker", "--request", "/input/request.json",
		"--workspace", "/workspace", "--bundle", "/output/bundle.json", "--result", "/output/result.json",
	}
	if err := runner.run.Run(ctx, runner.config.Binary, arguments...); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	bundleBytes, err := readBounded(filepath.Join(output, "bundle.json"), request.Limits.MaxOutputBytes)
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, errors.New("worker bundle output is unavailable or oversized")
	}
	resultBytes, err := readBounded(filepath.Join(output, "result.json"), request.Limits.MaxOutputBytes)
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, errors.New("worker result output is unavailable or oversized")
	}
	bundle, err := workercontract.DecodeBundle(bundleBytes)
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	result, err := workercontract.DecodeResult(resultBytes)
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	if err := os.RemoveAll(root); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, errors.New("worker cleanup could not be proven")
	}
	cleanupComplete = true
	result.CleanupOutcome = "completed"
	return bundle, result, nil
}

func readBounded(name string, maximum int64) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	value, err := io.ReadAll(io.LimitReader(file, maximum+1))
	if err != nil || int64(len(value)) > maximum {
		return nil, errors.New("bounded output limit exceeded")
	}
	return value, nil
}
