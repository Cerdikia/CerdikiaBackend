# Dokumentasi API Rekap Semester

Dokumentasi ini berisi informasi detail tentang endpoint-endpoint yang berkaitan dengan Rekap Semester pada aplikasi Cerdikia. Sistem Rekap Semester digunakan untuk mengarsipkan data progres belajar siswa dari tabel logs ke tabel data_siswa berdasarkan tahun ajaran dan semester.

## Daftar Endpoint

1. [Rekap Semester](#1-rekap-semester)
2. [Edit Tahun Ajaran](#2-edit-tahun-ajaran)
3. [Get All Data Siswa](#3-get-all-data-siswa)
4. [Get Data Siswa](#4-get-data-siswa)
5. [Delete Data Siswa](#5-delete-data-siswa)

---

## 1. Rekap Semester

**Endpoint:** `POST /rekap-semester`

**Deskripsi:** Mengumpulkan data dari tabel logs dan menyimpannya ke tabel data_siswa dengan pengelompokan berdasarkan siswa, tahun ajaran, dan semester.

**Autentikasi:** Diperlukan (Protected Route)

### Request Parameters

Parameter dapat dikirim melalui query parameters atau JSON body.

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| tahun_ajaran | string | Ya | Tahun ajaran dalam format "YYYY/YYYY", contoh: "2024/2025" |
| semester | string | Ya | Semester: "Ganjil" atau "Genap" |
| delete_logs_data | boolean | Tidak | Apakah data logs akan dihapus setelah rekap (default: false) |
| filter_kelas | integer | Tidak | Filter berdasarkan ID kelas tertentu |
| filter_mapel | integer | Tidak | Filter berdasarkan ID mapel tertentu |
| start_date | string | Tidak | Filter berdasarkan tanggal mulai (format: YYYY-MM-DD) |
| end_date | string | Tidak | Filter berdasarkan tanggal akhir (format: YYYY-MM-DD) |

### Request Body (JSON)

```json
{
  "tahun_ajaran": "2024/2025",
  "semester": "Ganjil",
  "delete_logs_data": true,
  "filter_kelas": 1,
  "filter_mapel": 2,
  "start_date": "2024-07-01",
  "end_date": "2024-12-31"
}
```

### Response

```json
{
  "message": "Data siswa berhasil direkap",
  "tahun_ajaran": "2024/2025",
  "semester": "Ganjil",
  "success_count": 25,
  "skipped_count": 2,
  "logs_deleted": true
}
```

### Catatan

- Endpoint ini akan mengelompokkan data logs berdasarkan email siswa
- Data progres disimpan dalam format JSON di kolom progres
- Jika data untuk siswa dengan tahun ajaran dan semester yang sama sudah ada, data tersebut akan dilewati (skipped)
- Jika parameter `delete_logs_data` bernilai true, data logs yang sudah direkap akan dihapus

---

## 2. Edit Tahun Ajaran

**Endpoint:** `POST /edit-tahun-ajaran`

**Deskripsi:** Memperbarui tahun ajaran pada data rekap semester yang sudah ada.

**Autentikasi:** Diperlukan (Protected Route)

### Request Body (JSON)

```json
{
  "tahun_ajaran_lama": "2024/2025",
  "tahun_ajaran_baru": "2025/2026",
  "semester": "Ganjil"
}
```

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| tahun_ajaran_lama | string | Ya | Tahun ajaran lama dalam format "YYYY/YYYY" |
| tahun_ajaran_baru | string | Ya | Tahun ajaran baru dalam format "YYYY/YYYY" |
| semester | string | Ya | Semester: "Ganjil", "Genap", atau "" (kosong untuk semua semester) |

### Response

```json
{
  "message": "Tahun ajaran berhasil diperbarui",
  "rows_affected": 15,
  "tahun_ajaran_lama": "2024/2025",
  "tahun_ajaran_baru": "2025/2026",
  "semester": "Ganjil"
}
```

---

## 3. Get All Data Siswa

**Endpoint:** `GET /rekap-semester-all`

**Deskripsi:** Mendapatkan semua data rekap semester dengan filter opsional.

**Autentikasi:** Diperlukan (Protected Route)

### Query Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| tahun_ajaran | string | Tidak | Filter berdasarkan tahun ajaran |
| semester | string | Tidak | Filter berdasarkan semester |
| id_kelas | string | Tidak | Filter berdasarkan ID kelas |
| email | string | Tidak | Filter berdasarkan email siswa |

### Response

```json
{
  "message": "Data rekap semester berhasil diambil",
  "count": 2,
  "data": [
    {
      "id_data": 1,
      "email": "siswa@example.com",
      "id_kelas": 1,
      "kelas": "X IPA 1",
      "nama_siswa": "Nama Siswa",
      "progres": [
        {
          "id_mapel": 1,
          "id_module": 2,
          "skor": 85,
          "created_at": "2024-09-15 10:30:45"
        },
        {
          "id_mapel": 1,
          "id_module": 3,
          "skor": 90,
          "created_at": "2024-09-20 14:15:30"
        }
      ],
      "tahun_ajaran": "2024/2025",
      "semester": "Ganjil",
      "created_at": "2024-12-20 09:45:12"
    },
    {
      "id_data": 2,
      "email": "siswa2@example.com",
      "id_kelas": 1,
      "kelas": "X IPA 1",
      "nama_siswa": "Nama Siswa Dua",
      "progres": [
        {
          "id_mapel": 2,
          "id_module": 1,
          "skor": 75,
          "created_at": "2024-09-10 11:20:15"
        }
      ],
      "tahun_ajaran": "2024/2025",
      "semester": "Ganjil",
      "created_at": "2024-12-20 09:45:12"
    }
  ]
}
```

---

## 4. Get Data Siswa

**Endpoint:** `GET /rekap-semester/:id_data`

**Deskripsi:** Mendapatkan data rekap semester berdasarkan ID data.

**Autentikasi:** Diperlukan (Protected Route)

### Path Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_data | string | Ya | ID data rekap semester |

### Response

```json
{
  "message": "Data siswa berhasil diambil",
  "data": {
    "id_data": 1,
    "email": "siswa@example.com",
    "id_kelas": 1,
    "progres": [
      {
        "id_mapel": 1,
        "id_module": 2,
        "skor": 85,
        "created_at": "2024-09-15 10:30:45"
      },
      {
        "id_mapel": 1,
        "id_module": 3,
        "skor": 90,
        "created_at": "2024-09-20 14:15:30"
      }
    ],
    "tahun_ajaran": "2024/2025",
    "semester": "Ganjil",
    "created_at": "2024-12-20 09:45:12"
  }
}
```

---

## 5. Delete Data Siswa

**Endpoint:** `DELETE /rekap-semester/:id_data`

**Deskripsi:** Menghapus data rekap semester berdasarkan ID data.

**Autentikasi:** Diperlukan (Protected Route)

### Path Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_data | string | Ya | ID data rekap semester yang akan dihapus |

### Response

```json
{
  "message": "Data siswa berhasil dihapus",
  "id_data": "1"
}
```

---

## Struktur Data

### Data Siswa

Data siswa disimpan dalam tabel `data_siswa` dengan struktur sebagai berikut:

| Field | Tipe | Deskripsi |
|-------|------|------------|
| id_data | uint | Primary key |
| email | string | Email siswa |
| id_kelas | int | ID kelas siswa |
| progres | JSON | Data progres dalam format JSON |
| tahun_ajaran | string | Tahun ajaran (format: "YYYY/YYYY") |
| semester | string | Semester ("Ganjil" atau "Genap") |
| created_at | datetime | Waktu pembuatan data |

### Progres Item

Data progres disimpan dalam format JSON dengan struktur sebagai berikut:

```json
[
  {
    "id_mapel": 1,
    "id_module": 2,
    "skor": 85,
    "created_at": "2024-09-15 10:30:45"
  },
  {
    "id_mapel": 1,
    "id_module": 3,
    "skor": 90,
    "created_at": "2024-09-20 14:15:30"
  }
]
```

## Alur Kerja Rekap Semester

1. Data aktivitas belajar siswa dicatat dalam tabel `logs`
2. Pada akhir semester, admin menggunakan endpoint `POST /rekap-semester` untuk mengarsipkan data
3. Sistem mengelompokkan data logs berdasarkan email siswa
4. Data progres dikonversi ke format JSON dan disimpan dalam tabel `data_siswa`
5. Data logs dapat dihapus setelah proses rekap berhasil (opsional)
6. Data rekap dapat diakses melalui endpoint `GET /rekap-semester-all` atau `GET /rekap-semester/:id_data`

## Catatan Penting

- Format tahun ajaran harus mengikuti pola "YYYY/YYYY" (contoh: "2024/2025")
- Tahun akhir harus lebih besar 1 dari tahun awal (contoh: 2024/2025 valid, 2024/2026 tidak valid)
- Nilai semester hanya boleh "Ganjil" atau "Genap"
- Data rekap semester bersifat permanen dan tidak akan berubah meskipun data logs berubah setelah proses rekap
- Untuk mengubah tahun ajaran pada data yang sudah ada, gunakan endpoint `POST /edit-tahun-ajaran`
