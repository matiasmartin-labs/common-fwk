# Tasks: JSON 404/405 Default Handlers

## Phase 1: Foundation — Error Code Constants

- [ ] 1.1 In `errors/codes.go`, append `CodeNotFound = "not_found"` and `CodeMethodNotAllowed = "method_not_allowed"` to the existing const block
- [ ] 1.2 In `errors/codes_test.go`, extend the table-driven test to cover all 11 constants (add rows for `CodeNotFound` and `CodeMethodNotAllowed`)

## Phase 2: Core Implementation — UseServer() Handlers

- [ ] 2.1 In `app/application.go`, add import alias `fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"` (avoid collision with stdlib `errors`)
- [ ] 2.2 In `UseServer()`, set `a.handler.HandleMethodNotAllowed = true` before handler registration
- [ ] 2.3 In `UseServer()`, register `a.handler.NoRoute(...)` calling `c.AbortWithStatusJSON(http.StatusNotFound, httpgin.ErrorResponse{Code: fwkerrors.CodeNotFound, Message: "route not found"})`
- [ ] 2.4 In `UseServer()`, register `a.handler.NoMethod(...)` calling `c.AbortWithStatusJSON(http.StatusMethodNotAllowed, httpgin.ErrorResponse{Code: fwkerrors.CodeMethodNotAllowed, Message: "method not allowed"})`

## Phase 3: Testing

- [ ] 3.1 In `app/application_test.go`, add `TestUseServer_NoRoute_JSON`: GET `/nonexistent` → assert status 404 and `{"code":"not_found","message":"route not found"}`
- [ ] 3.2 In `app/application_test.go`, add `TestUseServer_NoMethod_JSON`: register POST `/ping`, send DELETE `/ping` → assert status 405 and `{"code":"method_not_allowed","message":"method not allowed"}`
- [ ] 3.3 In `app/application_test.go`, add `TestUseServer_CorrectMethod_NoInterference`: register GET `/ping` returning 200, send GET `/ping` → assert status 200 (regression guard)

## Phase 4: Verification

- [ ] 4.1 Run `go test ./...` and confirm all tests pass with no failures
