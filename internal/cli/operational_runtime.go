package cli

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	forgecontext "github.com/3corners/forge/internal/context"
	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/execution"
	"github.com/3corners/forge/internal/verification"
)

const (
	runtimeContract = "contract"
	runtimeDocker   = "docker"

	m1ArchitecturePath = "docs/architecture/vertical-slice.md"
	m1TargetPath       = "docs/forge-m1-developer-trial.md"
)

type operationalRuntime struct {
	root           string
	runtimeRoot    string
	sourceRevision string
	catalog        *forgecontext.RepositoryCatalog
	worker         execution.Runner
	verifier       verification.Runner
	operational    bool
}

func newOperationalRuntime(options journeyOptions, now func() time.Time) (operationalRuntime, error) {
	configuredRoot := options.RepositoryRoot
	root, err := resolveRepositoryRoot(options.RepositoryRoot)
	if err != nil {
		return operationalRuntime{}, err
	}
	runtimeRoot := options.RuntimeRoot
	if runtimeRoot == "" {
		runtimeRoot = filepath.Join(root, ".forge", "runtime")
	} else if configuredRoot != "" {
		rawRoot, rootErr := filepath.Abs(configuredRoot)
		rawRuntime, runtimeErr := filepath.Abs(runtimeRoot)
		if rootErr == nil && runtimeErr == nil {
			if relative, relativeErr := filepath.Rel(rawRoot, rawRuntime); relativeErr == nil &&
				relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				runtimeRoot = filepath.Join(root, relative)
			}
		}
	}
	runtimeRoot, err = filepath.Abs(runtimeRoot)
	if err != nil || !withinRoot(root, runtimeRoot) {
		return operationalRuntime{}, errors.New("runtime root must be project-local")
	}
	mode := options.RuntimeMode
	if mode == "" {
		mode = runtimeContract
	}
	revision := options.SourceRevision
	if revision == "" {
		revision = "developer-trial-base"
	}
	catalog, err := forgecontext.NewRepositoryCatalog(root, revision, []forgecontext.RepositorySource{
		{
			ID: "accepted-architecture", Path: m1ArchitecturePath,
			Authority: "accepted-knowledge", Classification: "Internal",
		},
		{
			ID: "planned-target", Path: m1TargetPath,
			Authority: "reference", Classification: "Internal", Optional: true,
		},
	}, now)
	if err != nil {
		return operationalRuntime{}, err
	}
	result := operationalRuntime{
		root: root, runtimeRoot: runtimeRoot, sourceRevision: revision, catalog: catalog,
		worker: trialWorkerRunner{now: now}, verifier: trialVerifierRunner{now: now},
	}
	switch mode {
	case runtimeContract:
		if options.WorkerImage != "" || options.VerifierImage != "" {
			return operationalRuntime{}, errors.New("Docker images require --runtime docker")
		}
	case runtimeDocker:
		if options.SourceRevision == "" || options.WorkerImage == "" || options.VerifierImage == "" {
			return operationalRuntime{}, errors.New("Docker runtime requires exact source revision and immutable worker and verifier images")
		}
		if err := validateExactLocalRevision(root, options.SourceRevision, []string{
			m1ArchitecturePath, m1TargetPath,
		}); err != nil {
			return operationalRuntime{}, err
		}
		if err := os.MkdirAll(runtimeRoot, 0o700); err != nil {
			return operationalRuntime{}, err
		}
		worker, workerErr := execution.NewDockerRunner(execution.DockerConfig{
			Image: options.WorkerImage, TempRoot: runtimeRoot, Memory: "256m", CPUs: "1",
		}, nil)
		if workerErr != nil {
			return operationalRuntime{}, workerErr
		}
		verifier, verifierErr := verification.NewDockerRunner(verification.DockerConfig{
			Image: options.VerifierImage, TempRoot: runtimeRoot, Memory: "256m", CPUs: "1",
			PIDs: "8", MaxOutputBytes: 1 << 20,
		}, nil)
		if verifierErr != nil {
			return operationalRuntime{}, verifierErr
		}
		result.worker, result.verifier, result.operational = worker, verifier, true
	default:
		return operationalRuntime{}, errors.New("runtime must be contract or docker")
	}
	return result, nil
}

func validateExactLocalRevision(root, revision string, paths []string) error {
	if len(revision) != 40 || len(paths) == 0 {
		return errors.New("Docker runtime requires a full exact source commit")
	}
	for _, character := range revision {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return errors.New("Docker runtime requires a full exact source commit")
		}
	}
	headCommand := exec.Command("git", "rev-parse", "--verify", "HEAD^{commit}")
	headCommand.Dir = root
	headCommand.Env = []string{"PATH=" + os.Getenv("PATH")}
	head, err := headCommand.Output()
	if err != nil || strings.TrimSpace(string(head)) != revision {
		return errors.New("source revision does not match the exact checkout commit")
	}
	arguments := []string{"status", "--porcelain=v1", "--untracked-files=all", "--"}
	arguments = append(arguments, paths...)
	statusCommand := exec.Command("git", arguments...)
	statusCommand.Dir = root
	statusCommand.Env = []string{"PATH=" + os.Getenv("PATH")}
	status, err := statusCommand.Output()
	if err != nil || len(strings.TrimSpace(string(status))) != 0 {
		return errors.New("authorized source paths differ from the exact source commit")
	}
	return nil
}

func resolveRepositoryRoot(configured string) (string, error) {
	if configured != "" {
		root, err := filepath.Abs(configured)
		if err != nil {
			return "", err
		}
		return filepath.EvalSymlinks(root)
	}
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if info, statErr := os.Stat(filepath.Join(current, "go.mod")); statErr == nil && info.Mode().IsRegular() {
			return filepath.EvalSymlinks(current)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("Forge repository root was not found")
		}
		current = parent
	}
}

func withinRoot(root, target string) bool {
	relative, err := filepath.Rel(root, target)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func workerSources(sources []changecontract.SourceFile, writePaths []string) []changecontract.SourceFile {
	result := make([]changecontract.SourceFile, 0, len(writePaths))
	for _, source := range sources {
		if slices.Contains(writePaths, source.Path) {
			result = append(result, source)
		}
	}
	return result
}
