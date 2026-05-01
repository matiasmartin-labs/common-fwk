# Delta for config-viper-adapter

## ADDED Requirements

### Requirement: Server runtime limits mapping and env overrides

The Viper adapter MUST load `server.read-timeout`, `server.write-timeout`, and `server.max-header-bytes` from configuration files and MUST support deterministic environment overrides for the same keys.

#### Scenario: File values are mapped into core config

- GIVEN a valid config file containing the three server runtime-limit keys
- WHEN the adapter loader runs
- THEN returned `config.Config.Server` contains mapped values for all three keys

#### Scenario: Env overrides take precedence when enabled

- GIVEN file values and environment values for the same runtime-limit keys
- WHEN env override is enabled
- THEN returned server runtime-limit values come from environment inputs
- AND behavior is deterministic for identical input snapshots

### Requirement: Typed failures for runtime-limit decoding and mapping

The adapter MUST return adapter-typed errors when runtime-limit values cannot be decoded or mapped into core types.

#### Scenario: Invalid duration format returns decode-typed error

- GIVEN `server.read-timeout` or `server.write-timeout` has an invalid duration string
- WHEN loading/decoding runs
- THEN the adapter returns a decode-typed error

#### Scenario: Invalid max-header-bytes type returns mapping/decode typed error

- GIVEN `server.max-header-bytes` is not representable as the required numeric type
- WHEN decoding/mapping runs
- THEN the adapter returns an adapter-typed error identifying load failure
