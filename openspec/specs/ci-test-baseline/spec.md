# CI Test Baseline Specification

## Purpose

Define the minimal CI behavior required by issue #1 to continuously validate that the bootstrap scaffold remains compilable.

## Requirements

### Requirement: CI executes Go test baseline

The system MUST provide a CI workflow that runs `go test ./...` for pull requests and SHALL fail the run if that command fails.

#### Scenario: Pull request triggers baseline Go tests

- GIVEN a pull request targeting the repository
- WHEN CI workflows execute for that pull request
- THEN CI runs `go test ./...`
- AND the workflow reports pass/fail based on the command result

#### Scenario: Failed tests fail the CI workflow

- GIVEN `go test ./...` returns a non-zero exit code in CI
- WHEN the workflow evaluates job status
- THEN the workflow result is failure

### Requirement: CI scope stays bootstrap-minimal

The baseline CI workflow SHOULD remain limited to bootstrap compile/test validation in this phase and MUST NOT require additional quality gates (for example lint, release, or coverage thresholds) before the bootstrap is accepted.

#### Scenario: Baseline CI has no extra mandatory gates

- GIVEN issue #1 bootstrap acceptance criteria
- WHEN reviewing the CI workflow for this phase
- THEN `go test ./...` is required
- AND extra gates beyond this baseline are not mandatory
