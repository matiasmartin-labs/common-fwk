# Framework Bootstrap Specification

## Purpose

Define the minimal Go module and package scaffold required by issue #1 so the repository compiles without introducing business behavior.

## Requirements

### Requirement: Module and package scaffold compiles

The repository MUST declare module path `github.com/matiasmartin-labs/common-fwk` and SHALL provide these Go packages as tracked, compilable package boundaries: `app`, `config`, `config/viper`, `security`, `http/gin`, and `errors`.

#### Scenario: Bootstrap scaffold compiles from repository root

- GIVEN the bootstrap scaffold is present in the repository
- WHEN `go test ./...` is executed from repository root
- THEN command execution succeeds with exit code `0`
- AND each listed package is discovered as a valid Go package

#### Scenario: Minimal package stubs remain compilable

- GIVEN each bootstrap package contains only minimal stub files (for example package docs)
- WHEN `go test ./...` is executed
- THEN compilation still succeeds without requiring runtime implementation

### Requirement: Bootstrap contains no business logic

Bootstrap artifacts MUST NOT include runtime/business behavior in this phase; they SHALL be limited to module metadata, package declarations/docs, and CI wiring needed for compile/test validation.

#### Scenario: Bootstrap files are structural only

- GIVEN files created by change `bootstrap-common-fwk`
- WHEN bootstrap artifacts are reviewed
- THEN no API handlers, auth flows, or configuration runtime logic is present
- AND scaffold remains structural/documentary in nature

#### Scenario: Business behavior is rejected during bootstrap phase

- GIVEN a bootstrap change attempts to add functional business code
- WHEN evaluating conformance to this specification
- THEN the change is considered non-compliant for phase `sdd-spec`
