# Delta Specs: issue-51-json-404-405-handlers

---

## Domain: app-bootstrap

### MODIFIED Requirements

### Requirement: Fluent setup methods

`UseConfig(cfg config.Config)`, `UseServer()`, and `UseServerSecurity(v security.Validator)` MUST support fluent chaining on the same `Application` instance. `UseServer()` MUST apply `cfg.Server.ReadTimeout`, `cfg.Server.WriteTimeout`, and `cfg.Server.MaxHeaderBytes` to the underlying `http.Server`. `UseServer()` MUST also register a JSON 404 handler, a JSON 405 handler, and enable `HandleMethodNotAllowed` on the underlying Gin engine.
(Previously: `UseServer()` did not register NoRoute/NoMethod handlers or enable HandleMethodNotAllowed)

#### Scenario: Fluent chain remains supported

- GIVEN an application instance
- WHEN caller chains `UseConfig(...).UseServer().UseServerSecurity(...)`
- THEN each method returns the same instance for continued chaining

#### Scenario: Server runtime limits are applied from config

- GIVEN `UseConfig` receives server runtime-limit values
- WHEN `UseServer()` initializes runtime server wiring
- THEN `http.Server.ReadTimeout` equals `cfg.Server.ReadTimeout`
- AND `http.Server.WriteTimeout` equals `cfg.Server.WriteTimeout`
- AND `http.Server.MaxHeaderBytes` equals `cfg.Server.MaxHeaderBytes`

#### Scenario: Default runtime limits are applied when config uses defaults

- GIVEN `UseConfig` receives a config built with core default server runtime limits
- WHEN `UseServer()` initializes runtime server wiring
- THEN `http.Server` runtime-limit fields equal the documented defaults

#### Scenario: GET on unregistered route returns 404 JSON

- GIVEN `UseServer()` has been called
- WHEN a GET request is sent to a path that has no registered route
- THEN the response status is `404 Not Found`
- AND the response body is `{"code":"not_found","message":"route not found"}`
- AND the `Content-Type` header is `application/json`

#### Scenario: Wrong HTTP method on registered route returns 405 JSON

- GIVEN `UseServer()` has been called and a route is registered for path `/ping` with method GET
- WHEN a DELETE request is sent to `/ping`
- THEN the response status is `405 Method Not Allowed`
- AND the response body is `{"code":"method_not_allowed","message":"method not allowed"}`
- AND the `Content-Type` header is `application/json`

#### Scenario: Correct method on registered route returns normal response (regression)

- GIVEN a route is registered for `GET /ping` returning `200`
- WHEN a GET request is sent to `/ping`
- THEN the response status is `200 OK`
- AND the JSON 404/405 handlers do not interfere

#### Scenario: Health and readiness routes still respond (regression)

- GIVEN `UseServer()` has been called with health/readiness presets registered
- WHEN GET requests are sent to `/health` and `/ready`
- THEN both routes respond with their normal status codes
- AND neither triggers the 404 JSON handler

---

## Domain: errors

### MODIFIED Requirements

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

(Previously: 9 constants — auth error codes only; routing error codes did not exist)

#### Scenario: Constant string stability

- GIVEN the `errors` package is imported
- WHEN any of the 11 constants is read
- THEN its value MUST equal the exact string in the table above

#### Scenario: All constants accessible at package level

- GIVEN the module is compiled
- WHEN a consumer references any of the 11 constants
- THEN compilation succeeds without additional imports

### MODIFIED Requirements

### Requirement: Test Coverage for String Stability

The `errors` package MUST include a table-driven test file `codes_test.go` that verifies the string value of all 11 constants.
(Previously: test covered 9 constants — auth error codes only)

#### Scenario: Table-driven test covers all 11 constants

- GIVEN `errors/codes_test.go` exists
- WHEN the test runs
- THEN each of the 11 constants is asserted against its expected string value
- AND the test passes with no failures
