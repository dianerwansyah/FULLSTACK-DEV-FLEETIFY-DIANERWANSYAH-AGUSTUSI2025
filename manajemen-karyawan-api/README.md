# Manajemen Karyawan API

Sistem API untuk manajemen karyawan, departemen, dan absensi berbasis **Golang + Gin** menggunakan **native SQL** dan dokumentasi **Swagger**.

## Fitur Utama
- **CRUD Karyawan**
- **CRUD Departemen**
- **Absensi Masuk (POST)**
- **Absensi Keluar (PUT)**
- **Log Absensi Karyawan** dengan ketepatan waktu berdasarkan aturan per departemen
- **Soft Delete** untuk semua entitas
- **Audit Log** (`created_by`, `updated_by`, `deleted_by`, `created_at`, `updated_at`, `deleted_at`)
- **JWT Authentication**
- **Swagger Documentation**

---

## Teknologi yang Digunakan
- **Golang** (Gin Framework)
- **MySQL** (Database)
- **database/sql** (Native SQL, tanpa ORM)
- **JWT** (Autentikasi)
- **Swagger** (Dokumentasi API)
- **godotenv** (Load konfigurasi dari `.env`)

---

##  Persiapan

### 1. Clone Repository
```bash
git clone https://github.com/dianerrwansyah/manajemen-karyawan-api.git
cd manajemen-karyawan-api
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Setup Environment Variable
Buat file `.env` di root project:
```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=manajemen_karyawan
JWT_SECRET=your_jwt_secret
```

### 4. Setup Database
Import SQL berikut ke MySQL:

```sql
CREATE DATABASE IF NOT EXISTS manajemen_karyawan;
USE manajemen_karyawan;

-- Tabel Departement
CREATE TABLE departement (
    id VARCHAR(50) PRIMARY KEY,
    departement_name VARCHAR(255) NOT NULL,
    max_clock_in_time TIME NOT NULL,
    max_clock_out_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by VARCHAR(50)
);

-- Tabel Employee
CREATE TABLE employee (
    id VARCHAR(50) PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL UNIQUE,
    departement_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by VARCHAR(50),
    FOREIGN KEY (departement_id) REFERENCES departement(id)
);

-- Tabel Attendance
CREATE TABLE attendance (
    id VARCHAR(50) PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL,
    clock_in TIMESTAMP NULL DEFAULT NULL,
    clock_out TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by VARCHAR(50),
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id)
);

-- Tabel Attendance History
CREATE TABLE attendance_history (
    id VARCHAR(50) PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL,
    attendance_id VARCHAR(50) NOT NULL,
    date_attendance TIMESTAMP NOT NULL,
    attendance_type TINYINT(1) NOT NULL COMMENT '1 = IN, 2 = OUT',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by VARCHAR(50),
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
    FOREIGN KEY (attendance_id) REFERENCES attendance(id)
);

-- Data Awal Departement
INSERT INTO departement (id, departement_name, max_clock_in_time, max_clock_out_time, created_by)
VALUES
(UUID(), 'IT', '09:00:00', '17:00:00', 'system'),
(UUID(), 'HRD', '08:30:00', '16:30:00', 'system');

-- Data Awal Employee
INSERT INTO employee (id, employee_id, departement_id, name, address, password, created_by)
VALUES
(UUID(), 'EMP001', (SELECT id FROM departement WHERE departement_name='IT'), 'Dian Erwansyah', 'Jl. Merdeka No. 10', '$2y$12$Sayj3fjn6J6XrPZvUs0zpuprWh6VuqRqOJORIS7uw9SjtFYIWez4G', 'system'),
(UUID(), 'EMP002', (SELECT id FROM departement WHERE departement_name='HRD'), 'Putra Pratama', 'Jl. Mawar No. 5', '$2y$12$Sayj3fjn6J6XrPZvUs0zpuprWh6VuqRqOJORIS7uw9SjtFYIWez4G', 'system');
```

### 5. Jalankan Aplikasi
```bash
go run main.go
```

---

## Dokumentasi API
Swagger UI akan tersedia di endpoint:
```
http://localhost:8080/swagger/index.html
```

---

## Skema Database
Mengacu pada ERD:
- **departement**: Informasi departemen & jam masuk/keluar maksimal
- **employee**: Data karyawan
- **attendance**: Data absensi harian
- **attendance_history**: Riwayat absensi (IN/OUT)

---

## Catatan Pengembangan
- Gunakan `soft delete` (`deleted_at`) untuk menghapus data
- Semua query menggunakan **native SQL** (`database/sql`)
- Audit log (`created_by`, `updated_by`, `deleted_by`) diisi otomatis oleh middleware dari JWT
