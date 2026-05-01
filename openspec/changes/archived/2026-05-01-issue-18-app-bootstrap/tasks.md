# Tasks: issue-18-app-bootstrap

## Phase 1: Application foundation

- [x] 1.1 Create `app/application.go` with package imports, constructor `NewApplication()`, and package sentinel errors (`ErrServerNotConfigured`, `ErrValidatorNotConfigured`, `ErrNilHandler`, `ErrInvalidPath`, `ErrNilListener`).
- [x] 1.2 Define `Application` struct in `app/application.go` with fields: `cfg config.Config`, `server http.Server`, `handler *gin.Engine`, `validator security.Validator`, `serverReady bool`, `securityReady bool`. (depends on 1.1)
- [x] 1.3 Initialize default engine/state in `NewApplication()` so a new instance is safe for fluent chaining and later route registration. (depends on 1.2)

## Phase 2: Fluent bootstrap methods

- [x] 2.1 Implement `UseConfig(cfg config.Config) *Application` to set config on the same instance and return receiver for chaining in `app/application.go`. (depends on 1.3)
- [x] 2.2 Implement `UseServer() *Application` to initialize server/handler wiring and set `serverReady=true`; return same instance. (depends on 2.1)
- [x] 2.3 Implement `UseServerSecurity(v security.Validator) *Application` to store validator, set `securityReady=true`, and return same instance. (depends on 2.2)

## Phase 3: Route and run operations

- [x] 3.1 Implement shared guards/helpers in `app/application.go` for path validation, nil handler validation, and server/security readiness checks with deterministic errors. (depends on 2.3)
- [x] 3.2 Implement `RegisterGET(path string, h gin.HandlerFunc) error` with guard checks and `handler.GET(...)` registration. (depends on 3.1)
- [x] 3.3 Implement `RegisterPOST(path string, h gin.HandlerFunc) error` with guard checks and `handler.POST(...)` registration. (depends on 3.1)
- [x] 3.4 Implement `RegisterProtectedGET(path string, h gin.HandlerFunc) error`, wiring `http/gin.NewAuthMiddleware(a.validator)` before handler. (depends on 3.1)
- [x] 3.5 Implement `Run() error` as blocking startup (`ListenAndServe`) with precondition checks and error returns only. (depends on 2.2)
- [x] 3.6 Implement `RunListener(l net.Listener) error` with nil-listener guard and `server.Serve(l)` for testable runtime path. (depends on 2.2)

## Phase 4: Tests by scenario

- [x] 4.1 Create `app/application_test.go` test group: bootstrap chain (`UseConfig().UseServer().UseServerSecurity()`) preserves same pointer and marks readiness state. (depends on 2.3)
- [x] 4.2 Add test group: route registration (`RegisterGET`, `RegisterPOST`, `RegisterProtectedGET`) succeeds after full bootstrap and mounts distinct paths. (depends on 3.2, 3.3, 3.4)
- [x] 4.3 Add protected-enforcement tests: missing token returns 401 and invalid token returns 401 on `RegisterProtectedGET` routes. (depends on 3.4)
- [x] 4.4 Add ordering-guard tests: calling `Register*`, `RegisterProtectedGET`, and `Run` before prerequisites returns expected sentinel errors. (depends on 3.1, 3.5)
- [x] 4.5 Add run-path tests: `RunListener` serves a request on ephemeral listener and returns runtime/startup errors cleanly. (depends on 3.6)
- [x] 4.6 Add `Run()` behavior test using controlled server setup to verify blocking path delegates to configured server and propagates errors. (depends on 3.5)
