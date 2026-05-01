# Delta for release-readiness-docs

## ADDED Requirements

### Requirement: Release checklist covers JWT mode behavior

Release-readiness documentation MUST describe RS256 keypair bootstrap behavior, HS256 backward compatibility expectations, and verification checkpoints for both JWT modes before release publication.

#### Scenario: Release checklist includes mode-specific checks

- GIVEN release checklist documentation for the targeted version
- WHEN maintainers execute pre-release verification
- THEN checklist items explicitly validate HS256 compatibility and RS256 bootstrap behavior
