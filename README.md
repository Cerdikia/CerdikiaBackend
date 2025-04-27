# CerdikiaBackend

API ini digunakan untuk manajemen pengguna, autentikasi, serta akses data pembelajaran seperti mapel dan modul, disesuaikan berdasarkan peran pengguna (admin, guru, siswa, dll).

---

## 🚀 Fitur Utama

- 🔐 Register & Login User (dengan Role)
- 🔄 Refresh Token
- 📘 Ambil Mapel & Modul Berdasarkan Kelas
- 👥 Manajemen Data Pengguna & Aktor

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

🔐 Register
Endpoint:

`POST /register/:role`

Description:
Daftar akun baru dengan role tertentu(siswa, guru, atau admin).

Example :

`localhost:81/register/siswa`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "nama": "myname",  // optional
  "kelas": "3"       // optional
}
```

`localhost:81/register/guru`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "id_mapel": 4,
  "nama": "User Name",
  "jabatan": "guru matematika" // optional
}
```

`localhost:81/register/admin`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "nama": "User Name",
  "keterangan": "admin baru" // optional
}
```

Response

```bash
{
  "message": "User dengan nama myname berhasil dibuat"
}
```

🔐 Login
Endpoint:
POST /login

Description:
User login to receive access and refresh tokens.

Body (JSON):

```bash
{
  "username": "user123",
  "role": "siswa"
}
```

Response :

```bash
{
  "Message": "E-mail user123@gmail.com Successfuly Login",
  "Data":
  {
    "email": "user123@gmail.com",
    "nama": "user123",
    "role": "siswa",
    "date_created": "2025-04-24T15:29:29+07:00",
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

🔄 Refresh Token
Endpoint:
POST /refresh

Description:
Refresh the access token using a valid refresh token.

Body (JSON):

```bash
{
  "refresh_token": "..."
}
```

Response :

```bash
{
  "Message": "E-mail user123@gmail.com Successfuly Login",
  "Data":
  {
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

# 🔒 Protected Routes (Require Authentication)

💯 Get User Point

Endpoint:
`GET /point`

Description :
Retrieve user point.

Example :
`localhost:81/point`
Response :

```bash
{
  "Message": "Data retrieved successfully",
  "Data":
  {
    "email": "user123@gmail.com",
    "diamond": 0,
    "exp": 0
  }
}
```

💯 Get User Point

Endpoint:
`PUT /point`

Description :
set user point, you can jus set 1 value instead of both.

Example :
`localhost:81/point`

Request Body :

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
