# prabodrive — Backend API

Platform manajemen dokumen berbasis cloud. Tugas akhir mata kuliah Cloud Computing.

---

## ⚠️ SEBELUM MENULIS SATU BARIS PUN — BACA INI DULU

**Langkah pertama wajib: periksa template yang sudah ada.**

```bash
# 1. Lihat seluruh struktur project
find . -type f | sort

# 2. Baca SEMUA file yang sudah ada sebelum menulis apapun
# Tidak ada pengecualian — termasuk go.mod, main.go, config, utils, dll.
```

Setelah membaca semua file yang ada:

- **Ikuti konvensi yang sudah ada** — naming, error handling style, package structure, komentar. Jangan ganti jika sudah ada pola yang konsisten.
- **Jangan timpa file yang sudah ada** kecuali ada konflik nyata yang harus diselesaikan. Jika file sudah ada, extend — jangan recreate.
- **Sesuaikan Build Order di bawah** dengan kondisi riil — skip langkah yang sudah selesai, lanjut dari yang belum ada.
- **Jika template punya folder/package berbeda** dari yang didefinisikan di sini, ikuti template. Spesifikasi di CLAUDE.md adalah panduan fungsional, bukan kewajiban struktural.
- **Laporkan apa yang ditemukan** sebelum mulai: "File yang sudah ada: ..., File yang perlu dibuat: ..."

---

## Code Quality Standards

Setiap file yang ditulis harus memenuhi standar berikut. Tidak ada kompromi.

### Struktur & Desain

- **Handler tipis** — handler hanya parsing input, panggil service, kembalikan response. Tidak ada business logic di handler.
- **Satu fungsi, satu tanggung jawab** — jika fungsi butuh komentar untuk menjelaskan "bagian A" dan "bagian B", itu harusnya dua fungsi.
- **Return early** — validasi dan error check di atas, happy path di bawah. Hindari nesting dalam.
- **Tidak ada duplikasi** — jika logika muncul lebih dari satu kali, ekstrak ke fungsi atau helper.

### Error Handling

- **Tidak ada `_` untuk error** — setiap error wajib di-handle atau di-wrap dengan konteks: `fmt.Errorf("createFolder: %w", err)`.
- **Jangan panic di luar startup** — panic hanya boleh di `main()` untuk konfigurasi kritis (JWT secret kosong, DB tidak bisa connect). Semua error runtime dikembalikan.
- **Error dari DB** — bedakan `pgx.ErrNoRows` (→ 404) dari error lain (→ 500). Jangan semua dijadikan 500.

### Kode

- **Tidak ada magic number** — gunakan konstanta bernama: `const maxFileSize = 5 * 1024 * 1024` bukan `5242880` tersebar di mana-mana.
- **Nama bermakna** — tidak ada `data`, `tmp`, `x`, `res2`. Nama harus menjelaskan isi.
- **Tidak ada import yang tidak dipakai** — jalankan `goimports` atau pastikan manual.
- **Komentar hanya untuk WHY** — jangan komentari apa yang sudah jelas dari kode:
  ```go
  // ❌ Increment counter
  count++

  // ✓ Quota dihitung dalam bytes bukan KB karena S3 melaporkan dalam bytes
  const quotaUnit = 1
  ```
- **`defer rows.Close()`** setelah setiap query yang return rows, tanpa pengecualian.
- **Parameterized query wajib** — tidak ada string concatenation untuk SQL. Selalu gunakan `$1, $2, ...`.
- **Context diteruskan** — semua fungsi yang melakukan I/O (DB, S3, HTTP) menerima `context.Context` sebagai parameter pertama.

### Verifikasi Sebelum Selesai

```bash
go build ./...          # harus zero error
go vet ./...            # harus zero warning
```

Jika ada warning dari `go vet`, perbaiki dulu sebelum lanjut ke langkah berikutnya.

---

## Stack & Versions

- **Language**: Go 1.24
- **Framework**: Gin
- **Database**: PostgreSQL 18
- **Storage**: Amazon S3
- **Auth**: JWT (access 1h + refresh 7d) + bcrypt
- **Email**: AWS SES
- **Runtime**: Docker on AWS EC2, behind CloudFront `/api/*` behavior

## Dependencies (go.mod)

```
github.com/gin-gonic/gin v1.10.0
github.com/golang-jwt/jwt/v5 v5.2.1
github.com/aws/aws-sdk-go-v2 v1.30.0
github.com/aws/aws-sdk-go-v2/config v1.27.0
github.com/aws/aws-sdk-go-v2/service/s3 v1.58.0
github.com/aws/aws-sdk-go-v2/service/ses v1.25.0
github.com/jackc/pgx/v5 v5.6.0
golang.org/x/crypto v0.24.0
github.com/gabriel-vasile/mimetype v1.4.4
github.com/ulule/limiter/v3 v3.11.2
github.com/swaggo/gin-swagger v1.6.0
github.com/swaggo/swag v1.16.3
github.com/joho/godotenv v1.5.1
github.com/google/uuid v1.6.0
```

---

## Project Structure
disesuaikan dengan current folder structure di repo ini

## Environment Variables (.env.example)

```env
# Server
PORT=8080
GIN_MODE=release
MAINTENANCE_MODE=false

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=prabodrive
DB_USER=postgres
DB_PASSWORD=secret
DB_SSLMODE=disable

# JWT
JWT_SECRET=change-this-secret-min-32-chars
JWT_ACCESS_EXPIRY=1h
JWT_REFRESH_EXPIRY=168h

# AWS
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
S3_BUCKET=prabodrive-prod
S3_PRESIGN_EXPIRY=15m

# SES
SES_FROM_EMAIL=no-reply@prabodrive.com

# CloudFront
CLOUDFRONT_DOMAIN=https://dxxxxx.cloudfront.net
```

---

## Database Schema (001_initial.sql)

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email        VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    quota_used   BIGINT NOT NULL DEFAULT 0,
    quota_max    BIGINT NOT NULL DEFAULT 3221225472, -- 3 GB in bytes
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE folders (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id  UUID REFERENCES folders(id) ON DELETE SET NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE documents (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id  UUID REFERENCES folders(id) ON DELETE SET NULL,
    name       VARCHAR(255) NOT NULL,
    size       BIGINT NOT NULL,
    mime_type  VARCHAR(127) NOT NULL,
    s3_key     VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE share_links (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id   UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    token         VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    expires_at    TIMESTAMPTZ NOT NULL,
    created_by    UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE activity_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users(id) ON DELETE SET NULL,
    action      VARCHAR(50) NOT NULL,
    document_id UUID REFERENCES documents(id) ON DELETE SET NULL,
    ip_address  INET,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_folder_id ON documents(folder_id);
CREATE INDEX idx_folders_user_id ON folders(user_id);
CREATE INDEX idx_share_links_token ON share_links(token);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
```

---

## Middleware

### 1. Maintenance Mode (`middleware/maintenance.go`)

```go
// Baca MAINTENANCE_MODE dari env setiap request (bukan cache startup)
// Jika true: izinkan hanya GET /health, semua lain abort 503
// Response 503:
// {"status": "maintenance", "message": "Service is under maintenance. Please try again later."}
```

Daftarkan sebagai middleware global PERTAMA sebelum semua route group.

### 2. Auth (`middleware/auth.go`)

```go
// Ambil Bearer token dari header Authorization
// Validasi JWT, extract user_id dan email ke gin.Context
// Key context: "user_id" (string UUID), "user_email" (string)
// Jika invalid/expired: 401 {"error": "unauthorized"}
```

### 3. Rate Limiter (`middleware/ratelimit.go`)

```go
// 100 request per menit per IP
// Gunakan limiter/v3 dengan in-memory store
// Jika exceeded: 429 {"error": "too many requests"}
```

---

## API Endpoints

### Response Conventions

Dua aturan yang berlaku di seluruh API tanpa pengecualian:

**1. PATCH untuk semua update parsial** — jangan gunakan PUT. PUT mengimplikasikan replace seluruh resource. PATCH untuk update field tertentu saja.

**2. Response minimal untuk operasi non-GET** — jangan kembalikan full payload jika client tidak butuhnya:

| Operasi | Response body |
|---|---|
| `GET` | Full resource object |
| `POST` create | `{"id": "uuid"}` — kecuali ada data server-generated yang client butuh sekarang (token, URL) |
| `PATCH` update | `{"id": "uuid"}` |
| `DELETE` | `{"id": "uuid"}` |

Alasan: client sudah punya data yang dia kirim. Yang dia butuh hanya ID untuk referensi selanjutnya. Mengirim balik seluruh object adalah pemborosan bandwidth dan menyebabkan coupling antara server schema dan client state.

**Pengecualian yang sah** — response boleh menyertakan field tambahan jika field tersebut di-generate di server dan client tidak mungkin mengetahuinya:
- Auth login/register → `access_token`, `refresh_token` (wajib dikembalikan)
- Presign upload → `upload_url`, `s3_key` (client butuh untuk upload ke S3)
- Create share link → `token`, `share_url` (hanya ada di server, tidak bisa di-derive client)

---

### Health

```
GET  /health
```
Response: `{"status": "ok", "version": "1.0.0"}`
**Tidak butuh auth. Selalu aktif meski maintenance mode ON.**

---

### Auth (`/api/v1/auth`)

```
POST /api/v1/auth/register
Body: { "name": string, "email": string, "password": string (min 8 char) }
Response 201: { "id": "uuid", "access_token": string, "refresh_token": string }

POST /api/v1/auth/login
Body: { "email": string, "password": string }
Response 200: { "id": "uuid", "access_token": string, "refresh_token": string }

POST /api/v1/auth/refresh
Body: { "refresh_token": string }
Response 200: { "access_token": string, "refresh_token": string }

POST /api/v1/auth/logout        [AUTH REQUIRED]
Body: { "refresh_token": string }
Response 200: { "id": "uuid" }
```

---

### Documents (`/api/v1/documents`) — semua AUTH REQUIRED

```
GET    /api/v1/documents
       Query: folder_id (optional UUID), search (optional string), page (int), limit (int, max 50)
       Response 200: { "data": [Document], "total": int, "page": int }

GET    /api/v1/documents/:id
       Response 200: { "data": Document }

POST   /api/v1/documents/presign-upload
       Body: { "name": string, "size": int (bytes), "mime_type": string, "folder_id": UUID|null }
       Validasi SEBELUM issue URL:
         1. mime_type harus ada dalam ALLOWED_MIME_TYPES
         2. size <= 5242880 (5 MB)
         3. quota_used + size <= quota_max → jika tidak: 403 {"error": "quota exceeded"}
       Response 200: { "upload_url": string, "s3_key": string, "expires_at": ISO8601 }

POST   /api/v1/documents/confirm-upload
       Body: { "s3_key": string, "name": string, "size": int, "mime_type": string, "folder_id": UUID|null }
       Aksi: simpan metadata ke DB, UPDATE users SET quota_used = quota_used + $size WHERE id = $user_id
       Response 201: { "id": "uuid" }

DELETE /api/v1/documents/:id
       Aksi: hapus dari S3, hapus dari DB, UPDATE quota_used - size
       Response 200: { "id": "uuid" }

GET    /api/v1/documents/:id/download
       Generate presigned GET URL (expiry 15 menit)
       Response 200: { "url": string, "expires_at": ISO8601 }
```

Document shape (untuk GET):
```json
{ "id": "uuid", "name": "string", "size": 1234, "mime_type": "application/pdf",
  "folder_id": "uuid|null", "created_at": "ISO8601", "updated_at": "ISO8601" }
```

---

### Folders (`/api/v1/folders`) — semua AUTH REQUIRED

```
GET    /api/v1/folders
       Response 200: { "data": [Folder] }

GET    /api/v1/folders/:id
       Response 200: { "data": Folder }

POST   /api/v1/folders
       Body: { "name": string, "parent_id": UUID|null }
       Response 201: { "id": "uuid" }

PATCH  /api/v1/folders/:id
       Body: { "name": string }
       Response 200: { "id": "uuid" }

DELETE /api/v1/folders/:id
       Aksi: dokumen di dalamnya dipindah ke parent folder atau root
       Response 200: { "id": "uuid" }
```

---

### Share Links (`/api/v1/share`) — semua AUTH REQUIRED kecuali access

```
POST   /api/v1/share
       Body: { "document_id": UUID, "expires_at": ISO8601, "password": string|null }
       Aksi: generate random 32-byte token, simpan hash password jika ada
       Kirim notifikasi email via SES ke owner dokumen
       Response 201: { "id": "uuid", "token": string, "share_url": string, "expires_at": ISO8601 }
       # token dan share_url dikembalikan karena di-generate server dan client butuh sekarang

GET    /api/v1/share/:token            # PUBLIC, no auth
       Query: password (optional)
       Validasi: token valid, belum expired, password cocok (jika ada)
       Response 200: { "download_url": string, "expires_at": ISO8601 }

DELETE /api/v1/share/:id               # AUTH REQUIRED, must be owner
       Response 200: { "id": "uuid" }
```

---

### Activity Log (`/api/v1/activity`) — AUTH REQUIRED

```
GET /api/v1/activity
    Query: page (int), limit (int, max 50)
    Response 200: { "data": [ActivityLog], "total": int }
```

Log action types: `upload`, `download`, `share_create`, `share_access`, `delete`, `login`, `logout`

---

## File Validation (utils/mime.go)

**ALLOWED_MIME_TYPES** (whitelist):
```go
var AllowedMIMETypes = map[string]bool{
    "application/pdf": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
    "application/msword": true,                                                       // .doc
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,       // .xlsx
    "application/vnd.ms-excel": true,                                                 // .xls
    "application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
    "application/vnd.ms-powerpoint": true,                                            // .ppt
    "text/plain": true,                                                               // .txt
    "image/jpeg": true,
    "image/png":  true,
    "image/webp": true,
}
```

Fungsi `ValidateMIME(data []byte, declaredMIME string) error`:
1. Deteksi MIME dari magic bytes: `mimetype.Detect(data)`
2. Cek `detectedMIME` ada di `AllowedMIMETypes`
3. Cek `declaredMIME` ada di `AllowedMIMETypes`
4. Jika detected != declared: return error "MIME type mismatch"
5. Jika salah satu tidak di whitelist: return error "file type not allowed"

---

## S3 Operations (services/s3.go)

```go
// GeneratePresignedPutURL(s3Key string, mimeType string, expiry time.Duration) (string, error)
//   → pakai s3.PresignClient, PUT method, expiry dari env S3_PRESIGN_EXPIRY
//
// GeneratePresignedGetURL(s3Key string, expiry time.Duration) (string, error)
//   → GET method, expiry 15 menit
//
// DeleteObject(s3Key string) error
//
// S3 key format: {user_id}/{folder_id_or_"root"}/{uuid}_{sanitized_filename}
```

Bucket: selalu ambil dari env `S3_BUCKET`. Bucket bersifat private — akses hanya via presigned URL.

---

## Quota Service (services/quota.go)

```go
// CheckQuota(userID string, fileSize int64) error
//   SELECT quota_used, quota_max FROM users WHERE id = $1
//   IF quota_used + fileSize > quota_max → return ErrQuotaExceeded
//
// AddQuota(tx pgx.Tx, userID string, delta int64) error
//   UPDATE users SET quota_used = quota_used + $delta WHERE id = $1
//   Gunakan dalam transaksi DB yang sama dengan insert/delete dokumen
```

---

## Standard Response Helpers (utils/response.go)

```go
func OK(c *gin.Context, data any)              // 200 - untuk GET dengan full data
func Created(c *gin.Context, id string)        // 201 - kembalikan {"id": id}
func Updated(c *gin.Context, id string)        // 200 - kembalikan {"id": id}
func Deleted(c *gin.Context, id string)        // 200 - kembalikan {"id": id}
func Data(c *gin.Context, data any)            // 200 - untuk response dengan server-generated fields
func BadRequest(c *gin.Context, err string)    // 400
func Unauthorized(c *gin.Context)              // 401
func Forbidden(c *gin.Context, err string)     // 403
func NotFound(c *gin.Context)                  // 404
func TooManyRequests(c *gin.Context)           // 429
func Maintenance(c *gin.Context)               // 503
func InternalError(c *gin.Context, err error)  // 500, log err ke stderr
```

Contoh penggunaan yang benar:
```go
// ✓ PATCH folder — hanya kembalikan id
func UpdateFolder(c *gin.Context) {
    id := c.Param("id")
    // ... update logic ...
    response.Updated(c, id)
}

// ✓ POST create share — kembalikan server-generated fields
func CreateShareLink(c *gin.Context) {
    // ... create logic ...
    response.Data(c, gin.H{
        "id":        link.ID,
        "token":     link.Token,
        "share_url": buildShareURL(link.Token),
        "expires_at": link.ExpiresAt,
    })
}

// ✓ GET list — kembalikan full data
func ListDocuments(c *gin.Context) {
    // ...
    response.OK(c, gin.H{"data": docs, "total": total, "page": page})
}
```

---

## main.go Route Registration

```
r := gin.New()
r.Use(gin.Logger(), gin.Recovery())
r.Use(middleware.RateLimit())
r.Use(middleware.MaintenanceMode())   ← PERTAMA, sebelum semua route

r.GET("/health", handlers.Health)

api := r.Group("/api/v1")
{
    auth := api.Group("/auth")
    {
        auth.POST("/register", handlers.Register)
        auth.POST("/login",    handlers.Login)
        auth.POST("/refresh",  handlers.RefreshToken)
        auth.POST("/logout",   middleware.Auth(), handlers.Logout)
    }

    // semua di bawah butuh auth
    protected := api.Group("/", middleware.Auth())
    {
        protected.GET("/documents",                  handlers.ListDocuments)
        protected.GET("/documents/:id",              handlers.GetDocument)
        protected.POST("/documents/presign-upload",  handlers.PresignUpload)
        protected.POST("/documents/confirm-upload",  handlers.ConfirmUpload)
        protected.DELETE("/documents/:id",           handlers.DeleteDocument)
        protected.GET("/documents/:id/download",     handlers.DownloadDocument)

        protected.GET("/folders",                    handlers.ListFolders)
        protected.GET("/folders/:id",                handlers.GetFolder)
        protected.POST("/folders",                   handlers.CreateFolder)
        protected.PATCH("/folders/:id",              handlers.UpdateFolder)
        protected.DELETE("/folders/:id",             handlers.DeleteFolder)

        protected.POST("/share",                     handlers.CreateShareLink)
        protected.DELETE("/share/:id",               handlers.DeleteShareLink)

        protected.GET("/activity",                   handlers.ListActivity)
    }

    // share access — public
    api.GET("/share/:token", handlers.AccessShareLink)
}
```

---

## Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o prabodrive-api ./main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/prabodrive-api .
EXPOSE 8080
CMD ["./prabodrive-api"]
```

---

## docker-compose.yml

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:18
    environment:
      POSTGRES_DB: prabodrive
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./database/migrations/001_initial.sql:/docker-entrypoint-initdb.d/001_initial.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

---

## Security Requirements

- Password: bcrypt cost 12
- JWT secret: minimum 32 karakter dari env, panic saat startup jika tidak ada
- S3 bucket: selalu private, akses hanya via presigned URL
- Share link token: `crypto/rand` 32 bytes → hex encode
- Ownership validation: setiap endpoint yang melibatkan document/folder/share harus verifikasi `user_id` dari JWT cocok dengan owner di DB
- CORS: izinkan origin dari `CLOUDFRONT_DOMAIN` env + `localhost:3000` untuk development
- Header `X-Forwarded-Proto`: log jika bukan HTTPS di production (GIN_MODE=release)
- Sanitasi nama file: hapus karakter selain `[a-zA-Z0-9._-]` sebelum dijadikan S3 key

---

## Build Order

**Langkah 0 wajib sebelum apapun:** baca semua file yang ada, laporkan kondisi saat ini, baru lanjut.

Kerjakan dalam urutan ini agar tidak ada dependency yang hilang:

1. Periksa `go.mod` — jika belum ada, `go mod init github.com/[username]/prabodrive-be` lalu tambah semua dependency. Jika sudah ada, cek apakah dependency yang dibutuhkan sudah tercantum.
2. `config/config.go` — load env, init DB pool, init AWS session
3. `database/` — koneksi + migration runner
4. `utils/response.go` — helper response
5. `utils/jwt.go` — generate + validate token
6. `utils/mime.go` — ALLOWED_MIME_TYPES + ValidateMIME
7. `models/` — semua struct
8. `middleware/maintenance.go` → `middleware/auth.go` → `middleware/ratelimit.go`
9. `services/s3.go` → `services/quota.go` → `services/email.go`
10. `handlers/health.go` → `handlers/auth.go` → `handlers/document.go` → `handlers/folder.go` → `handlers/share.go` → `handlers/activity.go`
11. `main.go` — route registration
12. `Dockerfile` + `docker-compose.yml`
13. `.env.example`
14. Jalankan `go build ./...` dan `go vet ./...` — harus zero error/warning
15. Jalankan `docker-compose up` dan verifikasi `GET /health` → 200

---

## Success Criteria

- [ ] `docker-compose up` berhasil tanpa error
- [ ] `GET /health` → `{"status":"ok"}` selalu 200
- [ ] `MAINTENANCE_MODE=true` → semua endpoint selain `/health` return 503
- [ ] Register + login + refresh token flow berjalan
- [ ] Presign upload → confirm upload → list documents berjalan
- [ ] Upload file > 5MB → 400 error
- [ ] Upload MIME yang tidak diizinkan → 400 error
- [ ] Upload ketika quota habis → 403 error
- [ ] Share link dengan password berjalan
- [ ] Rate limiter aktif (test dengan 101 request cepat)
- [ ] Ownership check: user A tidak bisa akses dokumen user B