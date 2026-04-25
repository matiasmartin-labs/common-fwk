# Delta for framework-bootstrap

## MODIFIED Requirements

### Requirement: Bootstrap contains no business logic

Bootstrap artifacts created for the initial bootstrap phase MUST NOT include runtime/business behavior; they SHALL be limited to module metadata, package declarations/docs, and CI wiring needed for compile/test validation. This guard SHALL remain applicable to bootstrap-only packages, and SHALL NOT prevent implementation growth in packages explicitly evolved by later approved capabilities (including `config` for `config-core`).
(Previously: Structural-only expectations could be read as applying broadly to `config`, blocking later capability work.)

#### Scenario: Bootstrap files are structural only

- GIVEN files created by change `bootstrap-common-fwk`
- WHEN bootstrap artifacts are reviewed
- THEN no API handlers, auth flows, or configuration runtime logic is present
- AND scaffold remains structural/documentary in nature

#### Scenario: Business behavior is rejected during bootstrap phase

- GIVEN a bootstrap change attempts to add functional business code
- WHEN evaluating conformance to this specification
- THEN the change is considered non-compliant for phase `sdd-spec`

#### Scenario: Bootstrap guard allows approved config evolution

- GIVEN change `issue-2-config-core` includes implementation files in `config/`
- WHEN bootstrap structural guards are evaluated
- THEN those guards do not fail solely because `config/` contains non-doc implementation files
- AND bootstrap-only package guard intent is preserved for unaffected packages
