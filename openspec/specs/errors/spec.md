# errors — Auth Error Codes

## Purpose

Stable, exported string constants for all auth error codes and HTTP routing error codes, consumable by any package inside or outside this module.

## Requirements

### Requirement: Export Auth Error Code Constants

The `errors` package MUST export string constants covering all auth error conditions AND HTTP routing error conditions.

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
| `CodeNotFound` | `"not_found"` |
| `CodeMethodNotAllowed` | `"method_not_allowed"` |

#### Scenario: Constant string stability

- GIVEN the `errors` package is imported
- WHEN any of the 11 constants is read
- THEN its value MUST equal the exact string in the table above

#### Scenario: All constants accessible at package level

- GIVEN the module is compiled
- WHEN a consumer references any of the 11 constants
- THEN compilation succeeds without additional imports

### Requirement: Test Coverage for String Stability

The `errors` package MUST include a table-driven test file `codes_test.go` that verifies the string value of all 11 constants.

#### Scenario: Table-driven test covers all 11 constants

- GIVEN `errors/codes_test.go` exists
- WHEN the test runs
- THEN each of the 11 constants is asserted against its expected string value
- AND the test passes with no failures
