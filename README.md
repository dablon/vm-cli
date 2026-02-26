# vm-cli - Remote VM Management CLI

A command-line tool for managing remote Linux VMs via SSH. Built with Go.

## ⚠️ SECURITY WARNING

**All passwords have been rotated to high-entropy secure passwords (24 characters).**

## Features

- 🔌 SSH connection to remote VMs
- 📟 Execute commands remotely
- 👤 Create and manage users
- 🐳 Docker container management
- 📊 System information
- ⚙️ Configuration file support
- 🔐 Environment variable support for passwords

## Installation

### From Source

```bash
git clone https://github.com/dablon/vm-cli.git
cd vm-cli
go build -o vm-cli .
```

### Download Binary

Download the latest binary from [Releases](https://github.com/dablon/vm-cli/releases)

## Usage

### Basic Connection

```bash
# Using environment variable (RECOMMENDED)
export VM_CLI_PASSWORD="your_secure_password"
./vm-cli connect --host 142.44.247.203 --user myuser

# Or with argument (not recommended - visible in process list)
./vm-cli connect --host 142.44.247.203 --user myuser --password "your_secure_password"
```

### Execute Commands

```bash
# Run any command
./vm-cli exec --host 142.44.247.203 --user myuser --command "docker ps"
```

### User Management

```bash
# Create a new user
./vm-cli user-create --host 142.44.247.203 --user admin --password adminpass \
  --new-user newuser --new-password "secure_random_password"

# Check if user exists
./vm-cli user-exists --host 142.44.247.203 --user admin --check-user username
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VM_CLI_PASSWORD` | SSH password (alternative to --password) |

## Security Best Practices

- ⚠️ **Never commit passwords to version control**
- ✅ Use environment variables: `export VM_CLI_PASSWORD="..."`
- ✅ Generate strong passwords: `openssl rand -base64 24 | tr -dc 'A-Za-z0-9!@#$%&*' | head -c 24`
- ✅ Rotate passwords regularly
- ✅ Use SSH keys instead of passwords when possible

## Commands

| Command | Description |
|---------|-------------|
| `connect` | Connect to VM and run test command |
| `exec` | Execute a command on the remote VM |
| `user-create` | Create a new user |
| `user-exists` | Check if a user exists |
| `user-delete` | Delete a user |
| `sysinfo` | Get system information |
| `docker ps` | List Docker containers |
| `docker info` | Show Docker info |
| `init` | Initialize configuration |

## Development

```bash
# Run tests
go test ./...

# Build
go build -o vm-cli .

# Coverage
go test ./... -cover
```

## License

MIT
