package intelligence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
)

const maximumActivationFileBytes = 64 << 10

type activationFile struct {
	SchemaVersion          string           `json:"schemaVersion"`
	ID                     string           `json:"id"`
	Provider               string           `json:"provider"`
	Endpoint               string           `json:"endpoint"`
	AdapterVersion         string           `json:"adapterVersion"`
	AllowedModelIDs        []string         `json:"allowedModelIds"`
	AllowedPurposes        []string         `json:"allowedPurposes"`
	AllowedClassifications []string         `json:"allowedClassifications"`
	CredentialRef          string           `json:"credentialRef"`
	HumanApprovalRef       string           `json:"humanApprovalRef"`
	PromptVersion          string           `json:"promptVersion"`
	Retention              string           `json:"retention"`
	TrainingUse            string           `json:"trainingUse"`
	Region                 string           `json:"region"`
	StoreProviderState     bool             `json:"storeProviderState"`
	ReasoningEffort        string           `json:"reasoningEffort"`
	Pricing                activationPrice  `json:"pricing"`
	Budgets                activationBudget `json:"budgets"`
	ActivatedAt            time.Time        `json:"activatedAt"`
	ExpiresAt              time.Time        `json:"expiresAt"`
}

type activationPrice struct {
	Currency                     string `json:"currency"`
	InputMicrosPerMillionTokens  int64  `json:"inputMicrosPerMillionTokens"`
	OutputMicrosPerMillionTokens int64  `json:"outputMicrosPerMillionTokens"`
}

type activationBudget struct {
	MaxInputBytes          int64 `json:"maxInputBytes"`
	MaxOutputBytes         int64 `json:"maxOutputBytes"`
	MaxTokens              int64 `json:"maxTokens"`
	MaxLatencyMs           int64 `json:"maxLatencyMs"`
	MaxCostMicros          int64 `json:"maxCostMicros"`
	MaxTurns               int   `json:"maxTurns"`
	MaxSessionTokens       int64 `json:"maxSessionTokens"`
	MaxSessionCostMicros   int64 `json:"maxSessionCostMicros"`
	MaxProjectTokensPerDay int64 `json:"maxProjectTokensPerDay"`
	MaxProjectCostPerDay   int64 `json:"maxProjectCostMicrosPerDay"`
	MaxConcurrentRequests  int   `json:"maxConcurrentRequests"`
}

// LoadProviderActivation reads one exact, non-secret provider activation.
// The returned digest binds runtime receipts to the exact bytes approved by a
// human. Credentials remain references and are resolved only at invocation.
func LoadProviderActivation(path string, now time.Time) (ProviderActivation, error) {
	if path == "" || now.IsZero() {
		return ProviderActivation{}, errors.New("provider activation file request is invalid")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return ProviderActivation{}, errors.New("read provider activation file")
	}
	if !info.Mode().IsRegular() {
		return ProviderActivation{}, errors.New("provider activation must be a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return ProviderActivation{}, errors.New("open provider activation file")
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return ProviderActivation{}, errors.New("provider activation file changed during open")
	}
	raw, err := io.ReadAll(io.LimitReader(file, maximumActivationFileBytes+1))
	if err != nil {
		return ProviderActivation{}, errors.New("read provider activation file")
	}
	if len(raw) == 0 || len(raw) > maximumActivationFileBytes {
		return ProviderActivation{}, errors.New("provider activation file size is invalid")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var encoded activationFile
	if err := decoder.Decode(&encoded); err != nil {
		return ProviderActivation{}, errors.New("decode provider activation file")
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF {
		return ProviderActivation{}, errors.New("provider activation file must contain one object")
	}
	sum := sha256.Sum256(raw)
	value := ProviderActivation{
		SchemaVersion:          encoded.SchemaVersion,
		ID:                     encoded.ID,
		Digest:                 "sha256:" + hex.EncodeToString(sum[:]),
		Provider:               encoded.Provider,
		Endpoint:               encoded.Endpoint,
		AdapterVersion:         encoded.AdapterVersion,
		AllowedModelIDs:        encoded.AllowedModelIDs,
		AllowedPurposes:        encoded.AllowedPurposes,
		AllowedClassifications: encoded.AllowedClassifications,
		CredentialRef:          encoded.CredentialRef,
		HumanApprovalRef:       encoded.HumanApprovalRef,
		PromptVersion:          encoded.PromptVersion,
		Retention:              encoded.Retention,
		TrainingUse:            encoded.TrainingUse,
		Region:                 encoded.Region,
		StoreProviderState:     encoded.StoreProviderState,
		ReasoningEffort:        encoded.ReasoningEffort,
		Pricing: Pricing{
			Currency:                     encoded.Pricing.Currency,
			InputMicrosPerMillionTokens:  encoded.Pricing.InputMicrosPerMillionTokens,
			OutputMicrosPerMillionTokens: encoded.Pricing.OutputMicrosPerMillionTokens,
		},
		Budgets: discovery.Limits{
			MaxInputBytes:  encoded.Budgets.MaxInputBytes,
			MaxOutputBytes: encoded.Budgets.MaxOutputBytes,
			MaxTokens:      encoded.Budgets.MaxTokens,
			MaxLatency:     time.Duration(encoded.Budgets.MaxLatencyMs) * time.Millisecond,
			MaxCostMicros:  encoded.Budgets.MaxCostMicros,
		},
		MaxTurns:               encoded.Budgets.MaxTurns,
		MaxSessionTokens:       encoded.Budgets.MaxSessionTokens,
		MaxSessionCostMicros:   encoded.Budgets.MaxSessionCostMicros,
		MaxProjectTokensPerDay: encoded.Budgets.MaxProjectTokensPerDay,
		MaxProjectCostPerDay:   encoded.Budgets.MaxProjectCostPerDay,
		MaxConcurrentRequests:  encoded.Budgets.MaxConcurrentRequests,
		ActivatedAt:            encoded.ActivatedAt,
		ExpiresAt:              encoded.ExpiresAt,
	}
	if err := value.validate(now); err != nil {
		return ProviderActivation{}, fmt.Errorf("provider activation is ineligible: %w", err)
	}
	return cloneActivation(value), nil
}
