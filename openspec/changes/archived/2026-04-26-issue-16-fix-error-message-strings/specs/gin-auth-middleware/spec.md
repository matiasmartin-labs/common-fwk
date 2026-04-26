# Delta for gin-auth-middleware

## MODIFIED Requirements

### Requirement: Unauthorized error response contract

When authentication fails, the middleware MUST return HTTP 401 with JSON payload `{ "code": string, "message": string }`. Missing token MUST use `code: auth_token_missing` and `message: "missing authentication token"`. Invalid, malformed, expired, invalid issuer, or invalid audience token outcomes MUST use `code: auth_token_invalid` and `message: "invalid or expired token"`.
(Previously: message string values were `"authentication token is missing"` and `"authentication token is invalid"` — not matching `auth-provider-ms` contract)

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

## ADDED Requirements

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

## Breaking Change Notice

Consumers asserting the previous message strings `"authentication token is missing"` or `"authentication token is invalid"` MUST update their assertions to use the new values or reference the exported constants `MsgTokenMissing` / `MsgTokenInvalid`.
