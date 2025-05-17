# Dokumentasi API Ranking

Dokumentasi ini berisi informasi detail tentang endpoint ranking pada aplikasi Cerdikia. Endpoint ranking digunakan untuk mendapatkan peringkat siswa berdasarkan poin pengalaman (EXP) yang mereka peroleh.

## Endpoint Ranking

### Get Ranking By Kelas

**Endpoint:** `GET /ranking`

**Deskripsi:** Mendapatkan daftar peringkat siswa berdasarkan poin pengalaman (EXP). Dapat difilter berdasarkan kelas tertentu atau menampilkan semua siswa dari semua kelas.

**Akses:** Public (tidak memerlukan autentikasi)

### Request Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|-------|------------|
| id_kelas | string | Tidak | ID kelas untuk filter peringkat. Jika tidak disediakan, akan menampilkan peringkat dari semua kelas |

### Response

```json
{
  "data": [
    {
      "ranking": 1,
      "nama": "Siswa Satu",
      "kelas": "X IPA 1",
      "exp": 1500
    },
    {
      "ranking": 2,
      "nama": "Siswa Dua",
      "kelas": "X IPA 1",
      "exp": 1350
    },
    {
      "ranking": 3,
      "nama": "Siswa Tiga",
      "kelas": "X IPA 2",
      "exp": 1200
    }
  ]
}
```

### Penjelasan Response

| Field | Tipe | Deskripsi |
|-------|------|------------|
| ranking | integer | Peringkat siswa berdasarkan EXP (urutan menurun) |
| nama | string | Nama siswa |
| kelas | string | Nama kelas siswa |
| exp | integer | Jumlah poin pengalaman (EXP) yang dimiliki siswa |

### Contoh Request

#### Mendapatkan Peringkat Semua Siswa

```
GET /ranking
```

#### Mendapatkan Peringkat Siswa dari Kelas Tertentu

```
GET /ranking?id_kelas=1
```

### Logika Peringkat

Peringkat ditentukan berdasarkan jumlah EXP yang dimiliki siswa, diurutkan dari yang tertinggi ke terendah. Sistem peringkat menggunakan fungsi `RANK()` SQL yang akan memberikan peringkat yang sama untuk nilai EXP yang sama.

### Integrasi dengan Sistem Poin

Endpoint ranking ini terintegrasi dengan sistem poin pengguna (user_points) yang menyimpan jumlah EXP yang diperoleh siswa dari berbagai aktivitas pembelajaran. Peringkat akan diperbarui secara otomatis saat siswa mendapatkan EXP baru.

### Catatan Implementasi

- Endpoint ini menggunakan JOIN antara tabel `siswa`, `kelas`, dan `user_points` untuk mendapatkan data lengkap
- Peringkat dihitung menggunakan fungsi window SQL `RANK() OVER (ORDER BY exp DESC)`
- Jika filter kelas digunakan, hanya siswa dari kelas tersebut yang akan ditampilkan dalam hasil peringkat
