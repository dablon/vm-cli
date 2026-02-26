# vm-cli

CLI para gestionar máquinas virtuales remotas via SSH.

## Uso

### Conectar a VM
```bash
./vm-cli connect --host 142.44.247.203 --user reviewer --password reviewer123
```

### Ejecutar comando
```bash
./vm-cli exec --host 142.44.247.203 --user reviewer --password reviewer123 --command "docker ps"
```

### Crear usuario
```bash
./vm-cli user-create --host 142.44.247.203 --user admin --password pass --new-user coder --new-password coder123
```

### Verificar usuario
```bash
./vm-cli user-exists --host 142.44.247.203 --user admin --password pass --check-user reviewer
```

## Configuración

Usar variable de entorno para contraseña:
```bash
export VM_CLI_PASSWORD="tu-contraseña"
./vm-cli connect --host 142.44.247.203 --user reviewer
```

## Requisitos

- Go 1.21+
- Dependencias: ver `go.mod`

## Build

```bash
go build -o vm-cli .
```

## Seguridad

⚠️ **No usar en producción:**
- Contraseñas en texto plano (usar SSH keys)
- HostKeyCallback ignorado (`InsecureIgnoreHostKey`)
- Necesita mejoras de seguridad antes de uso en producción
