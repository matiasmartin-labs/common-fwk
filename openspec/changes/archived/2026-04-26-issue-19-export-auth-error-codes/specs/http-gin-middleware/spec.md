# Delta for http/gin/middleware

## MODIFIED Requirements

### Requirement: Use Exported Error Codes from errors Package

The middleware MUST import `github.com/matiasmartin-labs/common-fwk/errors` (aliased as `fwkerrors`) and reference `fwkerrors.CodeTokenMissing` and `fwkerrors.CodeTokenInvalid` wherever the unexported `codeTokenMissing` and `codeTokenInvalid` constants were previously used.
(Previously: middleware defined and used two unexported string constants `codeTokenMissing` and `codeTokenInvalid` directly in `middleware.go`)

#### Scenario: Missing token returns correct error code

- GIVEN a request with no bearer token
- WHEN the auth middleware processes the request
- THEN the response body contains `"code": "auth_token_missing"`
- AND HTTP status is 401

#### Scenario: Invalid token returns correct error code

- GIVEN a request with a bearer token that fails validation
- WHEN the auth middleware processes the request
- THEN the response body contains `"code": "auth_token_invalid"`
- AND HTTP status is 401

#### Scenario: Existing middleware tests pass unchanged

- GIVEN `middleware_test.go` is not modified
- WHEN tests run
- THEN all tests pass (string values are identical — only source changed from literals to constants)
