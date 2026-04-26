# errors — Auth Error Codes

## Purpose

Stable, exported string constants for all auth error codes, consumable by any package inside or outside this module.

## Requirements

### Requirement: Export Auth Error Code Constants

The `errors` package MUST export 9 string constants covering all auth error conditions.

| Constant | String Value |
|---|---|
| `CodeTokenMissing` | `"auth_token_missing"` |
| `CodeTokenInvalid` | `"auth_token_invalid"` |
| `CodeCallbackStateInvalid` | `"auth_callback_state_invalid"` |
| `CodeCallbackCodeMissing` | `"auth_callback_code_missing"` |
| `CodeEmailNotAllowed` | `"auth_email_not_allowed"` |
| `CodeProviderFailure` | `"auth_provider_failure"` |
| `CodeTokenGenerationFailed` | `"auth_token_generation_failed"` |
| `CodeClaimsMissing` | `"auth_claims_missing"` |
| `CodeClaimsInvalid` | `"auth_claims_invalid"` |

#### Scenario: Constant string stability

- GIVEN the `errors` package is imported
- WHEN any of the 9 constants is read
- THEN its value MUST equal the exact string in the table above

#### Scenario: All constants accessible at package level

- GIVEN the module is compiled
- WHEN a consumer references `errors.CodeTokenMissing` (or any of the 9)
- THEN compilation succeeds without additional imports

### Requirement: Test Coverage for String Stability

The `errors` package MUST include a table-driven test file `codes_test.go` that verifies the string value of all 9 constants.

#### Scenario: Table-driven test covers all 9 constants

- GIVEN `errors/codes_test.go` exists
- WHEN the test runs
- THEN each of the 9 constants is asserted against its expected string value
- AND the test passes with no failures
