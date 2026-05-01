# Apply Progress: issue-29-config-viper-kebab-case-keys

## Completed

- Updated Viper raw mapping tags to canonical kebab-case keys in `config/viper/mapping.go`.
- Preserved deterministic env override behavior by switching internal override target paths to canonical kebab-case keys.
- Added legacy camelCase compatibility layer in `config/viper/loader.go` that backfills canonical paths only when canonical keys are absent.
- Enforced deterministic precedence for mixed key styles: canonical kebab-case wins over legacy aliases.
- Updated `config/viper/loader_test.go` fixtures to kebab-case for core scenarios.
- Added compatibility regression tests for legacy camelCase keys and mixed-style precedence behavior.
- Updated docs/examples to canonical kebab-case in `README.md`, migration guide, and release checklist verification section.
- Added docs index page `docs/home.md` and prepared `docs/releases/v0.2.0-checklist.md` as the target release checklist for this change.

## Verification

- `go test ./config/viper` passed.
- `go test ./...` passed.

## Notes

- Environment variable names were intentionally not changed.
- Legacy camelCase file keys remain compatibility-only and are no longer documented as canonical.
