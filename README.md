<div align="center">

  <img src="assets/imgs/logo.png" alt="Medical Exams API" width="220" />

  # Medical Exams Backend

  API para gestion de pacientes, ordenes y resultados de laboratorio.

</div>

## Contenido

- [Resumen](#resumen)
- [Requisitos](#requisitos)
- [Configuracion](#configuracion)
- [Ejecucion](#ejecucion)
- [Documentacion Swagger](#documentacion-swagger)
- [Autenticacion](#autenticacion)
- [Endpoints](#endpoints)
- [Pruebas](#pruebas)
- [Seguridad y buenas practicas](#seguridad-y-buenas-practicas)

## Resumen

Backend en Go con Gin y GORM para manejo de:

- autenticacion con JWT
- pacientes
- ordenes de laboratorio
- validacion y carga de resultados

Incluye migraciones y seeding controlado por entorno.

## Requisitos

- Go 1.25.6
- PostgreSQL

## Configuracion

Crea un archivo `.env` basado en estos valores (ejemplo):

```env
# Server configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=lab_admin
DB_PASSWORD=CHANGE_ME
DB_NAME=medical_exams_db
DB_SSLMODE=disable

# JWT configuration
JWT_SECRET=CHANGE_ME_TO_A_STRONG_SECRET
JWT_EXPIRATION_HOURS=24

# CORS configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOW_CREDENTIALS=true

# Rate limiting (login/register)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_PER_MINUTE=120
RATE_LIMIT_BURST=30

# Proxy trust (release)
TRUSTED_PROXIES=127.0.0.1,::1

# Seeding configuration
SEED_DB=true
```

Notas:

- En `GIN_MODE=release` se requiere `CORS_ALLOWED_ORIGINS`.
- `SEED_DB` por defecto se ejecuta en dev y se omite en release.

## Ejecucion

```bash
go run .
```

El servidor quedara en `http://localhost:8080`.

## Documentacion Swagger

Disponible en:

```
http://localhost:8080/swagger/index.html
```

## Autenticacion

La API usa JWT. Flujo recomendado:

1. `POST /api/v1/login` para obtener el token.
2. Enviar `Authorization: Bearer <token>` en rutas protegidas.

Ejemplo:

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin123!"}'
```

## Endpoints

Base path: `/api/v1`

### Publicos

- `POST /login`
- `POST /register`

### Protegidos

- `GET /me`

#### Pacientes

- `POST /patients`
- `GET /patients`
- `GET /patients/:id`

#### Ordenes

- `POST /orders`
- `GET /orders`

#### Laboratorio

- `GET /lab/exams/:id`
- `PATCH /lab/exams/:id/status`
- `POST /lab/exams/:id/validate`
- `POST /lab/exams/:id/results`
- `GET /lab/exams/catalog`

Ejemplo (GET pacientes):

```bash
curl http://localhost:8080/api/v1/patients \
  -H "Authorization: Bearer <token>"
```

## Ejemplos de algunos payloads

Base path: `/api/v1`

### Auth

**Login**

```json
{
  "username": "admin",
  "password": "Admin123!"
}
```

**Register**

```json
{
  "username": "user1",
  "email": "user1@test.com",
  "password": "Admin123!",
  "full_name": "User One",
  "role_id": 2
}
```

### Pacientes

**POST /patients**

```json
{
  "document_type": "cedula",
  "document_number": "V12345678",
  "first_name": "Maria",
  "last_name": "Delgado",
  "date_of_birth": "1990-05-15T00:00:00Z",
  "gender": "F",
  "phone": "04125555555",
  "email": "maria@test.com",
  "blood_type": "O+"
}
```

### Ordenes

**POST /orders**

```json
{
  "patient_id": 1,
  "priority": "normal",
  "referring_doctor": "Dr. Perez",
  "diagnosis": "Chequeo general",
  "exams": [
    {
      "exam_type_id": 1,
      "price": 10.00
    }
  ]
}
```

### Laboratorio

**PATCH /lab/exams/:id/status**

```json
{
  "status": "muestra_tomada"
}
```

**POST /lab/exams/:id/results**

```json
[
  {
    "exam_parameter_id": 1,
    "value_numeric": 5.6,
    "value_text": ""
  }
]
```

**POST /lab/exams/:id/validate**

No requiere body.

## Errores y respuestas

Formato estandar:

```json
{
  "status": "success",
  "code": 200,
  "message": "Mensaje",
  "data": {}
}
```

```json
{
  "status": "error",
  "code": 400,
  "message": "Mensaje",
  "errors": "detalle"
}
```

Tabla de codigos mas comunes:

| Codigo | Descripcion | Uso tipico |
| --- | --- | --- |
| 200 | OK | Respuesta exitosa |
| 201 | Created | Registro creado |
| 400 | Bad Request | Validacion o payload invalido |
| 401 | Unauthorized | Token faltante o invalido |
| 404 | Not Found | Recurso inexistente |
| 429 | Too Many Requests | Rate limit excedido |
| 500 | Internal Server Error | Error inesperado |

## Pruebas

```bash
go test ./...
```

Incluye smoke tests con `httptest` para auth, pacientes, ordenes y catalogo de examenes.

## Seguridad y buenas practicas

- JWT con secreto por entorno y verificacion de algoritmo.
- Rate limiting en `/login` y `/register`.
- Headers de seguridad agregados via middleware.
- CORS configurable por entorno.
- Seeding deshabilitado en `release` por defecto.

Si necesitas hardening adicional (TLS, audit logs, WAF o redis rate limit), abre un issue o solicita un ajuste.
