# CerdikiaBackend

API ini digunakan untuk manajemen pengguna, autentikasi, serta akses data pembelajaran seperti mapel dan modul, disesuaikan berdasarkan peran pengguna (admin, guru, siswa, dll).

---

## 🚀 Fitur Utama

- 🔐 Register & Login User (dengan Role)
- 🔄 Refresh Token
- 📘 Ambil Mapel & Modul Berdasarkan Kelas
- 👥 Manajemen Data Pengguna & Aktor
- 📊 Logs dan Reporting
- 🏆 Ranking dan Point System
- 🎁 Gift Management
- 📝 Soal dan Module Management
- 📚 Kelas Management
- ⚡ Energy System
- 💬 Messaging System

---

## ⚙️ Teknologi

- Golang (Gin Framework)
- GORM (ORM untuk Database)
- JWT untuk autentikasi
- Middleware untuk otorisasi

---

## 🔧 Instalasi & Setup

1. Clone repo ini:
   ```bash
   git clone https://github.com/Cerdikia/CerdikiaBackend.git
   cd CerdikiaBackend
   ```
2. Install dependensi:

   ```bash
   go mod tidy
   ```

3. Konfigurasi file `.env`:

   ```bash
   PORT=your_app_port
   DATABASE_URL=mysqldsn
   ```

4. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## 📘 API Dokumentasi

### 🔓 Public Endpoints

#### 🔐 Authentication

##### Register
**Endpoint:** `POST /register/:role`

**Description:** Daftar akun baru dengan role tertentu (siswa, guru, atau admin).

**Example:**
- `POST /register/siswa`
  ```json
  {
    "email": "user123@gmail.com",
    "nama": "myname",
    "kelas": "3"
  }
  ```

- `POST /register/guru`
  ```json
  {
    "email": "user123@gmail.com",
    "id_mapel": 4,
    "nama": "User Name",
    "jabatan": "guru matematika"
  }
  ```

- `POST /register/admin`
  ```json
  {
    "email": "user123@gmail.com",
    "nama": "User Name",
    "keterangan": "admin baru"
  }
  ```

**Response:**
```json
{
  "message": "User dengan nama myname berhasil dibuat"
}
```

##### Login
**Endpoint:** `POST /login`

**Description:** User login to receive access and refresh tokens.

**Body:**
```json
{
  "username": "user123",
  "role": "siswa"
}
```

**Response:**
```json
{
  "Message": "E-mail user123@gmail.com Successfuly Login",
  "Data": {
    "email": "user123@gmail.com",
    "nama": "user123",
    "role": "siswa",
    "date_created": "2025-04-24T15:29:29+07:00",
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

##### Refresh Token
**Endpoint:** `POST /refresh`

**Description:** Refresh the access token using a valid refresh token.

**Body:**
```json
{
  "refresh_token": "..."
}
```

**Response:**
```json
{
  "Message": "E-mail user123@gmail.com Successfuly Login",
  "Data": {
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

#### 🏆 Ranking

##### Get Ranking By Kelas
**Endpoint:** `GET /ranking`

**Description:** Get ranking data filtered by class ID.

**Query Parameters:**
- `id_kelas` - ID of the class to get rankings for

**Response:**
```json
{
  "Message": "Data ranking berhasil diambil",
  "Data": [
    {
      "id": 1,
      "email": "student1@example.com",
      "nama": "Student One",
      "exp": 1500,
      "diamond": 200,
      "rank": 1
    },
    ...
  ]
}
```

#### 🎁 Gifts Management

##### Get All Gifts
**Endpoint:** `GET /gifts`

**Description:** Get all available gifts.

**Response:**
```json
{
  "Message": "Gifts retrieved successfully",
  "Data": [
    {
      "id": 1,
      "name": "Gift Name",
      "description": "Gift Description",
      "price": 100,
      "image_url": "uploads/gifts/gift1.jpg",
      "stock": 10
    },
    ...
  ]
}
```

##### Get Gift By ID
**Endpoint:** `GET /gifts/:id`

**Description:** Get a specific gift by ID.

**Response:**
```json
{
  "Message": "Gift retrieved successfully",
  "Data": {
    "id": 1,
    "name": "Gift Name",
    "description": "Gift Description",
    "price": 100,
    "image_url": "uploads/gifts/gift1.jpg",
    "stock": 10
  }
}
```

##### Create Gift
**Endpoint:** `POST /gifts`

**Description:** Create a new gift.

**Body (Form Data):**
- `name` - Gift name
- `description` - Gift description
- `price` - Gift price
- `stock` - Gift stock
- `image` - Gift image file

##### Update Gift
**Endpoint:** `PUT /gifts/:id`

**Description:** Update an existing gift.

##### Delete Gift
**Endpoint:** `DELETE /gifts/:id`

**Description:** Delete a gift by ID.

#### 📊 Reports

##### Get Score Report
**Endpoint:** `GET /score-report`

**Description:** Get score report data.

##### Get Score Report Summary
**Endpoint:** `GET /score-report-summary`

**Description:** Get summary of score reports.

##### Get Student Score Comparison
**Endpoint:** `GET /student-score-comparison`

**Description:** Compare student scores.

##### Get Score Progress
**Endpoint:** `GET /score-progress`

**Description:** Get score progress over time.

##### Get All Students Report
**Endpoint:** `GET /all-students-report`

**Description:** Get reports for all students.

#### 👥 Verification

##### Get Verified Users
**Endpoint:** `GET /verified`

**Description:** Get users that have been verified.

##### Get All Verified Users
**Endpoint:** `GET /verifiedes`

**Description:** Get all verified users with additional information (student name, class ID, class name).

##### Update User Verification Batch
**Endpoint:** `PATCH /verifiedes`

**Description:** Update verification status for multiple users.

#### 🖼️ Upload

##### Upload Image
**Endpoint:** `POST /upload-image`

**Description:** Upload an image file.

### 🔒 Protected Routes (Require Authentication)

#### 💯 Point System

##### Get User Point
**Endpoint:** `GET /point`

**Description:** Retrieve user point.

**Response:**
```json
{
  "Message": "Data retrieved successfully",
  "Data": {
    "email": "user123@gmail.com",
    "diamond": 0,
    "exp": 0
  }
}
```

##### Update User Point
**Endpoint:** `PUT /point`

**Description:** Set user point, you can just set one value instead of both.

**Body:**
```json
{
  "diamond": 100,
  "exp": 500
}
```

#### 📝 Logs Management

##### Get Logs By Email
**Endpoint:** `GET /logs/:email`

**Description:** Get logs for a specific email (DEV MODE ONLY).

##### Get All Logs
**Endpoint:** `GET /gegeralLogs`

**Description:** Get all logs in the system.

##### Create Log
**Endpoint:** `POST /logs`

**Description:** Create a new log entry.

##### Get Logs By Email With Token
**Endpoint:** `GET /logs`

**Description:** Get logs for the authenticated user's email.

##### Get Logs By Parameters
**Endpoint:** `GET /logsBy`

**Description:** Get logs filtered by various parameters.

**Example:** `GET /logsBy/email/john@example.com/module/2`

##### Get Logs By Period
**Endpoint:** `GET /logs-periode`

**Description:** Get logs for a specific time period relative to the current time.

**Query Parameters:**
- `periode` - Time period ("today", "week", "month", "semester", "year")

**Period Definitions:**
- **Today**: From current time to 24 hours before
- **Week**: From current time to 7 days before
- **Month**: From current time to 30 days before
- **Semester**: From current time to 6 months before
- **Year**: From current time to 1 year before

#### 📚 Mapel (Subject) Management

##### Get All Mapel
**Endpoint:** `GET /genericAllMapels`

**Description:** Get all subjects (mapel).

##### Get Mapel By ID
**Endpoint:** `GET /genericMapel/:id`

**Description:** Get a specific subject by ID.

##### Get Generic Mapels
**Endpoint:** `GET /genericMapels`

**Description:** Get subjects with filtering options.

**Query Parameters:**
- `id_kelas` - Filter by class ID
- `finished` - Filter by completion status

##### Create Mapel
**Endpoint:** `POST /genericMapels`

**Description:** Create a new subject.

##### Update Mapel
**Endpoint:** `PUT /genericMapels/:id`

**Description:** Update an existing subject.

##### Delete Mapel
**Endpoint:** `DELETE /genericMapels/:id`

**Description:** Delete a subject.

#### 📝 Module Management

##### Get Generic Modules
**Endpoint:** `GET /genericModules`

**Description:** Get modules with filtering options.

##### Create Module
**Endpoint:** `POST /genericModules`

**Description:** Create a new module.

##### Update Module
**Endpoint:** `PUT /genericModules/:id`

**Description:** Update an existing module.

##### Delete Module
**Endpoint:** `DELETE /genericModules/:id`

**Description:** Delete a module.

##### Get Module By ID
**Endpoint:** `GET /genericModule/:id`

**Description:** Get a specific module by ID.

##### Toggle Module Ready Status
**Endpoint:** `PUT /togle-module/:id_module`

**Description:** Toggle the ready status of a module.

#### 📊 Statistics

##### Get Stats
**Endpoint:** `GET /stats`

**Description:** Get system statistics (legacy endpoint).

##### Get All Stats
**Endpoint:** `GET /all-stats`

**Description:** Get all system statistics.

##### Get Recent Activities
**Endpoint:** `GET /recent-activities`

**Description:** Get recent activities in the system.

#### 📝 Soal (Question) Management

##### Get Generic Soal
**Endpoint:** `GET /genericSoal/:id_module`

**Description:** Get questions for a specific module.

##### Upload Soal
**Endpoint:** `POST /upload-soal`

**Description:** Upload questions.

##### Get Data Soal
**Endpoint:** `GET /getDataSoal/:id_soal`

**Description:** Get data for a specific question.

##### Update Data Soal
**Endpoint:** `PUT /editDataSoal/:id_soal`

**Description:** Update a specific question.

##### Delete Soal
**Endpoint:** `DELETE /deleteDataSoal/:id_soal`

**Description:** Delete a specific question.

#### 👥 User Management

##### Get All Users
**Endpoint:** `GET /getAllUsers`

**Description:** Get all users in the system.

##### Get User Data by Role and Email
**Endpoint:** `GET /getDataActor/:role/:email`

**Description:** Get user data filtered by role and email.

##### Get Current User Data
**Endpoint:** `GET /getDataUser`

**Description:** Get data for the currently authenticated user.

##### Update User Data
**Endpoint:** `PUT /editDataUser/:role`

**Description:** Update user data for a specific role.

##### Delete User Data
**Endpoint:** `DELETE /deleteDataUser`

**Description:** Delete user data.

##### Update User Profile Image
**Endpoint:** `PATCH /patchImageProfile/:role/:email`

**Description:** Update the profile image for a user.

##### Change User Role
**Endpoint:** `POST /changeUserRole`

**Description:** Change the role of a user.

#### 👥 Guru-Mapel Relationship

##### Get Mapel By Guru
**Endpoint:** `GET /guru/:id_guru`

**Description:** Get subjects associated with a specific teacher.

##### Add Guru-Mapel Relationship
**Endpoint:** `POST /guru_mapel`

**Description:** Create a relationship between a teacher and a subject.

#### 🎁 Gift Exchange

##### Exchange Gift
**Endpoint:** `POST /tukar-barang`

**Description:** Exchange user points for a gift.

#### 📚 Kelas (Class) Management

##### Get All Kelas
**Endpoint:** `GET /kelas`

**Description:** Get all classes.

##### Get Kelas By ID
**Endpoint:** `GET /kelas/:id`

**Description:** Get a specific class by ID.

##### Create Kelas
**Endpoint:** `POST /kelas`

**Description:** Create a new class.

##### Update Kelas
**Endpoint:** `PUT /kelas/:id`

**Description:** Update an existing class.

##### Delete Kelas
**Endpoint:** `DELETE /kelas/:id`

**Description:** Delete a class.

#### ⚡ Energy Management

##### Get User Energy
**Endpoint:** `GET /user-energy/:email`

**Description:** Get energy level for a specific user.

##### Use Energy
**Endpoint:** `POST /user-energy/:email`

**Description:** Use energy for a user.

##### Add Energy
**Endpoint:** `POST /add-energy/:email`

**Description:** Add energy for a user.

#### 📊 Semester Recap

##### Create Semester Recap
**Endpoint:** `POST /rekap-semester`

**Description:** Create a semester recap.

##### Edit Tahun Ajaran
**Endpoint:** `POST /edit-tahun-ajaran`

**Description:** Edit the academic year.

**Body:**
```json
{
  "tahun_ajaran_lama": "2025/225",
  "tahun_ajaran_baru": "2025/2026"
}
```

##### Get All Student Data
**Endpoint:** `GET /rekap-semester-all`

**Description:** Get all student data for semester recap.

##### Get Student Data
**Endpoint:** `GET /rekap-semester/:id_data`

**Description:** Get specific student data for semester recap.

##### Delete Student Data
**Endpoint:** `DELETE /rekap-semester/:id_data`

**Description:** Delete specific student data from semester recap.

#### 💬 Messaging System

##### Create Message
**Endpoint:** `POST /messages`

**Description:** Create a new message.

##### Get All Messages
**Endpoint:** `GET /messages`

**Description:** Get all messages with query filtering.

##### Get Messages By Recipient
**Endpoint:** `GET /messages/recipient/:dest`

**Description:** Get messages for a specific recipient.

##### Get Messages By Sender
**Endpoint:** `GET /messages/sender/:form`

**Description:** Get messages from a specific sender.

##### Get Messages By Subject
**Endpoint:** `GET /messages/subject/:subject`

**Description:** Get messages with a specific subject.

##### Count Unread Messages
**Endpoint:** `GET /messages/unread/count`

**Description:** Count unread messages for the authenticated user.

##### Count All Unread Messages (Admin)
**Endpoint:** `GET /messages/unread/count/all`

**Description:** Count all unread messages (admin only).

##### Get Message By ID and Mark as Read
**Endpoint:** `GET /messages/:id`

**Description:** Get a specific message by ID and mark it as read.

##### Update Message Status
**Endpoint:** `PATCH /messages/:id/status`

**Description:** Update the status of a message.

##### Mark Message as Read
**Endpoint:** `POST /messages/:id/read`

**Description:** Mark a specific message as read.

```bash
{
  "diamond" : 0, // optional
  "exp": 0       // optional
}
```

Response :

```bash
{
  "Message": "User point updated",
  "Data":
  {
    "diamond": 0,
    "exp": 0
  }
}
```

📚 Get Generic Mapels by id_kelas

Endpoint:
`GET /genericMapels/:id_kelas`

Description :
Retrieve list of mapels (subjects) based on class.

Example :
`localhost:81/genericMapels/3`
Response :

```bash
{
  "Message": "Success",
  "Data": [
    {
      "id_mapel": 3,
      "nama_mapel": "Bahasa Indonesia",
      "kelas": "3",
      "jumlah_modul": 4
    },
    {
      "id_mapel": 4,
      "nama_mapel": "Matematika",
      "kelas": "3",
      "jumlah_modul": 3
    }
  ]
}
```

📘 Get Generic Modules For 1 Class

Endpoint:
`GET /genericModulesClass/:id_kelas/:id_mapel`

Description:
Get list of modules based on class and subject.

Example :
`localhost:81/genericModulesClass/3/3`
Response :

```bash
{
  "Message": "Success",
  "Data": [
    {
      "id_module": 10,
      "module": 1,
      "module_judul": "Pemahaman Bacaan",
      "module_deskripsi": "Latihan memahami isi bacaan pendek dan menjawab pertanyaan"
    },
    {
      "id_module": 11,
      "module": 2,
      "module_judul": "Menulis Narasi Sederhana",
      "module_deskripsi": "Penyusunan cerita pendek dengan struktur awal-tengah-akhir"
    }
  ]
}
```

📘 Get Generic Modules For All Class

Endpoint:
`GET /genericModules/:id_mapel`

Description:
Get list of All modules based on subject.

Example :
`http://localhost:81/genericModules/3`
Response :

```bash
{
  "Message": "Success",
  "Data": [
    {
      "kelas": "1",
      "id_module": 4,
      "module": 1,
      "module_judul": "Mengenal Huruf dan Suku Kata",
      "module_deskripsi": "Pengenalan huruf vokal dan konsonan serta pembentukan suku kata dasar"
    },
    {
      "kelas": "1",
      "id_module": 5,
      "module": 2,
      "module_judul": "Membaca Kata Sederhana",
      "module_deskripsi": "Belajar membaca kata-kata pendek dengan suku kata terbuka"
    }
  ]
}
```

📘 Get Generic Specific Modules

Endpoint:
`GET /genericModule/:id_mapel`

Description:
Get spesific module which contains data on questions for 1 module.

Example :
`localhost:81/genericModule/4`
Response :

```bash
{
  "Message": "Success",
  "Data": [
        {
      "id_soal": 22,
      "id_module": 4,
      "soal": "<p><img src=\"http://localhost:81/uploads/67fbd55b-9ca1-419d-b9e1-c6ba7e718808.png\"></p><p>dasdasd</p>",
      "jenis": "pilihan_ganda",
      "opsi_a": "<p>dasasdasd</p>",
      "opsi_b": "<p>asdasdas</p>",
      "opsi_c": "<p>asdasd</p>",
      "opsi_d": "<p>asdasd</p>",
      "jawaban": "b"
    },
    {
      "id_soal": 24,
      "id_module": 4,
      "soal": "<p><img src=\"http://localhost:81/uploads/ad82c46d-0604-40ea-ae5b-b2ac0ccc5dc1.png\"></p>",
      "jenis": "pilihan_ganda",
      "opsi_a": "<p>asdfasdf</p>",
      "opsi_b": "<p>gdfgsdf</p>",
      "opsi_c": "<p>gdfgs</p>",
      "opsi_d": "<p>dfgsdfgsdfg</p>",
      "jawaban": "c"
    },
  ]
}
```

👥 User & Actor Management
🧑‍🤝‍🧑 Get All Users
Endpoint:
GET /getAllUsers

Description:
Retrieve a list of all users.

Response:

```bash
{
  "Message": "Data retrieved successfully",
  "Data": [
    {
      "email": "user123@gmail.com",
      "nama": "myname",
      "role": "siswa",
      "date_created": "2025-04-13T13:59:11+07:00"
    }
  ]
}
```

👤 Get Actor Data by Role
Endpoint:
GET /getDataActor/:role

Description:
Get detailed user data by role.

URL Params:

role: Role of the user (guru, siswa, etc).

Example :
`localhost:81/getDataActor/siswa`
Response

```bash
{
  "Message": "Data retrieved successfully",
  "Data": [
    {
      "email": "user123@gmail.com",
      "nama": "myname",
      "role": "siswa",
      "kelas": "3",
      "date_created": "2025-04-13T13:59:11+07:00"
    }
  ]
}
```

👥 Get Current User Data
Endpoint:
GET /getDataUser?role

Description:
Get the currently authenticated user's data.

Example :
`http://localhost:81/getDataUser?role=siswa`
Response:

```bash
{
  "Message": "Data retrieved successfully",
  "Data":
  {
    "email": "user123@gmail.com",
    "nama": "myname",
    "role": "siswa",
    "kelas": "3",
    "date_created": "2025-04-13T13:59:11+07:00"
  }
}
```

✏️ Update User Data by Role
Endpoint:
PUT /editDataUser/:role

Description:
Update user data based on the role.

URL Params:

role: Role of the user (guru, siswa, etc).

Body (JSON):

Example :
`localhost:81/editDataUser/siswa`
Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "nama": "myname",
  "kelas": 3
}
```

`localhost:81/editDataUser/guru`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "id_mapel": 4,
  "nama": "User Name",
  "jabatan": "guru matematika"
}
```

`localhost:81/editDataUser/admin`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "nama": "User Name",
  "keterangan": "admin baru"
}
```

Response :

```bash
{
  "Message": "Success",
  "Data":
  {
    "email": "user123@gmail.com",
    "nama": "user123",
    "id_kelas": 3,
    "date_created": "0001-01-01T00:00:00Z"
  }
}
```
