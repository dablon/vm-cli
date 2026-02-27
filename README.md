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
- ⚙️ **Profile system** - save VM credentials and reuse with --profile
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

## Quick Start (with Profiles)

```bash
# 1. Save your first profile
./vm-cli profile-save --name myvm --host 192.168.1.100 --user admin --password secret123

# 2. Use the profile (no need to type credentials again!)
./vm-cli --profile myvm exec --command "docker ps"

# Or use as command flag
./vm-cli exec --profile myvm --command "uptime"
```

## Profile Management

```bash
# Save a new profile
vm-cli profile-save --name ovh --host 142.44.247.203 --user nalcaraz --password your_password

# List all saved profiles
vm-cli profile-list

# Delete a profile
vm-cli profile-delete --name oldvm
```

## Usage

### Using Profiles (Recommended)

```bash
# Save profile once, use forever
vm-cli profile-save --name production --host 142.44.247.203 --user nalcaraz --password "your_secure_password"

# Execute commands using profile
vm-cli --profile production exec --command "docker ps"
vm-cli --profile production exec --command "df -h"
vm-cli --profile production sysinfo
```

### Using Environment Variables

```bash
# Using environment variable
export VM_CLI_PASSWORD="your_secure_password"
vm-cli exec --host 142.44.247.203 --user nalcaraz --command "docker ps"
```

### Using Flags Directly

```bash
# With explicit credentials (not recommended - visible in process list)
vm-cli exec --host 142.44.247.203 --user nalcaraz --password "your_password" --command "uptime"
```

## Commands

| Command | Description |
|---------|-------------|
| `profile-save` | Save a VM profile for quick access |
| `profile-list` | List all saved profiles |
| `profile-delete` | Delete a saved profile |
| `connect` | Connect to VM and run test command |
| `exec` | Execute a command on the remote VM |
| `user-create` | Create a new user |
| `user-exists` | Check if a user exists |
| `user-delete` | Delete a user |
| `sysinfo` | Get system information (uptime, CPU, memory, disk) |
| `docker` | Docker management commands |
| `init` | Initialize configuration directory |

### Docker Commands

```bash
# List containers
vm-cli --profile myvm docker ps

# Get container logs
vm-cli --profile myvm docker logs --container myapp --lines 50
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VM_CLI_PASSWORD` | SSH password (alternative to --password) |

## Security Best Practices

- ⚠️ **Never commit passwords to version control**
- ✅ Use profiles: `vm-cli profile-save --name vm1 ...`
- ✅ Use environment variables: `export VM_CLI_PASSWORD="..."`
- ✅ Generate strong passwords: `openssl rand -base64 24 | tr -dc 'A-Za-z0-9!@#$%&*' | head -c 24`
- ✅ Rotate passwords regularly
- ✅ Use SSH keys instead of passwords when possible
- ✅ Profiles are stored in `~/.vm-cli/profiles.json` - keep this file secure!

## Profile Storage

Profiles are saved in: `~/.vm-cli/profiles.json`

```json
{
  "ovh": {
    "name": "ovh",
    "host": "142.44.247.203",
    "user": "nalcaraz",
    "password": "encrypted_or_plaintext",
    "port": "22"
  }
}
```

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
