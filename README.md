# vm-cli - Remote VM Management CLI

A command-line tool for managing remote Linux VMs via SSH. Built with Go.

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
export VM_CLI_PASSWORD="your_password"
./vm-cli connect --host 142.44.247.203 --user myuser

# Or with argument
./vm-cli connect --host 142.44.247.203 --user myuser --password mypass
```

### Execute Commands

```bash
# Run any command
./vm-cli exec --host 142.44.247.203 --user myuser --password mypass --command "docker ps"

# List all containers (including stopped)
./vm-cli exec --host 142.44.247.203 --user myuser --command "docker ps -a"
```

### User Management

```bash
# Create a new user
./vm-cli user-create --host 142.44.247.203 --user admin --password adminpass --new-user newuser --new-password newpass

# Check if user exists
./vm-cli user-exists --host 142.44.247.203 --user admin --password adminpass --check-user username

# Delete a user
./vm-cli user-delete --host 142.44.247.203 --user admin --password adminpass --username username
```

### Docker Management

```bash
# List containers
./vm-cli docker ps --host 142.44.247.203 --user myuser --password mypass

# Show Docker info
./vm-cli docker info --host 142.44.247.203 --user myuser --password mypass
```

### System Information

```bash
./vm-cli sysinfo --host 142.44.247.203 --user myuser --password mypass
```

### Configuration

```bash
# Initialize config file
./vm-cli init --host 142.44.247.203 --user myuser --password mypass --agent myagent
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VM_CLI_PASSWORD` | SSH password (alternative to --password) |

## Commands

| Command | Description |
|---------|-------------|
| `connect` | Connect to VM and run test command |
| `exec` | Execute a command on the remote VM |
| `user-create` | Create a new user on the VM |
| `user-exists` | Check if a user exists |
| `user-delete` | Delete a user from the VM |
| `sysinfo` | Get system information |
| `docker ps` | List Docker containers |
| `docker info` | Show Docker information |
| `init` | Initialize configuration file |

## Examples

### Full Workflow

```bash
# 1. Connect and check system
./vm-cli connect --host 142.44.247.203 --user nalcaraz --password "mypassword"

# 2. Create a user for your agent
./vm-cli user-create --host 142.44.247.203 --user nalcaraz --password "mypassword" \
  --new-user agent_001 --new-password "agent123"

# 3. Use the new user
./vm-cli exec --host 142.44.247.203 --user agent_001 --password "agent123" \
  --command "docker run hello-world"

# 4. List containers
./vm-cli docker ps --host 142.44.247.203 --user agent_001 --password "agent123"
```

## Security Notes

- ⚠️ Never commit passwords to version control
- Use environment variables: `export VM_CLI_PASSWORD="..."`
- Consider using SSH keys instead of passwords
- The `--password` flag is visible in process list

## Development

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build -o vm-cli .
```

## Project Structure

```
vm-cli/
├── cmd/
│   └── commands.go       # CLI commands
├── internal/
│   └── ssh/
│       ├── client.go    # SSH client
│       └── client_test.go
├── main.go              # Entry point
├── go.mod               # Dependencies
└── README.md            # This file
```

## License

MIT
