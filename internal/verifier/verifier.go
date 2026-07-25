// Package verifier implements Forge-owned independent candidate controls. It
// depends on shared contracts, never on the candidate-producing worker.
package verifier

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go/format"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

// Control identifiers are defined by the shared contract so the control plane
// can compose a profile without importing this deployment unit.
const (
	ControlCandidateIntegrity = verificationcontract.ControlCandidateIntegrity
	ControlDeclaredPathScope  = verificationcontract.ControlDeclaredPathScope
	ControlGoFormat           = verificationcontract.ControlGoFormat
	ControlDocumentMetadata   = verificationcontract.ControlDocumentMetadata
	ControlMarkdownLinks      = verificationcontract.ControlMarkdownLinks
)

var (
	ErrUnsupportedControl = errors.New("verifier control is unsupported")
	// ErrNoApplicableOracle reports that every independent consequential control
	// in the plan is incapable of failing for the paths this candidate changes.
	// Certifying such a candidate would be an empty assurance.
	ErrNoApplicableOracle = errors.New("no independent consequential control can fail for the changed paths")
	// ErrProducerIsVerifier reports that the principal that produced the
	// candidate is the principal asked to verify it.
	ErrProducerIsVerifier = errors.New("the candidate producer cannot verify its own candidate")
)

// metadataDirectories mirrors the canonical documentation directories whose
// documents must declare lifecycle status and owner.
var metadataDirectories = []string{
	"docs/foundation/", "docs/architecture/", "docs/governance/",
	"docs/decisions/", "docs/reference/", "docs/reference-applications/",
	"docs/context/",
}

var documentStatuses = []string{"Draft", "Proposed", "Accepted", "Superseded", "Deprecated", "Archive"}

// Applicable reports whether a control is capable of producing a failure for
// the given changed paths.
//
// The two integrity controls are deliberately never applicable. Both are
// already enforced before the verifier evaluates anything: the bundle is
// validated at entry, and the delivery lifecycle requires the assurance plan's
// changed paths to equal the candidate's own. Counting them as an oracle would
// certify a candidate against invariants that cannot fail.
func Applicable(controlID string, changedPaths []string) bool {
	switch controlID {
	case ControlCandidateIntegrity, ControlDeclaredPathScope:
		return false
	case ControlGoFormat:
		return anyPath(changedPaths, func(value string) bool { return path.Ext(value) == ".go" })
	case ControlDocumentMetadata:
		return anyPath(changedPaths, canonicalDocument)
	case ControlMarkdownLinks:
		return anyPath(changedPaths, func(value string) bool { return path.Ext(value) == ".md" })
	default:
		return false
	}
}

func anyPath(values []string, predicate func(string) bool) bool {
	for _, value := range values {
		if predicate(value) {
			return true
		}
	}
	return false
}

func canonicalDocument(value string) bool {
	if path.Ext(value) != ".md" {
		return false
	}
	for _, directory := range metadataDirectories {
		if strings.HasPrefix(value, directory) {
			return true
		}
	}
	return false
}

// DefaultIdentity names the deterministic Forge verifier. It is deliberately
// distinct from any candidate-producing principal so that the AC.8.1 invariant
// in the verification contract can be evaluated.
const DefaultIdentity = "forge.verifier.deterministic.v1"

type Engine struct {
	now      func() time.Time
	identity string
}

func NewEngine(now func() time.Time, identity string) (*Engine, error) {
	if now == nil {
		return nil, errors.New("verifier clock is required")
	}
	if identity == "" {
		return nil, errors.New("verifier identity is required")
	}
	return &Engine{now: now, identity: identity}, nil
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
	if plan.ProducedBy == engine.identity {
		return verificationcontract.Result{}, ErrProducerIsVerifier
	}
	if !hasApplicableOracle(plan) {
		return verificationcontract.Result{}, ErrNoApplicableOracle
	}
	findings := make([]verificationcontract.Finding, 0, len(plan.Controls))
	outcome := "passed"
	for _, control := range plan.Controls {
		if err := ctx.Err(); err != nil {
			return verificationcontract.Result{}, err
		}
		finding, err := evaluate(control, plan, bundle, workspace)
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
		Independent: true, VerifiedBy: engine.identity, CompletedAt: at,
		EvidenceCorrelationID: plan.EvidenceCorrelationID,
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

func hasApplicableOracle(plan verificationcontract.AssurancePlan) bool {
	for _, control := range plan.Controls {
		if control.Consequential && !control.CandidateOwned && Applicable(control.ID, plan.ChangedPaths) {
			return true
		}
	}
	return false
}

func evaluate(control verificationcontract.Control, plan verificationcontract.AssurancePlan, bundle workercontract.CandidateBundle, workspace string) (verificationcontract.Finding, error) {
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
	case ControlDocumentMetadata:
		passed = evaluateDocumentMetadata(bundle)
	case ControlMarkdownLinks:
		passed = evaluateMarkdownLinks(bundle, workspace)
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

// evaluateDocumentMetadata requires every canonical document the candidate
// touches to declare a valid lifecycle status and an owner. Authority in this
// repository is resolved by status, so a document without one has undefined
// authority.
func evaluateDocumentMetadata(bundle workercontract.CandidateBundle) bool {
	for _, file := range bundle.Files {
		if !canonicalDocument(file.Path) {
			continue
		}
		status, owner := "", ""
		for _, line := range strings.Split(string(file.Content), "\n") {
			line = strings.TrimRight(line, "\r")
			if value, found := strings.CutPrefix(line, "Status: "); found && status == "" {
				status = strings.TrimSpace(value)
			}
			if value, found := strings.CutPrefix(line, "Owner: "); found && owner == "" {
				owner = strings.TrimSpace(value)
			}
		}
		if owner == "" || !slices.Contains(documentStatuses, status) {
			return false
		}
	}
	return true
}

var markdownLink = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

// evaluateMarkdownLinks requires every relative link in a changed Markdown file
// to resolve, either against another candidate file or against the independently
// reconstructed source checkout.
func evaluateMarkdownLinks(bundle workercontract.CandidateBundle, workspace string) bool {
	candidatePaths := make(map[string]struct{}, len(bundle.Files))
	for _, file := range bundle.Files {
		candidatePaths[file.Path] = struct{}{}
	}
	root, err := filepath.Abs(workspace)
	if err != nil {
		return false
	}
	for _, file := range bundle.Files {
		if path.Ext(file.Path) != ".md" {
			continue
		}
		for _, match := range markdownLink.FindAllStringSubmatch(string(file.Content), -1) {
			target, ok := resolvableTarget(match[1])
			if !ok {
				continue
			}
			resolved := target
			if !strings.HasPrefix(target, "/") {
				resolved = path.Join(path.Dir(file.Path), target)
			}
			resolved = path.Clean(strings.TrimPrefix(resolved, "/"))
			if resolved == ".." || strings.HasPrefix(resolved, "../") {
				return false
			}
			if _, found := candidatePaths[resolved]; found {
				continue
			}
			if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(resolved))); statErr != nil {
				return false
			}
		}
	}
	return true
}

// resolvableTarget normalizes a link target and reports whether it names a
// repository path. Absolute URLs and pure fragments are not our concern.
func resolvableTarget(raw string) (string, bool) {
	target := strings.TrimSpace(raw)
	target = strings.TrimPrefix(strings.TrimSuffix(target, ">"), "<")
	if fields := strings.Fields(target); len(fields) > 0 {
		target = fields[0]
	}
	if target == "" || strings.HasPrefix(target, "#") {
		return "", false
	}
	if index := strings.IndexByte(target, '#'); index >= 0 {
		target = target[:index]
	}
	if target == "" {
		return "", false
	}
	if parsed, err := url.Parse(target); err != nil || parsed.Scheme != "" {
		return "", false
	}
	decoded, err := url.PathUnescape(target)
	if err != nil {
		return "", false
	}
	return decoded, true
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:])
}
