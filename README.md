# vm-cli 🖥️

> Lightweight CLI for remote Linux VM management via SSH. Built with Go.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dablon/vm-cli)](https://github.com/dablon/vm-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/dablon/vm-cli)](https://github.com/dablon/vm-cli/releases/latest)

A fast, dependency-free CLI tool for managing remote Linux VMs over SSH. Execute commands, manage Docker containers, transfer files, search for files - all with a single, intuitive command-line interface.

## ✨ Features

- 🔌 **SSH Connection** - Connect to remote VMs effortlessly
- 📟 **Remote Execution** - Run commands on remote servers
- 🐳 **Docker Management** - List containers, view logs, manage services
- 👤 **User Management** - Create, check, and delete remote users
- 💾 **Profile System** - Save and reuse VM credentials securely
- 📊 **System Info** - Quick access to system resources (CPU, RAM, disk)
- 📁 **File Transfer** - Copy files between local and remote
- 🔍 **File Search** - Find files and folders on remote VM
- 🔐 **Security** - Support for environment variables and secure credential handling

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/dablon/vm-cli.git
cd vm-cli

# Build the binary
go build -o vm-cli .

# Optional: Add to PATH
sudo mv vm-cli /usr/local/bin/
```

### First Run - Save a Profile

```bash
# Save your first VM profile (do once)
vm-cli profile-save \
  --name production \
  --host 192.168.1.100 \
  --user admin \
  --password your_password

# Or use a custom SSH port
vm-cli profile-save \
  --name myserver \
  --host 192.168.1.100 \
  --user root \
  --password secret \
  --port 2222
```

### Execute Commands

```bash
# Using saved profile (no credentials needed!)
vm-cli exec --profile production --command "uptime"

# List Docker containers
vm-cli exec --profile production --command "docker ps"

# Get system information
vm-cli exec --profile production --command "df -h"
```

## 📖 Usage

### Profile Management

```bash
# Save a new profile
vm-cli profile-save --name myvm --host 192.168.1.100 --user admin --password secret

# List all saved profiles
vm-cli profile-list

# Delete a profile
vm-cli profile-delete --name oldvm
```

### Command Execution

```bash
# Basic command execution (profile flag must come AFTER exec!)
vm-cli exec --profile myvm --command "df -h"

# With inline credentials (not recommended)
vm-cli exec --host 192.168.1.100 --user admin --password secret --command "docker ps"

# Using environment variable
export VM_CLI_PROFILE=myvm
vm-cli exec --command "uptime"
```

### Docker Management

```bash
# List containers
vm-cli exec --profile myvm --command "docker ps"

# View container logs
vm-cli exec --profile myvm --command "docker logs myapp --tail 50"
```

### File Transfer

```bash
# Copy file FROM local TO remote
vm-cli copy --profile myvm --source=./localfile.txt --dest=/tmp/remotefile.txt --to-remote

# Copy file FROM remote TO local
vm-cli copy --profile myvm --source=/tmp/remotefile.txt --dest=./localfile.txt
```

### File Search

```bash
# Search for files by name pattern
vm-cli find --profile myvm --name "*.json"

# Search for directories
vm-cli find --profile myvm --name "crypto*" --type d

# In specific directory with depth
vm-cli find --profile myvm --directory /home --name "*.log" --maxdepth 3

# Include hidden files
vm-cli find --profile myvm --name ".env*" --hidden
```

### User Management

```bash
# Create a new user
vm-cli exec --profile myvm --command "sudo useradd -m developer"

# Check if user exists
vm-cli exec --profile myvm --command "id developer"

# Delete a user
vm-cli exec --profile myvm --command "sudo userdel developer"
```

### System Information

```bash
# Get system stats (CPU, Memory, Disk)
vm-cli exec --profile myvm --command "echo CPU && lscpu | grep 'Model name' && echo MEMORY && free -h && echo DISK && df -h"
```

## 🔐 Security

> **Important:** Never commit passwords or sensitive data to version control.

### Recommended Practices

1. **Use Profiles** - Store credentials locally with `profile-save`
2. **Environment Variables** - Use `VM_CLI_PASSWORD` or `VM_CLI_PROFILE` instead of CLI flags
3. **SSH Keys** - Prefer SSH key-based authentication when possible
4. **Rotate Passwords** - Regularly update stored credentials

```bash
# Using environment variable
export VM_CLI_PASSWORD="your_secure_password"
vm-cli exec --profile myvm --command "whoami"
```

Profile data is stored in `~/.vm-cli/profiles.json` - ensure this file has appropriate permissions:

```bash
chmod 600 ~/.vm-cli/profiles.json
```

## 🛠️ Development

```bash
# Clone and setup
git clone https://github.com/dablon/vm-cli.git
cd vm-cli

# Run tests
go test ./...

# Build
go build -o vm-cli .

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o vm-cli-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o vm-cli-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o vm-cli.exe .
```

## 📝 Commands Reference

| Command | Alias | Description |
|---------|-------|-------------|
| `profile-save` | - | Save a VM profile |
| `profile-list` | - | List saved profiles |
| `profile-delete` | - | Delete a profile |
| `exec` | - | Execute a command on remote VM |
| `connect` | - | Test SSH connection |
| `user-create` | - | Create a remote user |
| `user-exists` | - | Check if user exists |
| `user-delete` | - | Delete a remote user |
| `sysinfo` | - | Get system information |
| `docker` | - | Docker management (ps, logs) |
| `init` | - | Initialize config directory |
| `copy` | - | Copy files between local and remote |
| `find` | - | Search for files/folders on remote VM |

### Find Options

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--name` | `-n` | (required) | File/folder name pattern (* and ? supported) |
| `--directory` | - | `/` | Directory to search |
| `--type` | `-t` | `f` | Type: `f` (file), `d` (directory), `b` (both) |
| `--maxdepth` | - | 3 | Maximum directory depth |
| `--hidden` | - | false | Include hidden files |

### Copy Options

| Flag | Alias | Description |
|------|-------|-------------|
| `--source` | `-s` | Source file path |
| `--dest` | `-d` | Destination file path |
| `--to-remote` | `-to` | Copy to remote (default is from remote) |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - CLI framework for Go
- [ssh](golang.org/x/crypto/ssh) - SSH implementation for Go

---

<p align="center">Made with ❤️ for remote server management</p>
