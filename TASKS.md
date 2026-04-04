# Tasks

## High Priority

- [ ] **Remove local replace directive** - Clean up `go.mod` before pushing to origin
- [ ] **Add tests** - Unit tests for `pkg/review` and integration tests for CLI commands

## Medium Priority

- [ ] **Add CI workflow** - GitHub Actions workflow for lint/test/build
- [ ] **Consider error wrapping** - Use `fmt.Errorf("context: %w", err)` for better debugging and stack traces

## Low Priority

- [ ] **Add golangci-lint config** - Create `.golangci.yml` with project-specific linting rules
- [ ] **Add Makefile** - Common targets for build, test, lint, install
