// Package verifier implements Forge-owned independent candidate controls. It
// depends on shared contracts, never on the candidate-producing worker.
package verifier

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

const (
	ControlCandidateIntegrity = "forge.candidate-integrity"
	ControlDeclaredPathScope  = "forge.declared-path-scope"
	ControlGoFormat           = "forge.go-format"
)

var ErrUnsupportedControl = errors.New("verifier control is unsupported")

type Engine struct {
	now func() time.Time
}

func NewEngine(now func() time.Time) (*Engine, error) {
	if now == nil {
		return nil, errors.New("verifier clock is required")
	}
	return &Engine{now: now}, nil
}

func (engine *Engine) Verify(ctx context.Context, plan verificationcontract.AssurancePlan, bundle workercontract.CandidateBundle, workspace string) (verificationcontract.Result, error) {
	at := engine.now()
	if ctx == nil || at.IsZero() {
		return verificationcontract.Result{}, verificationcontract.ErrInvalidContract
	}
	if err := plan.Validate(at); err != nil {
		return verificationcontract.Result{}, err
	}
	if workspace == "" {
		return verificationcontract.Result{}, verificationcontract.ErrInvalidContract
	}
	if err := bundle.Validate(); err != nil || bundle.Digest != plan.CandidateDigest || bundle.PlanDigest != plan.PlanDigest {
		return verificationcontract.Result{}, verificationcontract.ErrInvalidContract
	}
	if err := verifyFreshSource(workspace, bundle); err != nil {
		return verificationcontract.Result{}, verificationcontract.ErrInvalidContract
	}
	findings := make([]verificationcontract.Finding, 0, len(plan.Controls))
	outcome := "passed"
	for _, control := range plan.Controls {
		if err := ctx.Err(); err != nil {
			return verificationcontract.Result{}, err
		}
		finding, err := evaluate(control, plan, bundle)
		if err != nil {
			return verificationcontract.Result{}, err
		}
		if finding.Consequential && finding.Outcome != "passed" {
			outcome = "failed"
		}
		findings = append(findings, finding)
	}
	return verificationcontract.BindResult(verificationcontract.Result{
		SchemaVersion: verificationcontract.ResultVersion, AssurancePlanID: plan.ID,
		AssurancePlanDigest: plan.Digest, CandidateDigest: bundle.Digest,
		Environment: plan.Environment, Findings: findings, Outcome: outcome,
		Independent: true, CompletedAt: at, EvidenceCorrelationID: plan.EvidenceCorrelationID,
	})
}

func verifyFreshSource(workspace string, bundle workercontract.CandidateBundle) error {
	root, err := filepath.Abs(workspace)
	if err != nil {
		return err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	for _, file := range bundle.Files {
		target := filepath.Join(root, filepath.FromSlash(file.Path))
		relative, err := filepath.Rel(root, target)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return verificationcontract.ErrInvalidContract
		}
		if info, statErr := os.Lstat(target); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return verificationcontract.ErrInvalidContract
		}
		current, readErr := os.ReadFile(target)
		if file.SourceDigest == "" {
			if readErr == nil || !errors.Is(readErr, os.ErrNotExist) {
				return verificationcontract.ErrInvalidContract
			}
			continue
		}
		if readErr != nil || digest(string(current)) != file.SourceDigest {
			return verificationcontract.ErrInvalidContract
		}
	}
	return nil
}

func evaluate(control verificationcontract.Control, plan verificationcontract.AssurancePlan, bundle workercontract.CandidateBundle) (verificationcontract.Finding, error) {
	passed := false
	code := "failed"
	switch control.ID {
	case ControlCandidateIntegrity:
		passed = bundle.Validate() == nil
	case ControlDeclaredPathScope:
		passed = true
		for _, file := range bundle.Files {
			if !slices.Contains(plan.ChangedPaths, file.Path) {
				passed = false
				break
			}
		}
	case ControlGoFormat:
		passed = true
		for _, file := range bundle.Files {
			if path.Ext(file.Path) != ".go" {
				continue
			}
			formatted, err := format.Source(file.Content)
			if err != nil || string(formatted) != string(file.Content) {
				passed = false
				break
			}
		}
	default:
		if !control.CandidateOwned {
			return verificationcontract.Finding{}, ErrUnsupportedControl
		}
		// Candidate-owned controls can be recorded but never become an
		// independent consequential oracle.
		passed = true
	}
	if passed {
		code = "passed"
	}
	evidence := digest(strings.Join([]string{control.ID, control.Version, code, bundle.Digest}, "\x00"))
	return verificationcontract.Finding{
		ControlID: control.ID, Outcome: code, Code: code, EvidenceDigest: evidence,
		CandidateOwned: control.CandidateOwned, Consequential: control.Consequential,
	}, nil
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:])
}
