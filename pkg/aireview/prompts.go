// Package aireview provides AI-powered code review functionality.
package aireview

import "fmt"

// Scope represents the focus area for a code review.
type Scope string

const (
	// ScopeFull performs a comprehensive review covering all aspects.
	ScopeFull Scope = "full"
	// ScopeSecurity focuses on security vulnerabilities and risks.
	ScopeSecurity Scope = "security"
	// ScopeStyle focuses on code style, readability, and best practices.
	ScopeStyle Scope = "style"
	// ScopePerformance focuses on performance issues and optimizations.
	ScopePerformance Scope = "performance"
)

// ValidScopes returns all valid review scopes.
func ValidScopes() []Scope {
	return []Scope{ScopeFull, ScopeSecurity, ScopeStyle, ScopePerformance}
}

// IsValidScope checks if a scope string is valid.
func IsValidScope(s string) bool {
	switch Scope(s) {
	case ScopeFull, ScopeSecurity, ScopeStyle, ScopePerformance:
		return true
	default:
		return false
	}
}

// systemPrompt is the base system prompt for all reviews.
const systemPrompt = `You are an expert code reviewer. You analyze pull request diffs and provide thorough, actionable feedback.

Your reviews should be:
- Specific: Point to exact lines and explain why something is an issue
- Constructive: Suggest solutions, not just problems
- Respectful: Critique the code, not the author
- Prioritized: Distinguish critical issues from nice-to-haves
- Balanced: Acknowledge good patterns and practices

Always structure your review with clear sections and use markdown formatting.`

// scopeInstructions contains focus-specific instructions for each scope.
var scopeInstructions = map[Scope]string{
	ScopeFull: `Perform a comprehensive code review covering:

1. **Correctness**
   - Does the code do what it's supposed to do?
   - Are there edge cases that aren't handled?
   - Are there off-by-one errors or boundary conditions?

2. **Security**
   - Check for injection vulnerabilities (SQL, XSS, command injection)
   - Validate authentication and authorization logic
   - Look for hardcoded secrets or credentials
   - Check for sensitive data exposure

3. **Performance**
   - Identify N+1 queries or unnecessary database calls
   - Look for inefficient algorithms or data structures
   - Check for memory leaks or resource cleanup issues

4. **Maintainability**
   - Is the code readable and self-documenting?
   - Are functions and methods appropriately sized?
   - Is there unnecessary complexity or over-engineering?

5. **Testing**
   - Are there adequate tests for the changes?
   - Do tests cover edge cases?
   - Are tests meaningful (not just for coverage)?`,

	ScopeSecurity: `Focus exclusively on security issues:

1. **Injection Vulnerabilities**
   - SQL injection
   - Cross-site scripting (XSS)
   - Command injection
   - Path traversal

2. **Authentication & Authorization**
   - Broken authentication
   - Missing authorization checks
   - Privilege escalation risks

3. **Sensitive Data**
   - Hardcoded secrets, API keys, passwords
   - Sensitive data in logs
   - Unencrypted data transmission
   - PII exposure

4. **Input Validation**
   - Missing input sanitization
   - Type confusion
   - Buffer overflows

5. **Dependencies**
   - Known vulnerable dependencies
   - Outdated security patches

Rate severity: CRITICAL, HIGH, MEDIUM, LOW`,

	ScopeStyle: `Focus on code style, readability, and best practices:

1. **Naming**
   - Are variable, function, and class names descriptive?
   - Do names follow language conventions?
   - Are abbreviations clear and consistent?

2. **Readability**
   - Is the code self-documenting?
   - Are complex sections adequately commented?
   - Is the code flow easy to follow?

3. **Structure**
   - Are functions appropriately sized?
   - Is there proper separation of concerns?
   - Are abstractions at the right level?

4. **Consistency**
   - Does the code follow project conventions?
   - Is formatting consistent?
   - Are patterns used consistently?

5. **Best Practices**
   - Are language idioms used correctly?
   - Is error handling appropriate?
   - Are there any anti-patterns?`,

	ScopePerformance: `Focus on performance issues and optimizations:

1. **Database**
   - N+1 query problems
   - Missing indexes
   - Unnecessary queries
   - Large result sets without pagination

2. **Algorithms**
   - Time complexity issues
   - Inefficient data structures
   - Unnecessary iterations
   - Redundant computations

3. **Memory**
   - Memory leaks
   - Large allocations
   - Missing cleanup/disposal
   - Unbounded caches

4. **I/O**
   - Synchronous blocking calls
   - Missing connection pooling
   - Unbuffered I/O
   - Missing timeouts

5. **Concurrency**
   - Race conditions
   - Deadlock risks
   - Thread safety issues
   - Inefficient locking`,
}

// outputFormat is the expected output format for reviews.
const outputFormat = `
Structure your response as:

## Summary
[1-2 sentence overview of the changes and overall assessment]

## Findings

### Critical
[Issues that must be fixed before merging - leave empty if none]

### Suggestions
[Recommendations for improvement, not blocking - leave empty if none]

### Positive
[Things done well, good patterns observed - include at least one if applicable]

## Verdict
[One of: APPROVE, COMMENT, or REQUEST_CHANGES]
[Brief justification for the verdict]`

// BuildPrompt constructs the full prompt for a code review.
func BuildPrompt(scope Scope, prTitle, prBody, diff string) string {
	instructions, ok := scopeInstructions[scope]
	if !ok {
		instructions = scopeInstructions[ScopeFull]
	}

	return fmt.Sprintf(`%s

## Review Focus
%s

## Output Format
%s

---

## Pull Request

**Title:** %s

**Description:**
%s

## Diff

%s`, systemPrompt, instructions, outputFormat, prTitle, prBody, diff)
}

// BuildSystemPrompt returns the system prompt for the LLM.
func BuildSystemPrompt() string {
	return systemPrompt
}

// BuildUserPrompt constructs the user message for a code review.
func BuildUserPrompt(scope Scope, prTitle, prBody, diff string) string {
	instructions, ok := scopeInstructions[scope]
	if !ok {
		instructions = scopeInstructions[ScopeFull]
	}

	return fmt.Sprintf(`Review this pull request with the following focus:

## Review Focus
%s

## Output Format
%s

---

## Pull Request

**Title:** %s

**Description:**
%s

## Diff

`+"```diff\n%s\n```", instructions, outputFormat, prTitle, prBody, diff)
}
