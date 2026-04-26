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

When authentication fails, the middleware MUST return HTTP 401 with JSON payload `{ "code": string, "message": string }`. Missing token MUST use `auth_token_missing`. Invalid, malformed, expired, invalid issuer, or invalid audience token outcomes MUST use `auth_token_invalid`.

#### Scenario: Missing token returns missing code

- GIVEN neither configured header nor configured cookie contains a token
- WHEN the middleware processes the request with auth enabled
- THEN response status is 401 with code `auth_token_missing`

#### Scenario: Malformed token returns invalid code

- GIVEN a malformed token in the selected source
- WHEN validation is attempted
- THEN response status is 401 with code `auth_token_invalid`

#### Scenario: Expired token returns invalid code

- GIVEN an expired token in the selected source
- WHEN validation is attempted
- THEN response status is 401 with code `auth_token_invalid`

#### Scenario: Invalid issuer or audience returns invalid code

- GIVEN a token failing issuer or audience policy
- WHEN validation is attempted
- THEN response status is 401 with code `auth_token_invalid`

### Requirement: Claims injection on successful authentication

On successful validation, the middleware MUST inject validated `claims.Claims` into `gin.Context` under a configurable key via `WithContextKey`, and request handling MUST continue.

#### Scenario: Claims available to downstream handlers

- GIVEN a valid token from header or cookie
- WHEN middleware validation succeeds
- THEN downstream handlers can retrieve `claims.Claims` from `gin.Context`
- AND the request reaches the next handler without unauthorized response
