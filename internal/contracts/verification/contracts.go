// Package verification defines immutable contracts shared with the isolated
// verifier. It has no dependency on candidate-producing runtime code.
package verification

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"sort"
	"strings"
	"time"
)

// ADR-0025 versioned these contracts when provenance became mandatory. v1
// remains readable for historical validation but cannot satisfy a current gate.
const (
	AssurancePlanVersion   = "forge.assurance-plan.v2"
	ResultVersion          = "forge.verification-result.v2"
	AssurancePlanVersionV1 = "forge.assurance-plan.v1"
	ResultVersionV1        = "forge.verification-result.v1"
)

var ErrInvalidContract = errors.New("verification contract is invalid")

func DecodePlan(encoded []byte) (AssurancePlan, error) {
	var value AssurancePlan
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return AssurancePlan{}, ErrInvalidContract
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return AssurancePlan{}, ErrInvalidContract
	}
	return value, nil
}

func DecodeResult(encoded []byte) (Result, error) {
	var value Result
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return Result{}, ErrInvalidContract
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Result{}, ErrInvalidContract
	}
	return value, nil
}

// Forge-owned control identifiers. They live in the shared contract because
// both the control plane, when it composes an assurance profile, and the
// isolated verifier, when it evaluates one, must name the same controls without
// the control plane depending on the verifier deployment unit.
const (
	ControlCandidateIntegrity = "forge.candidate-integrity"
	ControlDeclaredPathScope  = "forge.declared-path-scope"
	ControlGoFormat           = "forge.go-format"
	ControlDocumentMetadata   = "forge.document-metadata"
	ControlMarkdownLinks      = "forge.markdown-links"
)

type Control struct {
	ID             string `json:"id"`
	Version        string `json:"version"`
	Consequential  bool   `json:"consequential"`
	CandidateOwned bool   `json:"candidateOwned"`
}

type AssurancePlan struct {
	SchemaVersion         string    `json:"schemaVersion"`
	ID                    string    `json:"id"`
	Digest                string    `json:"digest"`
	CandidateDigest       string    `json:"candidateDigest"`
	PlanDigest            string    `json:"planDigest"`
	RiskClass             string    `json:"riskClass"`
	Classification        string    `json:"classification"`
	AffectedCapabilities  []string  `json:"affectedCapabilities"`
	ChangedPaths          []string  `json:"changedPaths"`
	Controls              []Control `json:"controls"`
	ProducedBy            string    `json:"producedBy"`
	Environment           string    `json:"environment"`
	ExpiresAt             time.Time `json:"expiresAt"`
	EvidenceCorrelationID string    `json:"evidenceCorrelationId"`
}

type Finding struct {
	ControlID      string `json:"controlId"`
	Outcome        string `json:"outcome"`
	Code           string `json:"code"`
	EvidenceDigest string `json:"evidenceDigest"`
	CandidateOwned bool   `json:"candidateOwned"`
	Consequential  bool   `json:"consequential"`
}

type Result struct {
	SchemaVersion         string    `json:"schemaVersion"`
	ID                    string    `json:"id"`
	Digest                string    `json:"digest"`
	AssurancePlanID       string    `json:"assurancePlanId"`
	AssurancePlanDigest   string    `json:"assurancePlanDigest"`
	CandidateDigest       string    `json:"candidateDigest"`
	Environment           string    `json:"environment"`
	Findings              []Finding `json:"findings"`
	Outcome               string    `json:"outcome"`
	Independent           bool      `json:"independent"`
	VerifiedBy            string    `json:"verifiedBy"`
	CompletedAt           time.Time `json:"completedAt"`
	EvidenceCorrelationID string    `json:"evidenceCorrelationId"`
}

func BindPlan(value AssurancePlan) (AssurancePlan, error) {
	value.ID, value.Digest = "", ""
	sort.Strings(value.AffectedCapabilities)
	value.AffectedCapabilities = slices.Compact(value.AffectedCapabilities)
	sort.Strings(value.ChangedPaths)
	value.ChangedPaths = slices.Compact(value.ChangedPaths)
	sort.Slice(value.Controls, func(i, j int) bool { return value.Controls[i].ID < value.Controls[j].ID })
	if value.SchemaVersion != AssurancePlanVersion || !validDigest(value.CandidateDigest) ||
		!validDigest(value.PlanDigest) || value.RiskClass == "" || value.Classification == "" ||
		value.Environment != "developer-trial" || value.ExpiresAt.IsZero() ||
		value.EvidenceCorrelationID == "" || len(value.AffectedCapabilities) == 0 ||
		len(value.ChangedPaths) == 0 || len(value.Controls) == 0 || value.ProducedBy == "" {
		return AssurancePlan{}, ErrInvalidContract
	}
	hasIndependentOracle := false
	for index, control := range value.Controls {
		if control.ID == "" || control.Version == "" ||
			(index > 0 && value.Controls[index-1].ID == control.ID) {
			return AssurancePlan{}, ErrInvalidContract
		}
		if control.Consequential && !control.CandidateOwned {
			hasIndependentOracle = true
		}
	}
	if !hasIndependentOracle {
		return AssurancePlan{}, ErrInvalidContract
	}
	value.Digest = digest(value)
	value.ID = "forge.assurance-plan." + strings.TrimPrefix(value.Digest, "sha256:")
	return value, nil
}

func (value AssurancePlan) Validate(at time.Time) error {
	bound, err := BindPlan(value)
	if err != nil || bound.ID != value.ID || bound.Digest != value.Digest ||
		at.IsZero() || !at.Before(value.ExpiresAt) {
		return ErrInvalidContract
	}
	return nil
}

func BindResult(value Result) (Result, error) {
	value.ID, value.Digest = "", ""
	sort.Slice(value.Findings, func(i, j int) bool { return value.Findings[i].ControlID < value.Findings[j].ControlID })
	if value.SchemaVersion != ResultVersion || value.AssurancePlanID == "" ||
		!validDigest(value.AssurancePlanDigest) || !validDigest(value.CandidateDigest) ||
		value.Environment != "developer-trial" || value.CompletedAt.IsZero() ||
		value.EvidenceCorrelationID == "" || !value.Independent || value.VerifiedBy == "" ||
		(value.Outcome != "passed" && value.Outcome != "failed") || len(value.Findings) == 0 {
		return Result{}, ErrInvalidContract
	}
	independentPass := false
	consequentialFailed := false
	for index, finding := range value.Findings {
		if finding.ControlID == "" || finding.Code == "" || !validDigest(finding.EvidenceDigest) ||
			(finding.Outcome != "passed" && finding.Outcome != "failed") ||
			(index > 0 && value.Findings[index-1].ControlID == finding.ControlID) {
			return Result{}, ErrInvalidContract
		}
		if finding.Consequential && !finding.CandidateOwned && finding.Outcome == "passed" {
			independentPass = true
		}
		if finding.Consequential && finding.Outcome == "failed" {
			consequentialFailed = true
		}
	}
	if value.Outcome == "passed" && (!independentPass || consequentialFailed) {
		return Result{}, ErrInvalidContract
	}
	value.Digest = digest(value)
	value.ID = "forge.verification-result." + strings.TrimPrefix(value.Digest, "sha256:")
	return value, nil
}

func (value Result) Validate(plan AssurancePlan) error {
	bound, err := BindResult(value)
	if err != nil || bound.ID != value.ID || bound.Digest != value.Digest ||
		value.AssurancePlanID != plan.ID || value.AssurancePlanDigest != plan.Digest ||
		value.CandidateDigest != plan.CandidateDigest || value.Environment != plan.Environment ||
		value.EvidenceCorrelationID != plan.EvidenceCorrelationID ||
		!value.CompletedAt.Before(plan.ExpiresAt) || len(value.Findings) != len(plan.Controls) {
		return ErrInvalidContract
	}
	// OWASP AISVS Appendix C, AC.8.1: an agent may not approve what it produced.
	// The producing principal is issued by Identity under Policy, so a candidate
	// cannot assert it. Today's verifier is deterministic and therefore always a
	// different principal, so this invariant is dormant; it becomes load-bearing
	// the moment verification itself becomes model-assisted. It lives here
	// rather than in a control because a profile can omit a control.
	if value.VerifiedBy == plan.ProducedBy {
		return ErrInvalidContract
	}
	for index, control := range plan.Controls {
		finding := value.Findings[index]
		if finding.ControlID != control.ID || finding.CandidateOwned != control.CandidateOwned ||
			finding.Consequential != control.Consequential {
			return ErrInvalidContract
		}
	}
	return nil
}

func digest(value any) string {
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil
}
