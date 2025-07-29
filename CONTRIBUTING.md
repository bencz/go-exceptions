# Contributing to GoExceptions

Thank you for your interest in contributing to GoExceptions! This document provides guidelines and information for contributors.

## Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-exceptions.git
   cd go-exceptions
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** and commit them
5. **Push to your fork** and create a Pull Request

## Development Guidelines

### Code Style

- Follow standard Go conventions and `gofmt` formatting
- Use clear, descriptive variable and function names
- Add comprehensive comments for public APIs
- Keep functions focused and single-purpose
- Maintain consistency with existing code style

### Testing Requirements

We maintain **97.2% code coverage**. All contributions must:

- Include comprehensive unit tests
- Maintain or improve existing coverage
- Pass all existing tests
- Include integration tests for complex features
- Add benchmark tests for performance-critical code

**Running Tests:**
```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Documentation

- Update `doc.go` for new public APIs
- Add examples for new features
- Update README.md if needed
- Include inline code comments
- Create demo examples in `cmd/demo/` for significant features

## What We're Looking For

### High Priority
- **Bug fixes** with test cases
- **Performance improvements** with benchmarks
- **Additional test coverage** for edge cases
- **Documentation improvements**
- **Code optimization** without breaking changes

### Medium Priority
- **New exception types** for common use cases
- **Helper functions** for validation
- **Better error messages** and stack traces
- **API improvements** (backward compatible)

### Lower Priority
- **Breaking changes** (require strong justification)
- **Experimental features** (may be rejected)

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass**:
   ```bash
   go test ./...
   ```

2. **Check code formatting**:
   ```bash
   gofmt -w .
   go vet ./...
   ```

3. **Verify coverage**:
   ```bash
   go test -cover ./...
   ```

4. **Run benchmarks** (if performance-related):
   ```bash
   go test -bench=. ./...
   ```

### PR Requirements

- **Clear title** describing the change
- **Detailed description** of what and why
- **Test coverage** for new code
- **Documentation updates** if needed
- **Link to related issues** (if applicable)
- **All CI checks passing**

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Performance improvement
- [ ] Documentation update
- [ ] Test improvement

## Testing
- [ ] Added unit tests
- [ ] Added integration tests
- [ ] All existing tests pass
- [ ] Coverage maintained/improved

## Checklist
- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or justified)
```

## Bug Reports

### Before Reporting
- Check existing issues for duplicates
- Verify the bug with latest version
- Create minimal reproduction case

### Bug Report Template
```markdown
**Bug Description**
Clear description of the bug

**Reproduction Steps**
1. Step one
2. Step two
3. Step three

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- Go version: 
- OS: 
- GoExceptions version: 

**Code Sample**
```go
// Minimal code to reproduce
```

## Feature Requests

### Before Requesting
- Check existing issues and discussions
- Consider if it fits the project scope
- Think about backward compatibility

### Feature Request Template
```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why is this feature needed?

**Proposed API**
```go
// Example of how it would be used
```

**Alternatives Considered**
Other approaches you've considered
```

## Development Setup

### Prerequisites
- Go 1.18+ (for generics support)
- Git
- Your favorite Go IDE/editor

### Project Structure
```
go-exceptions/
├── goexceptions.go      # Core implementation
├── doc.go              # Package documentation
├── README.md           # Project overview
├── LICENSE             # MIT license
├── go.mod              # Go module file
├── package_test.go     # Internal tests
├── cmd/demo/           # Example applications
├── tests/              # External test suite
└── coverage.html       # Coverage report
```

### Building and Testing
```bash
# Build the project
go build ./...

# Run all tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestSpecificFunction

# Run benchmarks
go test -bench=BenchmarkTry -benchmem

# Check for race conditions
go test -race ./...
```

## Performance Considerations

- Benchmark any performance-related changes
- Avoid unnecessary allocations
- Consider reflection caching for type operations
- Profile memory usage for large-scale operations
- Maintain backward compatibility

## Code Review Process

### For Contributors
- Be responsive to feedback
- Make requested changes promptly
- Ask questions if feedback is unclear
- Be patient during the review process

### For Reviewers
- Be constructive and helpful
- Focus on code quality and maintainability
- Check test coverage and documentation
- Verify backward compatibility
- Test the changes locally

## Getting Help

- **GitHub Discussions** for questions and ideas
- **GitHub Issues** for bugs and feature requests
- **Email maintainers** for sensitive issues
- **Read the documentation** in `doc.go` and README

## Recognition

Contributors will be:
- Listed in the project contributors
- Mentioned in release notes for significant contributions
- Credited in documentation for major features

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professional communication
- Follow GitHub's community guidelines

## Thank You!

Every contribution, no matter how small, helps make GoExceptions better for everyone. We appreciate your time and effort in improving this project!

---

**Happy Coding!**
