# Dokumentasi API Laporan Nilai (Score Report)

Dokumentasi ini berisi informasi detail tentang endpoint-endpoint yang berkaitan dengan Laporan Nilai (Score Report) pada aplikasi Cerdikia.

## Daftar Endpoint

1. [Get Score Report](#1-get-score-report)
2. [Get Score Report Summary](#2-get-score-report-summary)
3. [Get Student Score Comparison](#3-get-student-score-comparison)
4. [Get Score Progress](#4-get-score-progress)
5. [Get All Students Report](#5-get-all-students-report)

---

## 1. Get Score Report

**Endpoint:** `GET /score-report`

**Deskripsi:** Menampilkan laporan nilai siswa berdasarkan filter tertentu.

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_kelas | string | Ya (jika id_mapel kosong) | ID kelas untuk filter laporan |
| id_mapel | string | Ya (jika id_kelas kosong) | ID mata pelajaran untuk filter laporan |
| sort_by | string | Tidak | Pengurutan data: `latest` (default) atau `highest` |
| aggregate_by | string | Tidak | Metode agregasi: `first` (default) atau `highest` |

> **Catatan**: Minimal satu dari parameter `id_kelas` atau `id_mapel` harus diisi.

### Response

```json
{
  "Message": "Data laporan nilai berhasil diambil",
  "Data": [
    {
      "id_logs": 123,
      "email": "siswa@example.com",
      "id_kelas": 1,
      "id_mapel": 2,
      "id_module": 5,
      "skor": 85,
      "created_at": "2025-05-10T14:30:45Z"
    },
    ...
  ]
}
```

### Contoh Request

```
GET /score-report?id_kelas=1&sort_by=highest&aggregate_by=highest
```

---

## 2. Get Score Report Summary

**Endpoint:** `GET /score-report-summary`

**Deskripsi:** Menampilkan ringkasan statistik dari laporan nilai siswa.

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_kelas | string | Ya (jika id_mapel kosong) | ID kelas untuk filter laporan |
| id_mapel | string | Ya (jika id_kelas kosong) | ID mata pelajaran untuk filter laporan |
| aggregate_by | string | Tidak | Metode agregasi: `highest` (default) atau `first` |

> **Catatan**: Minimal satu dari parameter `id_kelas` atau `id_mapel` harus diisi.

### Response

```json
{
  "Message": "Ringkasan laporan nilai berhasil diambil",
  "Data": {
    "total_siswa": 25,
    "nilai_rata2": 78,
    "nilai_min": 45,
    "nilai_max": 98,
    "distribusi": [
      {
        "kategori": "A (90-100)",
        "jumlah": 5
      },
      {
        "kategori": "B (80-89)",
        "jumlah": 8
      },
      {
        "kategori": "C (70-79)",
        "jumlah": 7
      },
      {
        "kategori": "D (60-69)",
        "jumlah": 3
      },
      {
        "kategori": "E (0-59)",
        "jumlah": 2
      }
    ]
  }
}
```

### Contoh Request

```
GET /score-report-summary?id_kelas=1&id_mapel=2
```

---

## 3. Get Student Score Comparison

**Endpoint:** `GET /student-score-comparison`

**Deskripsi:** Membandingkan nilai seorang siswa dengan rata-rata kelas.

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| email | string | Ya | Email siswa yang akan dibandingkan nilainya |
| id_kelas | string | Ya (jika id_mapel kosong) | ID kelas untuk filter laporan |
| id_mapel | string | Ya (jika id_kelas kosong) | ID mata pelajaran untuk filter laporan |
| aggregate_by | string | Tidak | Metode agregasi: `highest` (default) atau `first` |

> **Catatan**: Minimal satu dari parameter `id_kelas` atau `id_mapel` harus diisi.

### Response

```json
{
  "Message": "Perbandingan nilai siswa berhasil diambil",
  "Data": [
    {
      "id_mapel": 2,
      "student_score": 85,
      "class_average": 78.5,
      "difference": 6.5,
      "status": "Di atas rata-rata"
    },
    {
      "id_mapel": 3,
      "student_score": 70,
      "class_average": 75.2,
      "difference": -5.2,
      "status": "Di bawah rata-rata"
    }
  ]
}
```

### Contoh Request

```
GET /student-score-comparison?email=siswa@example.com&id_kelas=1
```

---

## 4. Get Score Progress

**Endpoint:** `GET /score-progress`

**Deskripsi:** Mendapatkan progres nilai siswa dari waktu ke waktu untuk mata pelajaran tertentu.

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| email | string | Ya | Email siswa |
| id_mapel | string | Ya | ID mata pelajaran |

### Response

```json
{
  "Message": "Data progres nilai berhasil diambil",
  "Data": [
    {
      "id_logs": 101,
      "email": "siswa@example.com",
      "id_mapel": 2,
      "id_module": 5,
      "skor": 75,
      "created_at": "2025-05-01 10:15:30"
    },
    {
      "id_logs": 120,
      "email": "siswa@example.com",
      "id_mapel": 2,
      "id_module": 6,
      "skor": 80,
      "created_at": "2025-05-05 14:20:45"
    },
    {
      "id_logs": 135,
      "email": "siswa@example.com",
      "id_mapel": 2,
      "id_module": 7,
      "skor": 85,
      "created_at": "2025-05-10 09:30:15"
    }
  ]
}
```

### Contoh Request

```
GET /score-progress?email=siswa@example.com&id_mapel=2
```

---

## 5. Get All Students Report

**Endpoint:** `GET /all-students-report`

**Deskripsi:** Mendapatkan laporan semua siswa termasuk yang belum mengerjakan mata pelajaran.

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_kelas | string | Ya (jika id_mapel kosong) | ID kelas untuk filter laporan |
| id_mapel | string | Ya (jika id_kelas kosong) | ID mata pelajaran untuk filter laporan |
| aggregate_by | string | Tidak | Metode agregasi: `highest` (default) atau `first` |

> **Catatan**: Minimal satu dari parameter `id_kelas` atau `id_mapel` harus diisi.

### Response

```json
{
  "Message": "Data laporan semua siswa berhasil diambil",
  "Data": [
    {
      "email": "siswa1@example.com",
      "nama": "Siswa Satu",
      "id_kelas": 1,
      "kelas": "X IPA 1",
      "has_activity": true,
      "scores": [
        {
          "id_logs": 123,
          "id_mapel": 2,
          "id_module": 5,
          "skor": 85,
          "created_at": "2025-05-10 14:30:45"
        }
      ]
    },
    {
      "email": "siswa2@example.com",
      "nama": "Siswa Dua",
      "id_kelas": 1,
      "kelas": "X IPA 1",
      "has_activity": false,
      "scores": []
    }
  ]
}
```

### Contoh Request

```
GET /all-students-report?id_kelas=1&id_mapel=2
```

---

## Catatan Tambahan

### Metode Agregasi

- **first**: Mengambil nilai pertama (berdasarkan waktu) untuk setiap mata pelajaran.
- **highest**: Mengambil nilai tertinggi untuk setiap mata pelajaran.

### Pengurutan

- **latest**: Mengurutkan berdasarkan waktu terbaru.
- **highest**: Mengurutkan berdasarkan nilai tertinggi.

### Kategori Nilai

- A (90-100)
- B (80-89)
- C (70-79)
- D (60-69)
- E (0-59)
