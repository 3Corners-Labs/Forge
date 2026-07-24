package intelligence_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/intelligence"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

type tokenSourceStub struct {
	calls int
	token string
}

func (source *tokenSourceStub) Token(context.Context, string) (string, error) {
	source.calls++
	if source.token == "" {
		return "", errors.New("missing token")
	}
	return source.token, nil
}

func TestOpenAIAdapterDisablesStorageAndParsesStructuredOutput(t *testing.T) {
	request := requestFixture()
	modelOutput, err := json.Marshal(map[string]any{
		"response": "Who chooses the desk?", "draft": emptyModelDraft(),
		"questions": []map[string]any{{
			"id": "q-1", "text": "Who chooses the desk?",
			"consequence": "The primary actor remains unclear", "material": true,
		}},
		"uncertainty": []string{"primary actor"}, "limitations": []string{"no repository context"},
		"provenanceRefs": []string{"interaction:user:1"}, "refusal": false, "incomplete": true,
	})
	if err != nil {
		t.Fatal(err)
	}
	tokens := &tokenSourceStub{token: "FORGE_TEST_SECRET_DO_NOT_RECORD"}
	client := &http.Client{Transport: roundTripFunc(func(httpRequest *http.Request) (*http.Response, error) {
		if httpRequest.URL.String() != "https://eu.api.openai.com/v1/responses" ||
			httpRequest.Header.Get("Authorization") != "Bearer FORGE_TEST_SECRET_DO_NOT_RECORD" {
			t.Fatalf("unexpected provider request: %s headers=%v", httpRequest.URL, httpRequest.Header)
		}
		body, err := io.ReadAll(httpRequest.Body)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(body, []byte("FORGE_TEST_SECRET_DO_NOT_RECORD")) {
			t.Fatal("credential leaked into provider request body")
		}
		var value map[string]any
		if err := json.Unmarshal(body, &value); err != nil {
			t.Fatal(err)
		}
		if stored, ok := value["store"].(bool); !ok || stored {
			t.Fatalf("provider storage was not disabled: %s", body)
		}
		if maximum, ok := value["max_output_tokens"].(float64); !ok ||
			int64(maximum) != request.Limits.MaxTokens {
			t.Fatalf("provider output token budget was not applied: %s", body)
		}
		reasoning, ok := value["reasoning"].(map[string]any)
		if !ok || reasoning["effort"] != "medium" {
			t.Fatalf("approved reasoning profile was not used: %s", body)
		}
		text, ok := value["text"].(map[string]any)
		if !ok {
			t.Fatalf("structured output configuration missing: %s", body)
		}
		format, ok := text["format"].(map[string]any)
		if !ok || format["type"] != "json_schema" || format["strict"] != true {
			t.Fatalf("strict schema was not requested: %s", body)
		}
		response := map[string]any{
			"id": "resp-1", "model": "model-approved",
			"output": []map[string]any{{"content": []map[string]any{{
				"type": "output_text", "text": string(modelOutput),
			}}}},
			"usage": map[string]any{"input_tokens": 100, "output_tokens": 50},
		}
		encoded, _ := json.Marshal(response)
		return &http.Response{
			StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(encoded)),
			Header: make(http.Header),
		}, nil
	})}
	adapter, err := intelligence.NewOpenAIAdapter(client, tokens)
	if err != nil {
		t.Fatal(err)
	}
	activation := openAIActivationFixture()
	activation.Endpoint = "https://eu.api.openai.com/v1/responses"
	activation.Region = "eu"
	proposal, err := adapter.Propose(context.Background(), activation, request)
	if err != nil {
		t.Fatal(err)
	}
	if tokens.calls != 1 || proposal.Receipt.ProviderRequest != "resp-1" ||
		proposal.Receipt.ModelID != "model-approved" || proposal.Questions[0].ID != "q-1" ||
		proposal.Receipt.ActivationID != activation.ID ||
		proposal.Receipt.ActivationDigest != activation.Digest ||
		proposal.Usage.CostMicros != 1000 {
		t.Fatalf("unexpected proposal: %#v", proposal)
	}
}

func TestOpenAIAdapterRejectsUnapprovedEndpointBeforeCredentialResolution(t *testing.T) {
	tokens := &tokenSourceStub{token: "secret"}
	adapter, err := intelligence.NewOpenAIAdapter(&http.Client{}, tokens)
	if err != nil {
		t.Fatal(err)
	}
	activation := openAIActivationFixture()
	activation.Endpoint = "https://attacker.invalid/v1/responses"
	_, err = adapter.Propose(context.Background(), activation, requestFixture())
	if err == nil || tokens.calls != 0 {
		t.Fatalf("unapproved endpoint reached credential source: calls=%d err=%v", tokens.calls, err)
	}
}

func TestOpenAIAdapterEnforcesWireInputAndOutputBudgets(t *testing.T) {
	request := requestFixture()
	tokens := &tokenSourceStub{token: "secret"}
	networkCalls := 0
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		networkCalls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(strings.Repeat("x", 1024))),
			Header:     make(http.Header),
		}, nil
	})}
	adapter, err := intelligence.NewOpenAIAdapter(client, tokens)
	if err != nil {
		t.Fatal(err)
	}
	activation := openAIActivationFixture()
	request.Limits.MaxInputBytes = 1
	if _, err := adapter.Propose(context.Background(), activation, request); err == nil ||
		networkCalls != 0 || tokens.calls != 0 {
		t.Fatalf("oversized wire input reached credential or provider: credentials=%d calls=%d err=%v",
			tokens.calls, networkCalls, err)
	}

	request = requestFixture()
	request.Limits.MaxOutputBytes = 64
	if _, err := adapter.Propose(context.Background(), activation, request); err == nil ||
		networkCalls != 1 || tokens.calls != 1 {
		t.Fatalf("oversized wire output was accepted: credentials=%d calls=%d err=%v",
			tokens.calls, networkCalls, err)
	}
}

func emptyModelDraft() map[string]any {
	return map[string]any{
		"outcome": "", "affectedUsers": []string{}, "problem": "",
		"functionalBehavior": []string{}, "uxExpectations": []string{},
		"accessibility": []string{}, "privacy": []string{}, "security": []string{},
		"operations": []string{}, "constraints": []string{}, "exclusions": []string{},
		"acceptanceCriteria": []string{}, "inferences": []string{}, "assumptions": []string{},
		"hypotheses": []string{}, "experiments": []string{}, "productAreaRefs": []string{},
	}
}

func TestEnvironmentTokenSourceAcceptsOnlyNamedEnvironmentReference(t *testing.T) {
	t.Setenv("FORGE_DISCOVERY_TEST_TOKEN", "value")
	source := intelligence.EnvironmentTokenSource{}
	token, err := source.Token(context.Background(), "env:FORGE_DISCOVERY_TEST_TOKEN")
	if err != nil || token != "value" {
		t.Fatalf("token=%q err=%v", token, err)
	}
	for _, reference := range []string{"FORGE_DISCOVERY_TEST_TOKEN", "env:", "env:BAD NAME"} {
		if _, err := source.Token(context.Background(), reference); err == nil {
			t.Fatalf("invalid reference %q was accepted", reference)
		}
	}
}

func TestOpenAIProviderErrorDoesNotIncludeResponseBody(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader("FORGE_TEST_SECRET_DO_NOT_RECORD")),
			Header:     make(http.Header),
		}, nil
	})}
	tokens := &tokenSourceStub{token: "secret"}
	adapter, err := intelligence.NewOpenAIAdapter(client, tokens)
	if err != nil {
		t.Fatal(err)
	}
	activation := openAIActivationFixture()
	_, err = adapter.Propose(context.Background(), activation, requestFixture())
	if err == nil || strings.Contains(err.Error(), "FORGE_TEST_SECRET_DO_NOT_RECORD") {
		t.Fatalf("provider response body leaked: %v", err)
	}
}

func TestOpenAIAdapterDoesNotFollowRedirects(t *testing.T) {
	requests := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests++
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     http.Header{"Location": []string{"https://attacker.invalid/capture"}},
			Request:    request,
		}, nil
	})}
	tokens := &tokenSourceStub{token: "secret"}
	adapter, err := intelligence.NewOpenAIAdapter(client, tokens)
	if err != nil {
		t.Fatal(err)
	}
	activation := openAIActivationFixture()
	_, err = adapter.Propose(context.Background(), activation, requestFixture())
	if err == nil || requests != 1 || tokens.calls != 1 {
		t.Fatalf("redirect was not stopped: requests=%d tokenCalls=%d err=%v", requests, tokens.calls, err)
	}
}

func openAIActivationFixture() intelligence.ProviderActivation {
	return intelligence.ProviderActivation{
		SchemaVersion:          intelligence.ProviderActivationVersion,
		ID:                     "activation-openai-1",
		Digest:                 "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Provider:               "openai",
		Endpoint:               "https://api.openai.com/v1/responses",
		AdapterVersion:         "forge.openai-discovery-adapter.v1",
		AllowedModelIDs:        []string{"model-approved"},
		AllowedPurposes:        []string{discovery.PurposeV1},
		AllowedClassifications: []string{"Internal"},
		CredentialRef:          "env:OPENAI_API_KEY",
		HumanApprovalRef:       "approval:provider-1",
		PromptVersion:          "forge.discovery.prompt.v1",
		Retention:              "provider-default-up-to-30-days",
		TrainingUse:            "disabled",
		Region:                 "global",
		ReasoningEffort:        "medium",
		Pricing: intelligence.Pricing{
			Currency:                     "USD",
			InputMicrosPerMillionTokens:  2_500_000,
			OutputMicrosPerMillionTokens: 15_000_000,
		},
		Budgets: discovery.Limits{
			MaxInputBytes: 16 << 10, MaxOutputBytes: 16 << 10, MaxTokens: 4096,
			MaxLatency: 10 * time.Second, MaxCostMicros: 100_000,
		},
		MaxTurns:         8,
		MaxSessionTokens: 32768, MaxSessionCostMicros: 1_000_000,
		MaxProjectTokensPerDay: 131072, MaxProjectCostPerDay: 4_000_000,
		MaxConcurrentRequests: 1,
		ActivatedAt:           intelligenceTime.Add(-time.Minute),
		ExpiresAt:             intelligenceTime.Add(time.Hour),
	}
}
