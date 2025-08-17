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


### Auth
- **POST** `/api/auth/login`  
  ![Login](https://github.com/user-attachments/assets/54ba0d29-03af-488e-97c9-c22dbb99a3d9)

- **POST** `/api/auth/logout`

- **GET** `/api/auth/me`

---

### Attendance
- **POST** `/api/attendance`

- **POST** `/api/attendance/GetData`  
  ![Attendance GetData](https://github.com/user-attachments/assets/4008a420-041a-4848-82a9-9057d6bfa3b7)

- **POST** `/api/attendance/logs`  
  ![Attendance Logs](https://github.com/user-attachments/assets/df748ef3-8f03-4510-9e60-334a043a331c)

- **GET** `/api/attendance/today`

---

### Departement
- **POST** `/api/departement`  
  ![Departement Create](https://github.com/user-attachments/assets/82bd9a24-cdba-46ba-960f-77139dded20e)

- **POST** `/api/departement/GetData`  
  ![Departement GetData](https://github.com/user-attachments/assets/2beed3bf-9de9-4113-86a7-34fa49d08e1c)

- **GET** `/api/departement/{id}`  
  ![Departement Detail](https://github.com/user-attachments/assets/533ceec8-4295-4570-91a9-ab4b7c43ee33)

- **PUT** `/api/departement/{id}`  
  ![Departement Update](https://github.com/user-attachments/assets/e8458ab9-14ea-444f-848d-3fee31e85cd7)

- **DELETE** `/api/departement/{id}`  
  ![Departement Delete](https://github.com/user-attachments/assets/982bf7bd-5941-4d99-a5f1-26d789c83237)

---

### Employee
- **POST** `/api/employee`  
  ![Employee Create](https://github.com/user-attachments/assets/f8a5dcad-545e-4d97-90f1-2ea23cf37930)

- **GET** `/api/employee/{id}`  
  ![Employee Detail](https://github.com/user-attachments/assets/0e77f629-21e7-4a7b-9157-96e37d7f2224)

- **PUT** `/api/employee/{id}`  
  ![Employee Update](https://github.com/user-attachments/assets/daf87c92-efce-42b7-b882-268970bcf279)


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
