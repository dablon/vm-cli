# vm-cli

CLI para ejecutar comandos en VMs Linux remotas via SSH.

## Características

- 🔌 Conexión SSH segura
- 📟 Ejecución de comandos remotos
- 👤 Gestión de usuarios (crear, eliminar, verificar)
- 🐳 Comandos Docker integrados
- 📊 Información del sistema
- 🔐 Soporte para variable de entorno `VM_CLI_PASSWORD`
- ✅ Código estructurado y testeable

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/dablon/vm-cli.git
cd vm-cli

# Compilar
go build -o vm-cli .

# O usar el binario precompilado
./vm-cli
```

## Uso

### Conexión básica

```bash
# Usando variable de entorno (RECOMENDADO)
export VM_CLI_PASSWORD="tu_contraseña"
./vm-cli connect --host 192.168.1.100 --user mi_usuario

# O con argumento
./vm-cli connect --host 192.168.1.100 --user mi_usuario --password "tu_contraseña"
```

### Ejecutar comandos

```bash
./vm-cli exec --host 192.168.1.100 --user mi_usuario --command "ls -la"
```

### Gestión de usuarios

```bash
# Crear usuario
./vm-cli user-create --host 192.168.1.100 --user admin --new-user nuevo_usuario --new-password "pass123"

# Verificar si existe
./vm-cli user-exists --host 192.168.1.100 --user admin --check-user nuevo_usuario

# Eliminar usuario
./vm-cli user-delete --host 192.168.1.100 --user admin --username nuevo_usuario
```

### Docker

```bash
# Listar contenedores
./vm-cli docker ps --host 192.168.1.100 --user mi_usuario

# Información de Docker
./vm-cli docker info --host 192.168.1.100 --user mi_usuario
```

### Información del sistema

```bash
./vm-cli sysinfo --host 192.168.1.100 --user mi_usuario
```

### Inicializar configuración

```bash
./vm-cli init --host 192.168.1.100 --user mi_usuario --agent mi_agente
```

## Comandos disponibles

| Comando | Descripción |
|---------|-------------|
| `connect` | Conectar a la VM y ejecutar comando de prueba |
| `exec` | Ejecutar un comando en la VM |
| `user-create` | Crear un nuevo usuario |
| `user-exists` | Verificar si un usuario existe |
| `user-delete` | Eliminar un usuario |
| `sysinfo` | Obtener información del sistema |
| `docker ps` | Listar contenedores Docker |
| `docker info` | Información de Docker |
| `init` | Inicializar configuración |

## Variables de entorno

| Variable | Descripción |
|----------|-------------|
| `VM_CLI_PASSWORD` | Contraseña SSH (alternativa a --password) |

## Seguridad

⚠️ **Recomendaciones:**
- Usa la variable de entorno `VM_CLI_PASSWORD` en lugar de pasar la contraseña como argumento
- Las contraseñas nunca se almacenan en disco
- Soporte para claves SSH

## Desarrollo

```bash
# Estructura del proyecto
vm-cli/
├── cmd/              # Comandos CLI
├── internal/
│   └── ssh/          # Cliente SSH
├── main.go           # Punto de entrada
└── go.mod            # Dependencias
```

## Licencia

MIT
