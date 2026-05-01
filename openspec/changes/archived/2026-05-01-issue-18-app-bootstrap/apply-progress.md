# Apply Progress: issue-18-app-bootstrap

Mode: Standard (strict_tdd=false)

## Task Checklist

### Phase 1: Application foundation
- ✅ 1.1 Create `app/application.go` with imports, `NewApplication()`, and sentinel errors.
- ✅ 1.2 Define `Application` struct with required fields.
- ✅ 1.3 Initialize default gin engine/state in `NewApplication()`.

### Phase 2: Fluent bootstrap methods
- ✅ 2.1 Implement `UseConfig(cfg config.Config) *Application`.
- ✅ 2.2 Implement `UseServer() *Application`.
- ✅ 2.3 Implement `UseServerSecurity(v security.Validator) *Application`.

### Phase 3: Route and run operations
- ✅ 3.1 Implement shared guard helpers for readiness/path/handler.
- ✅ 3.2 Implement `RegisterGET`.
- ✅ 3.3 Implement `RegisterPOST`.
- ✅ 3.4 Implement `RegisterProtectedGET` with `http/gin.NewAuthMiddleware(a.validator)`.
- ✅ 3.5 Implement `Run() error` with server readiness check and `ListenAndServe` delegation.
- ✅ 3.6 Implement `RunListener(l net.Listener) error` with nil-listener guard and `Serve` delegation.

### Phase 4: Tests
- ✅ 4.1 Bootstrap chain tests.
- ✅ 4.2 Route registration tests.
- ✅ 4.3 Protected enforcement tests (missing/invalid token => 401).
- ✅ 4.4 Ordering guard tests.
- ✅ 4.5 `RunListener` behavior tests.
- ✅ 4.6 `Run()` delegation/error propagation tests.

## Decisions / Notes

- Implemented sentinel names per current change key decisions and user constraints:
  - `ErrServerNotReady`
  - `ErrSecurityNotReady`
  - `ErrInvalidPath`
  - `ErrNilHandler`
  - `ErrNilListener`
- `UseServerSecurity(nil)` keeps `securityReady=false`; protected registration fails deterministically with `ErrSecurityNotReady`.
- `Run()` sets `server.Addr` from `cfg.Server` using `net.JoinHostPort(host, port)` and delegates directly to `ListenAndServe()`.
- No package-global `App` singleton was introduced.

## Verification

- `go build ./app/...` ✅
- `go test ./app/...` ✅

## Follow-up Fixes

### Bootstrap guard conformance update
- ✅ Updated `bootstrap_guard_test.go` so `app` is no longer treated as a structural-only bootstrap package.
- ✅ Added explicit evolution assertion: `TestAppPackageCanEvolveBeyondBootstrapDocs`.
- ✅ Searched for additional guard/conformance references to `app` structural-only assumption; no other files required changes.

## Verification (follow-up)

- `go test ./...` ✅
