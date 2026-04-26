# Gin Auth Middleware Specification

## Purpose

Define a Gin middleware capability that authenticates JWT Bearer tokens using the security core validator contract, returns standardized auth errors, and injects validated claims into request context.

## Requirements

### Requirement: Middleware dependency boundary and factory contract

The system MUST expose `NewAuthMiddleware(validator security.Validator, opts ...Option) gin.HandlerFunc`. The middleware MUST depend on `security.Validator` and MUST NOT depend on service globals.

#### Scenario: Adapter uses security core interface only

- GIVEN a validator implementation conforming to `security.Validator`
- WHEN `NewAuthMiddleware` is created and executed
- THEN the request authentication flow uses only the passed validator instance
- AND no global singleton or service-level auth dependency is required

### Requirement: Token extraction precedence and configurable sources

The middleware MUST extract tokens from a configured header and cookie, with header value taking precedence over cookie value when both are present. The system SHALL support `WithHeaderName` and `WithCookieName` options.

#### Scenario: Valid token from header

- GIVEN a valid Bearer token in the configured auth header
- WHEN a protected route is called
- THEN the middleware authenticates using the header token

#### Scenario: Valid token from cookie fallback

- GIVEN no auth header and a valid token in the configured cookie
- WHEN a protected route is called
- THEN the middleware authenticates using the cookie token

#### Scenario: Header wins over cookie

- GIVEN both header and cookie tokens are present
- WHEN the request is authenticated
- THEN the middleware validates the header token
- AND ignores cookie token value for selection precedence

### Requirement: Auth enablement toggle

Authentication MUST be enabled by default. The system SHALL support `WithAuthEnabled(false)` to bypass token extraction and validation.

#### Scenario: Auth disabled passes through

- GIVEN middleware configured with `WithAuthEnabled(false)`
- WHEN any request reaches the middleware
- THEN the request continues to the next handler without auth checks

### Requirement: Unauthorized error response contract

When authentication fails, the middleware MUST return HTTP 401 with JSON payload `{ "code": string, "message": string }`. Missing token MUST use `code: auth_token_missing` and `message: "missing authentication token"`. Invalid, malformed, expired, invalid issuer, or invalid audience token outcomes MUST use `code: auth_token_invalid` and `message: "invalid or expired token"`.

#### Scenario: Missing token returns missing code and canonical message

- GIVEN neither configured header nor configured cookie contains a token
- WHEN the middleware processes the request with auth enabled
- THEN response status is 401
- AND response body contains `"code": "auth_token_missing"`
- AND response body contains `"message": "missing authentication token"`

#### Scenario: Malformed token returns invalid code and canonical message

- GIVEN a malformed token in the selected source
- WHEN validation is attempted
- THEN response status is 401
- AND response body contains `"code": "auth_token_invalid"`
- AND response body contains `"message": "invalid or expired token"`

#### Scenario: Expired token returns invalid code and canonical message

- GIVEN an expired token in the selected source
- WHEN validation is attempted
- THEN response status is 401
- AND response body contains `"code": "auth_token_invalid"`
- AND response body contains `"message": "invalid or expired token"`

#### Scenario: Invalid issuer or audience returns invalid code and canonical message

- GIVEN a token failing issuer or audience policy
- WHEN validation is attempted
- THEN response status is 401
- AND response body contains `"code": "auth_token_invalid"`
- AND response body contains `"message": "invalid or expired token"`

### Requirement: Exported message string constants

The package MUST export `MsgTokenMissing` and `MsgTokenInvalid` as public string constants so consumers can assert error messages without magic string literals.

| Constant | Value |
|---|---|
| `MsgTokenMissing` | `"missing authentication token"` |
| `MsgTokenInvalid` | `"invalid or expired token"` |

#### Scenario: Consumer uses exported constant for assertion

- GIVEN a consumer test that asserts error message on 401 response
- WHEN the consumer imports the `http/gin` package
- THEN `MsgTokenMissing` and `MsgTokenInvalid` are accessible as exported constants
- AND their values match the response body `"message"` field

### Requirement: Claims injection on successful authentication

On successful validation, the middleware MUST inject validated `claims.Claims` into `gin.Context` under a configurable key via `WithContextKey`, and request handling MUST continue.

#### Scenario: Claims available to downstream handlers

- GIVEN a valid token from header or cookie
- WHEN middleware validation succeeds
- THEN downstream handlers can retrieve `claims.Claims` from `gin.Context`
- AND the request reaches the next handler without unauthorized response

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
