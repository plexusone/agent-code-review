# Tasks

## Open

- [ ] **Add MCP server tests** - Unit tests for `internal/mcp` with mock review client
- [ ] **Add CLI integration tests** - Test CLI commands with mock GitHub server
- [ ] **Add CHANGELOG** - Track changes for releases

## Completed

- [x] **Remove local replace directive** - Clean up `go.mod` before pushing to origin
- [x] **Add tests** - Unit tests for `pkg/review`
- [x] **Add CI workflow** - GitHub Actions workflow for lint/test/build
- [x] **Add error wrapping** - Wrap errors in SDK write operations for better debugging
- [x] **Add golangci-lint config** - Create `.golangci.yml` with project-specific linting rules
- [x] **Add Makefile** - Common targets for build, test, lint, install
