# Tasks

## High Priority

- [x] **Remove local replace directive** - Clean up `go.mod` before pushing to origin (already done)
- [x] **Add tests** - Unit tests for `pkg/review` and integration tests for CLI commands

## Medium Priority

- [x] **Add CI workflow** - GitHub Actions workflow for lint/test/build (already exists in .github/workflows/)
- [x] **Consider error wrapping** - Use `fmt.Errorf("context: %w", err)` for better debugging and stack traces

## Low Priority

- [x] **Add golangci-lint config** - Create `.golangci.yml` with project-specific linting rules
- [x] **Add Makefile** - Common targets for build, test, lint, install
