package errors

// Auth error codes returned in JSON error responses.
// Consumers can compare against these constants instead of raw strings.
const (
	CodeTokenMissing          = "auth_token_missing"
	CodeTokenInvalid          = "auth_token_invalid"
	CodeCallbackStateInvalid  = "auth_callback_state_invalid"
	CodeCallbackCodeMissing   = "auth_callback_code_missing"
	CodeEmailNotAllowed       = "auth_email_not_allowed"
	CodeProviderFailure       = "auth_provider_failure"
	CodeTokenGenerationFailed = "auth_token_generation_failed"
	CodeClaimsMissing         = "auth_claims_missing"
	CodeClaimsInvalid         = "auth_claims_invalid"
)
