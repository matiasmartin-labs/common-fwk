# Delta Spec: issue-29-config-viper-kebab-case-keys

## ADDED Requirements — config-viper-kebab-case-keys

### Requirement: Canonical kebab-case file keys

The Viper adapter MUST treat kebab-case as the canonical file-key contract for mapped configuration fields.

#### Scenario: Kebab-case fixture decodes successfully
- GIVEN a valid adapter configuration file using kebab-case keys
- WHEN `config/viper.Load` is executed
- THEN the adapter returns a valid mapped `config.Config`
- AND mapped values are identical to equivalent logical data previously expressed in legacy camelCase

### Requirement: Deterministic legacy-key compatibility

The adapter MUST define explicit compatibility behavior for legacy camelCase keys and preserve deterministic outcomes.

#### Scenario: Legacy camelCase remains compatible
- GIVEN a valid file that uses legacy camelCase keys for historical compatibility
- WHEN `config/viper.Load` is executed
- THEN the adapter still maps successfully to `config.Config`
- AND behavior is documented as compatibility mode

#### Scenario: Canonical and legacy keys both present
- GIVEN a file where canonical kebab-case and legacy camelCase aliases exist for the same logical field
- WHEN `config/viper.Load` is executed
- THEN the resulting value follows a documented deterministic precedence rule

### Requirement: Environment override determinism preserved

The key-style migration MUST NOT alter deterministic environment override behavior.

#### Scenario: Env override still wins with override enabled
- GIVEN file values in canonical kebab-case
- AND environment variables for the same logical fields
- WHEN `EnvOverride=true`
- THEN environment values override file values exactly as before

#### Scenario: Env override disabled preserves file source of truth
- GIVEN file values in canonical kebab-case
- AND environment variables for the same logical fields
- WHEN `EnvOverride=false`
- THEN file values remain effective exactly as before

### Requirement: Documentation uses kebab-case examples

All configuration examples under `/docs/*` and `README.md` MUST use kebab-case keys for file-based configuration snippets.

#### Scenario: Docs aligned with canonical key style
- GIVEN project documentation containing config snippets
- WHEN documentation is reviewed after this change
- THEN config examples use kebab-case keys
- AND no primary example presents camelCase as canonical style

## Out of Scope
- Renaming environment variable keys.
- Changes to core config validation rules.
