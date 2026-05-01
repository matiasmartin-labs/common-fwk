# Delta for app-bootstrap

## ADDED Requirements

### Requirement: Optional config-based security bootstrap convenience

The application bootstrap API MAY provide an optional helper that derives validator wiring from already loaded config. This helper MUST preserve existing explicit `UseServerSecurity` behavior and MUST fail deterministically with contextual errors when config prerequisites are invalid or incomplete.

#### Scenario: Config-based helper succeeds with valid JWT mode configuration

- GIVEN an application instance with valid security config for HS256 or RS256
- WHEN the optional config-based security bootstrap helper is invoked
- THEN validator wiring is configured successfully for protected routes

#### Scenario: Config-based helper fails deterministically on invalid security config

- GIVEN an application instance with incomplete or invalid JWT mode configuration
- WHEN the optional config-based security bootstrap helper is invoked
- THEN it returns a contextual error
- AND no partial security wiring is applied
