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
  "nama": "myname",
  "kelas": "3"
}
```

`localhost:81/register/guru`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "id_mapel": 4,
  "nama": "User Name",
  "jabatan": "guru matematika"
}
```

`localhost:81/register/admin`

Body (JSON):

```bash
{
  "email": "user123@gmail.com",
  "nama": "User Name",
  "keterangan": "admin baru"
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

📚 Get Generic Mapels by Kelas

Endpoint:
`GET /genericMapels/:kelas`

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

📘 Get Generic Modules
Endpoint:
`GET /genericModules/:kelas/:id_mapel`

Description:
Get list of modules based on class and subject.

Example :
`localhost:81/genericModules/3/3`
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
GET /getDataUser

Description:
Get the currently authenticated user's data.

Example :
`http://localhost:81/getDataUser`
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
  "kelas": "3"
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
    "kelas": "3",
    "date_created": "0001-01-01T00:00:00Z"
  }
}
```
