# Contributing to rfcquery

Thank you for your interest in contributing to rfcquery! This document provides guidelines to ensure a smooth contribution process.

## Code of Conduct

By participating, you are expected to uphold our code of conduct:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Respect differing viewpoints and experiences

## How to Contribute

### Reporting Bugs

1. **Check existing issues**: Search the [issue tracker](https://github.com/yourusername/rfcquery/issues) first
2. **Create a minimal reproduction**: Provide a small code snippet that demonstrates the issue
3. **Include environment details**: Go version, operating system
4. **Expected vs actual behavior**: Clearly describe what you expected vs what happened

**Bug report template:**
```markdown
**Description**
A clear description of the bug.

**Reproduction**
```go
// Minimal code that reproduces the issue Contributing to rfcquery
```
#### Environment
Go version: 1.x
OS: Linux/Windows/macOS

#### Expected behavior
What should happen.

#### Actual behavior
What actually happens.

#### Additional context
Any other relevant information.

### Suggesting Features

1. **Check the roadmap**: See if your idea is already planned
2. **Open a discussion**: Use GitHub Discussions for initial ideas
3. **Focus on the problem**: Describe the pain point your feature solves
4. **Consider the scope**: RFC3986 compliance is core; plugins are extensible

### Pull Requests

#### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/rfcquery.git
cd rfcquery

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run benchmarks
go test -bench=. -benchmem

# Run linter (if you have golangci-lint)
golangci-lint run
```

### PR Process
    1. Fork and create a branch:
    ```bash
    git chekout -b feature/token-stream-api
    ```

    2. Follow the code style:
     - Use `gofmt` and `go vet`
     - Match existing naming conventions
     - Add comments for exported symbols

    3. Add Tets:
         - Unit tests for new functionality
         - Benchmarks (optional)
         - Test error cases and edge conditions

    4. Update documentation:
        - Update README.md if adding user-facing features
        - Add examples for new APIs
        - Update this file if changing contribution process

    5. Commit with clear messages:
        Usually following [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) will do

    6. Submit PR:
        - Reference related issue(s)


### Code Style Guidelines

#### General
 - Follow Go conventions
 - Always check errors

#### Specific to rfcquery

 - RFC3986 Compliance
 - Error messages: include position, be actionable
 - Backwards compatibility: Don't break existing APIs


 ### License

 By contributing, you agree that your contributions will be licensed under the MIT License.
 
 Thank you for making rfcquery better!
