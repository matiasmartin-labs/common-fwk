package errors_test

import (
	"testing"

	fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"
)

func TestAuthErrorCodes(t *testing.T) {
	tests := []struct {
		name  string
		got   string
		want  string
	}{
		{"CodeTokenMissing", fwkerrors.CodeTokenMissing, "auth_token_missing"},
		{"CodeTokenInvalid", fwkerrors.CodeTokenInvalid, "auth_token_invalid"},
		{"CodeCallbackStateInvalid", fwkerrors.CodeCallbackStateInvalid, "auth_callback_state_invalid"},
		{"CodeCallbackCodeMissing", fwkerrors.CodeCallbackCodeMissing, "auth_callback_code_missing"},
		{"CodeEmailNotAllowed", fwkerrors.CodeEmailNotAllowed, "auth_email_not_allowed"},
		{"CodeProviderFailure", fwkerrors.CodeProviderFailure, "auth_provider_failure"},
		{"CodeTokenGenerationFailed", fwkerrors.CodeTokenGenerationFailed, "auth_token_generation_failed"},
		{"CodeClaimsMissing", fwkerrors.CodeClaimsMissing, "auth_claims_missing"},
		{"CodeClaimsInvalid", fwkerrors.CodeClaimsInvalid, "auth_claims_invalid"},
		{name: "not_found", got: fwkerrors.CodeNotFound, want: "not_found"},
		{name: "method_not_allowed", got: fwkerrors.CodeMethodNotAllowed, want: "method_not_allowed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
			}
		})
	}
}
