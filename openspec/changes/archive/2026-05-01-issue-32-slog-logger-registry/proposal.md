# Proposal: issue-32-slog-logger-registry

## Intent
`common-fwk` lacks a logging subsystem and `app.GetLogger(name)`. This change introduces deterministic named logging with explicit core→adapter boundaries and documented configuration behavior.

## Scope
### In Scope
- Add a core logging contract plus slog-backed registry (`Debugf/Infof/Warnf/Errorf`).
- Add logging config in core + viper adapter: `logging.enabled`, `logging.level`, `logging.format`, `logging.loggers.<name>.*`.
- Enforce deterministic precedence/validation and test acceptance for filtering, fields, and logger isolation.
- Expose `Application.GetLogger(name)` and synchronize `README.md` + `/docs/*` (including Loki collector-first guidance).

### Out of Scope
- Alternate backends, async pipelines, runtime hot-reload, remote sink implementation.
- Any non-logging app/security behavior change.

## Capabilities
### New Capabilities
- `logging-registry`: Named logger registry contract + slog adapter with deterministic filtering/formatting.

### Modified Capabilities
- `app-bootstrap`: logger accessor lifecycle/guard behavior.
- `config-core`: logging model, defaults, validation, precedence.
- `config-viper-adapter`: logging key mapping + env/file deterministic precedence.
- `adoption-migration-guide`: logging adoption notes and Loki-oriented usage guidance.

## Approach
Use exploration Approach 2: core contract + adapter. Root config defines defaults; per-logger config overrides deterministically:
- `enabled`: root gate with optional logger override.
- `level`: logger override else root level.
- `format`: root-level `json|text` handler selection.

## Affected Areas
| Area | Impact | Description |
|------|--------|-------------|
| `app/application.go`, `app/application_test.go`, `app/doc.go` | Modified | Registry lifecycle, `GetLogger(name)`, acceptance tests/docs |
| `config/types.go`, `config/constructors.go`, `config/validate.go` | Modified | Logging config model/defaults/validation |
| `config/viper/mapping.go`, `config/viper/loader.go`, `config/viper/*_test.go` | Modified | Deterministic logging key mapping/overrides |
| `README.md`, `docs/*` | Modified | Config examples, precedence notes, Loki guidance |

## Risks
| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Root vs logger precedence ambiguity | Med | Spec precedence matrix + acceptance tests |
| Text/JSON output drift | Med | Assert required fields: `logger`,`ts`,`level`,`msg` |
| Logger state leakage across names | Low-Med | Instance-scoped registry + isolation tests |
| Docs drift from behavior | Med | Docs as acceptance criteria in same change |

## Rollback Plan
Revert logger registry wiring and logging config additions together: remove `GetLogger`, remove logging mapping/validation fields, and rollback logging docs. Existing bootstrap/security flows remain intact.

## Dependencies
- Go stdlib `log/slog` only.

## Success Criteria
- [ ] `app.GetLogger(name)` returns deterministic named loggers.
- [ ] Precedence + level filtering pass acceptance tests.
- [ ] Text/JSON include `logger`,`ts`,`level`,`msg`.
- [ ] Logger instances remain isolated.
- [ ] `README.md` and `/docs/*` reflect final logging contract and Loki guidance.
