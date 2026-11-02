# 📖 Novelku API

Backend RESTful API untuk platform web novel digital, dibangun dengan performa tinggi menggunakan **Go (Golang)**, **Fiber v2**, **GORM**, dan **PostgreSQL**.

---

## 🌟 Fitur Utama

- **Authentication & User Management:** Mendukung sistem multi-role (`reader`, `author`, `admin`) dengan autentikasi berbasis JWT (*Access Token* & *Refresh Token*).
- **Novel & Genre Management:** Pengelolaan novel dengan status penerbitan (`draft`, `published`, `completed`), otomatisasi *slug*, dan relasi *Many-to-Many* dengan Genre.
- **Chapter System:** Konten chapter fleksibel berbasis `JSONB` (`datatypes.JSON`), penomoran urut otomatis (*Composite Unique Index*), dan *Cascade Deletion*.
- **Reader Engagement & Interactions:**
  - **Library / Bookmark:** Menyimpan novel ke daftar bacaan pengguna (*Composite Primary Key*).
  - **Nested Comments:** Sistem komentar bertingkat/bersarang (*parent-child reply*).
  - **Polymorphic Likes:** Sistem like untuk Novel, Chapter, maupun Komentar.
- **Reliability & Production-Ready:**
  - **Connection Pooling:** Manajemen pool PostgreSQL teroptimasi (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`).
  - **Graceful Shutdown:** Menangani sinyal `SIGINT` / `SIGTERM` dengan batas waktu aman agar transaksi berjalan tuntas sebelum server berhenti.
  - **Standardized Response:** Kontrak JSON terstandarisasi (`data`, `error`, `meta` untuk pagination).
  - **Global Middleware:** Pengaman *panic* (`recover`), *request tracing* (`requestid`), pencatat aktivitas (`logger`), dan keamanan domain (`cors`).

---

## 🛠️ Tech Stack

| Komponen | Teknologi | Keterangan |
| :--- | :--- | :--- |
| **Bahasa** | Go (Golang) | v1.23+ |
| **HTTP Framework** | [GoFiber v2](https://gofiber.io/) | Framework web berbasis Fasthttp dengan performa tinggi |
| **ORM** | [GORM](https://gorm.io/) | ORM developer-friendly dengan fitur AutoMigrate & Preload |
| **Database** | PostgreSQL 16 | Driver Postgres dengan pooling bawaan |
| **Containerization**| Docker & Docker Compose | Manajemen database & environment lokal |

---

## 📁 Struktur Direktori (Layered Architecture)

Proyek ini menerapkan struktur arsitektur berlapis (*Clean / Layered Architecture*) yang modular:

```text
├── cmd/
│   └── api/
│       └── main.go         # Entrypoint: Wiring dependency & startup server
├── internal/               # Kode privat aplikasi
│   ├── config/             # Pemuat konfigurasi .env
│   ├── database/           # Koneksi GORM, Pool, & AutoMigrate
│   ├── domain/             # Struct entitas/model tabel database
│   ├── dto/                # Data Transfer Object (Request & Validasi)
│   ├── repository/         # Query database murni (GORM)
│   ├── service/            # Logika bisnis aplikasi
│   ├── handler/            # HTTP Controller / Handler Fiber
│   └── middleware/         # Middleware Auth JWT, RBAC, dll
├── pkg/                    # Package publik pembantu
│   └── response/           # Standarisasi format JSON response
├── go.mod                  # Manajemen dependency Go
├── go.sum
└── README.md
```

---

## 🚀 Memulai (Quick Start)

### 1. Prasyarat
Pastikan Anda telah menginstal:
- [Go](https://go.dev/dl/) (versi 1.23 atau lebih baru)
- [Docker](https://www.docker.com/) & Docker Compose
- [Git](https://git-scm.com/)

### 2. Kloning Repositori & Setup Environment
```bash
git clone https://github.com/bagasriatmaputera/novelku.git
cd novelku

# Buat file .env di folder project
cp .env.example .env
```

Contoh konfigurasi di file `.env`:
```env
APP_ENV=development
APP_PORT=8080
CORS_ORIGINS=http://localhost:3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=novelku
DB_PASSWORD=novelku_dev_password
DB_NAME=novelku
```

### 3. Menjalankan Database PostgreSQL
Nyalakan container database menggunakan Docker Compose:
```bash
docker compose up -d
```

### 4. Menjalankan Server API
Jalankan server aplikasi:
```bash
# Download dependensi
go mod tidy

# Jalankan server API
go run cmd/api/main.go
```

Server akan aktif di `http://localhost:8080` dan secara otomatis menjalankan migrasi tabel database (**AutoMigrate**).

---

## 🩺 Cek Status Server (Health Check)

Uji apakah server dan koneksi database berjalan baik:

```bash
curl http://localhost:8080/health
```

**Response Sukses (HTTP 200):**
```json
{
  "data": {
    "status": "OK"
  },
  "error": null,
  "meta": null
}
```

---

## 📐 Format Standar API Response

Semua endpoint mengembalikan kontrak response yang konsisten:

### 1. Single / Detail Data (HTTP 200 / 201)
```json
{
  "data": { ... },
  "error": null,
  "meta": null
}
```

### 2. List Data dengan Pagination (HTTP 200)
```json
{
  "data": [ ... ],
  "error": null,
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

### 3. Error Response (HTTP 4xx / 5xx)
```json
{
  "data": null,
  "error": {
    "code": "DB_DOWN",
    "message": "database tidak dapat dihubungi"
  },
  "meta": null
}
```

---

## 🗺️ Roadmap Pengembangan

- [x] Inisialisasi Project & Struktur Folder
- [x] Setup Docker Compose PostgreSQL
- [x] Config Loader & Helper Environment
- [x] Format Response Terstandarisasi & Global Error Handler
- [x] Connection Pooling Database GORM
- [x] Pemodelan Domain (User, Novel, Chapter, Genre, Interaction) & AutoMigrate
- [x] Server Entrypoint dengan Fiber & Graceful Shutdown
- [ ] Implementasi DTO & Validator Request
- [ ] Modul Autentikasi (Register, Login, JWT, Refresh Token)
- [ ] Modul Master Genre (CRUD + Role Admin)
- [ ] Modul Novel (CRUD, Upload Cover, Filter & Search, Pagination)
- [ ] Modul Chapter (Nested Route, Rich Text JSONB)
- [ ] Modul Interaksi (Library, Like, & Nested Comment)
