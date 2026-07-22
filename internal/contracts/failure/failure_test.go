package failure

import "testing"

func TestNewProducesValidDisclosureSafeError(t *testing.T) {
	t.Parallel()

	got := New(Forbidden, "action is not authorized", "correlation-1", RetryNone)
	if err := got.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got.Error() != "forbidden: action is not authorized" {
		t.Fatalf("Error() = %q", got.Error())
	}
}

func TestNewUsesNonEmptyFallbackCorrelation(t *testing.T) {
	t.Parallel()

	got := New(InvalidRequest, "request is invalid", "", RetryNone)
	if got.CorrelationID != "forge.validation.uncorrelated" {
		t.Fatalf("CorrelationID = %q", got.CorrelationID)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsUnknownValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Error)
	}{
		{"schema", func(value *Error) { value.SchemaVersion = "other" }},
		{"code", func(value *Error) { value.Code = "other" }},
		{"message", func(value *Error) { value.SafeMessage = "" }},
		{"correlation", func(value *Error) { value.CorrelationID = "" }},
		{"retry", func(value *Error) { value.Retry = "other" }},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			value := New(InvalidRequest, "request is invalid", "correlation-1", RetryNone)
			test.mutate(value)
			if err := value.Validate(); err == nil {
				t.Fatal("Validate() error = nil")
			}
		})
	}
}
