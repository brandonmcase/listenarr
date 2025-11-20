# Coding Best Practices

## General Guidelines

- Write clean, maintainable, and well-documented code
- Follow language-specific conventions and style guides
- Prioritize readability and clarity over cleverness
- Use meaningful variable and function names
- Keep functions small and focused on a single responsibility
- Avoid premature optimization
- Write self-documenting code with clear intent
- **Write tests for all code** - Testing is not optional
- **Run tests frequently** - Don't wait until the end to test

## Workflows

### Development Cycle
1. **Before Starting Work**
   - Read and understand `/ai/CodingBestPractices.md`
   - Review relevant context files in `/ai/Contexts/`
   - Check existing TODOs and follow-ups
   - Ensure development environment is set up with required tools

2. **During Development**
   - Write code that follows best practices
   - **Write tests for all new code** - Every function, module, and feature must have corresponding tests
   - Update context files as the codebase evolves
   - Document complex logic and decisions
   - Run tests frequently during development

3. **After Each Unit of Work**
   - **Always read, update, and maintain `/ai/Contexts` files as part of the development cycle.**
   - **Run all tests and ensure they pass** - `make test` or `go test ./...`
   - **Run linting and ensure it passes** - `make lint` or `golangci-lint run`
   - **Run `build-all.sh` to validate changes** (unless documentation-only)
   - Update relevant context files to reflect changes
   - Document any follow-ups or TODOs

4. **Before Committing**
   - Follow `/ai/CommitMessageTemplate.md` for commit messages
   - **MANDATORY: All tests must pass** - `make test` must succeed
   - **MANDATORY: All linting must pass** - `make lint` must succeed
   - **MANDATORY: Code must be formatted** - `make fmt` must be run
   - Verify `build-all.sh` completes successfully
   - **DO NOT commit if tests or linting fail**

5. **Before Opening Pull Requests**
   - Follow `/ai/PullRequestTemplate.md`
   - Ensure all context files are up to date
   - **All tests must pass** - Include test results in PR description
   - **All linting must pass** - Include lint results in PR description
   - Include testing steps and verification
   - Run API endpoint tests if server changes were made

## Build System

### build-all.sh
- Located at: `/ai/build/build-all.sh`
- **Must run at the end of every development cycle** to validate the codebase
- **All tests and linting must pass** - Build will fail if tests fail
- If issues are found, they must be **troubleshooted and fixed before work is considered complete**
- **Exception**: If a task only involves creating documentation or retrieving information (no code changes), `build-all.sh` does not need to be run

### Build Process
1. Preflight checks (dependencies, tools: Go, Node.js, curl)
2. Go format check (`go fmt`)
3. Go vet (`go vet`)
4. Backend build compilation
5. **Unit tests with coverage** (`go test`)
6. Frontend linting and build (if Node.js available)
7. API endpoint testing (if server is running)

### Running build-all.sh
```bash
make build-all
# or
bash .ai/build/build-all.sh
```

## Testing Requirements

### Test Coverage
- **All code must have tests** - Every package, function, and feature
- **Test coverage should be maintained** - Aim for >80% coverage on critical paths
- **Tests must pass before committing** - No exceptions
- **Integration tests** - Test API endpoints and external integrations
- **Unit tests** - Test individual functions and modules in isolation

### Running Tests
```bash
# Run all tests
make test
# or
go test ./...

# Run tests with coverage
make test-coverage

# Run tests with verbose output
make test-verbose

# Test API endpoints (requires running server)
make test-endpoints
```

### Test Structure
- Tests should be in `*_test.go` files alongside the code they test
- Use `github.com/stretchr/testify` for assertions
- Test files should mirror package structure
- Integration tests can be in separate `*_integration_test.go` files

### Test Best Practices
- Write tests before or alongside code (TDD when possible)
- Test both success and failure cases
- Test edge cases and boundary conditions
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Keep tests fast and independent

## Documentation Rules

- Keep documentation up to date with code changes
- Update context files after significant changes
- Document architectural decisions
- Include examples for complex APIs
- Maintain clear README files
- **Document test coverage and testing approach**

## Context Files

- **Location**: `/ai/Contexts/`
- **Purpose**: Maintain project state and knowledge
- **Requirement**: **Always read, update, and maintain `/ai/Contexts` files as part of the development cycle.**
- Agents must update this folder after each cycle to reflect project state, so future prompts don't need to scan the entire codebase
- Start with `overview.md` containing:
  - An outline of project areas (auth, api, ui, etc.)
  - A mapping of contexts to corresponding parts of the codebase
- Add more context files as needed (`auth.md`, `api.md`, etc.)

## Commit Messages

- **When creating commits, follow `/ai/CommitMessageTemplate.md`**
- Use conventional commit format
- Include clear, descriptive messages
- Reference related issues when applicable

## Pull Requests

- **When opening pull requests, follow `/ai/PullRequestTemplate.md`**
- Include summary, related issues, changes, and testing steps
- Ensure checklist items are completed

## Reviews System

- **Location**: `/ai/reviews/{year-month-day}/{year_month_day_hour_minute_title}.md`
- At the end of every review, add a section:
  ```markdown
  ## Suggested Next Prompt
  - [ ] List unfinished tasks, follow-ups, or improvements
  ```

## Statistics and Reporting

- Track code metrics (lines of code, test coverage, etc.)
- Update stats after significant changes
- Store reports in appropriate locations

## Required Tools

### Backend Development (Go)
- **Go 1.21+** - Programming language and toolchain
  - Install: https://go.dev/dl/
  - Verify: `go version`
- **golangci-lint** (recommended) - Comprehensive linter
  - Install: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
  - Verify: `golangci-lint --version`

### Frontend Development (React/TypeScript)
- **Node.js 18+** - JavaScript runtime
  - Install: https://nodejs.org/
  - Verify: `node --version`
- **npm or yarn** - Package manager (comes with Node.js)

### Testing & Verification
- **curl** - For API endpoint testing
  - Usually pre-installed on macOS/Linux
  - Verify: `curl --version`

### Docker (Optional but Recommended)
- **Docker** - For containerized development and deployment
  - Install: https://www.docker.com/
  - Verify: `docker --version`

### Development Tools
- **Make** - For running build commands (usually pre-installed)
  - Verify: `make --version`
- **Git** - Version control (usually pre-installed)
  - Verify: `git --version`

## Enforcement Rules

- **Context files must be updated** as part of "After Each Unit of Work"
- **All tests must pass** before committing code
- **All linting must pass** before committing code
- **Code must be formatted** (`go fmt`) before committing
- **`build-all.sh` must pass** before declaring work complete, **unless** the task was *documentation-only*
- **New code must include tests** - No code should be committed without corresponding tests
- Reviews **must** include a "Suggested Next Prompt" section
- **DO NOT skip tests or linting** - These are mandatory quality gates

