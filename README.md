# Auth-Go API

Autenticación y gestión de sesiones multi-dispositivo para aplicaciones móviles y web, desarrollada en Go usando Fiber y GORM.

---

## Características principales
- **Registro y autenticación de usuarios** con email y contraseña.
- **Tokens JWT**: Acceso y Refresh, con expiración y secretos configurables.
- **Gestión de sesiones por dispositivo** (`device_id`), ideal para apps móviles.
- **Revocación de sesiones** individuales y logout global.
- **Límites configurables de sesiones activas por usuario**.
- **Pruebas QA exhaustivas** (colección Postman incluida).
- **Seguridad**: Hash de contraseñas, validaciones, manejo de errores.

---

## Tecnologías
- Go 1.20+
- [Fiber](https://gofiber.io/) (web framework)
- [GORM](https://gorm.io/) (ORM)
- JWT (github.com/golang-jwt/jwt)
- MySQL/MariaDB

---

# Crea la base de datos en MySQL

```sql
CREATE DATABASE IF NOT EXISTS test_db;
```

---

## Configuración
1. **Clona el repositorio y entra al directorio:**
   ```bash
   git clone https://github.com/DavidDevGt/auth-go
   cd auth-go
   ```

2. **Configura el entorno:**
   - Copia `.env.example` a `.env` y ajusta los valores:
     ```env
     PORT=8080
     DATABASE_URL=admin:password@tcp(localhost:3306)/test_db?parseTime=true&loc=Local
     ACCESS_TOKEN_SECRET=secretthings_ok
     REFRESH_TOKEN_SECRET=secretthings_ok
     ACCESS_TOKEN_EXPIRY=30m
     REFRESH_TOKEN_EXPIRY=960h
     MAX_SESSIONS_PER_USER=2
     ```

3. **Prepara la base de datos:**
   - Crea la base de datos `test_db` en tu MySQL/MariaDB.
   - Las migraciones automáticas se ejecutan al iniciar el servicio.

4. **Instala dependencias y ejecuta:**
   ```bash
   go mod tidy
   go run main.go
   ```

---

## Endpoints principales

### Liveness
```http
GET /healthz
```

### Readiness
```http
GET /readyz
```

### Registro
```http
POST /api/register
{
  "name": "Nombre",
  "email": "usuario@mail.com",
  "password": "PasswordFuerte123!"
}
```

### Login
```http
POST /api/login
{
  "email": "usuario@mail.com",
  "password": "PasswordFuerte123!",
  "device_id": "uuid-o-identificador-unico-del-dispositivo"
}
```

### Refresh Token
```http
POST /api/refresh
{
  "refresh_token": "<refresh_token>",
  "device_id": "uuid-o-identificador-unico-del-dispositivo"
}
```

### Logout (revoca el refresh_token)
```http
POST /api/logout
{
  "refresh_token": "<refresh_token>"
}
```

### Revocar sesión específica (por refresh_token)
```http
POST /api/revoke-session
{
  "refresh_token": "<refresh_token>"
}
```

### Acceso protegido
```http
GET /api/protected
Authorization: Bearer <access_token>
```

---

## Kubernetes

Los manifiestos de Kubernetes se encuentran en el directorio `k8s/`:

### Despliegue en Kubernetes
```bash
# Aplicar ConfigMap
kubectl apply -f k8s/configmap.yaml

# Aplicar Deployment y Service
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Verificar el estado
kubectl get pods -l app=auth-go
kubectl get svc auth-go
```

**Nota**: Ajusta los valores en `configmap.yaml` según tu entorno.

---

## Pruebas QA
- Incluye una colección Postman profesional (`test-postman.json`) con flujos positivos y negativos, datos dinámicos y validaciones automáticas.
- **Recomendado:** Ejecutar la suite con delay entre peticiones para evitar sobrecargar el servidor.

---

## Seguridad y mejores prácticas
- Usa secretos únicos y fuertes en `.env`.
- Limita el número de sesiones por usuario.
- Implementa HTTPS en producción.
- Revisa y ajusta los CORS según tu frontend.

---

## Extensión y personalización
- Puedes agregar campos personalizados al modelo de usuario (`internal/database/models/user.go`).
- Para notificaciones, MFA o integración con OAuth, extiende los servicios actuales y los modelos.
- Los servicios y modelos están diseñados para ser testeables de forma fácil.

---

## Estructura del proyecto

```
.
├── internal
│   ├── config         # Configuración de entorno
│   ├── database
│   │   └── models     # Modelos GORM: User, Session
│   ├── middleware     # Middlewares de autenticación
│   ├── routes         # Definición de endpoints
│   ├── services       # Lógica de negocio: Auth, Session, User
│   └── utils          # Utilidades: tokens, helpers
├── main.go            # Entry point
├── test-postman.json  # Colección QA
├── .env, .env.example # Configuración de entorno
└── README.md
```

---

## Licencia
MIT
