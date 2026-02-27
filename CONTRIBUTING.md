# Contributing to vm-cli

Thank you for your interest in contributing to vm-cli! This document provides guidelines for contributing.

## 🤝 How to Contribute

### Reporting Bugs

1. **Search existing issues** to avoid duplicates
2. **Create a new issue** with:
   - Clear title describing the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS

### Suggesting Features

1. **Open a discussion** before creating a feature request
2. **Explain the use case** - why do you need this feature?
3. **Provide examples** of how it should work

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature
   ```
3. **Make your changes** with clear commits
4. **Write tests** for new functionality
5. **Submit a Pull Request** with:
   - Clear description
   - Related issue number
   - Screenshots (if UI changes)

## 🛠️ Development Setup

```bash
# Clone repository
git clone https://github.com/dablon/vm-cli.git
cd vm-cli

# Install dependencies
go mod download

# Run tests
go test ./...

# Build locally
go build -o vm-cli .
```

## 📝 Coding Standards

- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Write tests for new features

## 📂 Project Structure

```
vm-cli/
├── cmd/           # CLI command implementations
├── internal/      # Internal packages
│   ├── config/   # Configuration management
│   └── ssh/     # SSH client implementation
├── .github/      # GitHub workflows
├── main.go       # Application entry point
├── README.md     # Documentation
└── LICENSE       # MIT License
```

## ❓ Questions?

- Open a GitHub Discussion
- Check existing issues and discussions

---

<p align="center">Your contributions make a difference. Thank you! 🚀</p>
