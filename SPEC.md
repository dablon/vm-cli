# vm-cli Specification

## Overview
CLI tool para ejecutar comandos en VMs remotas via SSH.

## Funcionalidades

### 1. Connect
Conectar a VM y ejecutar comando de prueba.
- Flags: `--host`, `--user`, `--password`, `--port`
- Output: información del sistema remoto

### 2. Exec
Ejecutar comando arbitrario en VM.
- Flags: `--host`, `--user`, `--password`, `--port`, `--command`
- Output: stdout del comando

### 3. User-Create
Crear usuario nuevo en VM con clave SSH.
- Flags: `--host`, `--user`, `--password`, `--new-user`, `--new-password`
- Output: confirmación + SSH public key

### 4. User-Exists
Verificar si usuario existe.
- Flags: `--host`, `--user`, `--password`, `--check-user`
- Output: boolean

### 5. Init
Guardar configuración local.
- Flags: `--host`, `--user`, `--agent`
- Output: path de config

## Arquitectura

```
vm-cli/
├── main.go           # Entry point + CLI flags
├── internal/
│   ├── ssh/
│   │   └── client.go # SSH client implementation
│   └── config/
│       └── config.go # Config management
├── go.mod
└── README.md
```

## Estado Actual

| Componente | Estado |
|------------|--------|
| Connect | ✅ Funcionando |
| Exec | ✅ Funcionando |
| User-Create | ✅ Funcionando |
| User-Exists | ✅ Funcionando |
| Init | ✅ Funcionando |
| Tests | ❌ Faltan |
| SSH Keys auth | ⚠️ Parcial |
| Config file | ⚠️ Básico |

## Tech Stack

- **Language:** Go 1.21
- **CLI Framework:** urfave/cli/v2
- **SSH:** golang.org/x/crypto/ssh

## Roadmap

- [ ] Agregar autenticación por SSH keys
- [ ] Agregar tests unitarios
- [ ] Mejorar manejo de errores
- [ ] Agregar comando `docker` integrado
- [ ] Soporte para múltiples VMs
- [ ] Configuración via YAML
