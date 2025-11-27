# Contributing to Excel TUI

First off, thank you for considering contributing to Excel TUI! 🎉

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples**
- **Describe the behavior you observed and what you expected**
- **Include screenshots if possible**
- **Include your environment details** (OS, Go version, terminal)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List any alternatives you've considered**

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code, add tests
3. Ensure the test suite passes
4. Make sure your code follows the existing style
5. Write a clear commit message
6. Create a pull request!

## Development Setup

```bash
# Clone your fork
git clone https://github.com/CodeOne45/vex-tui.git
cd vex-tui

# Add upstream remote
git remote add upstream https://github.com/vex-tui/vex-tui.git
# Install dependencies
go mod download

# Build
go build -o vex-tui .

# Run
./vex sample_data.csv
```

## Project Structure

```
vex/
├── main.go                 # Entry point
├── internal/               # Private application code
│   ├── app/               # Application logic
│   ├── loader/            # File operations
│   ├── theme/             # Theme management
│   └── ui/                # UI utilities
└── pkg/                   # Public packages
    └── models/            # Data models
```

## Coding Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `go fmt` before committing
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names
- Add comments for exported functions

### Code Organization

- Keep functions small and focused
- Separate concerns (UI, logic, data)
- Use interfaces for testability
- Handle errors explicitly

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...
```

### Commit Messages

Follow the conventional commits specification:

```
feat: add new theme
fix: resolve search crash
docs: update README
style: format code
refactor: reorganize UI package
test: add loader tests
chore: update dependencies
```

## Areas for Contribution

### Good First Issues

- Documentation improvements
- Adding new themes
- UI polish and refinements
- Performance optimizations
- Test coverage

### Feature Ideas

- Column resizing
- Cell editing
- Filter capabilities
- Custom keybindings
- Configuration file support
- Chart visualization
- Diff mode

## Testing Your Changes

1. Test with various file formats (.xlsx, .csv)
2. Test with large files (50k+ rows)
3. Test all themes
4. Test all keyboard shortcuts
5. Test on different terminals (iTerm, Terminal.app, alacritty, etc.)
6. Test on different OS (macOS, Linux, Windows)

## Documentation

When adding new features:

1. Update README.md
2. Add examples
3. Update keyboard shortcuts section
4. Add inline code comments
5. Consider adding a tutorial

## Performance Considerations

- Use lazy loading for large files
- Minimize allocations in hot paths
- Profile before optimizing
- Consider memory usage for large datasets

## Questions?

Feel free to:

- Open an issue with your question
- Join discussions in existing issues
- Reach out to maintainers

## Recognition

Contributors will be:

- Listed in the README
- Mentioned in release notes
- Credited in commit history

Thank you for contributing! 🚀
