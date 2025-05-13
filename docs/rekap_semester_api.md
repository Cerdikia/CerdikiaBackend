# Rekap Semester API Documentation

This document provides information about the Rekap Semester API endpoints, which are used to manage semester recap data for students.

## Base URL

All endpoints are relative to the base URL of your API server.

## Authentication

All endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Endpoints

### 1. Create Semester Recap

Creates a semester recap for a specific academic year.

- **URL**: `/rekap-semester`
- **Method**: `POST`
- **Auth Required**: Yes

#### Request Body

```json
{
  "tahun_ajaran": "2025/2026",
  "delete_logs_data": false,
  "start_date": "2025-01-01",
  "end_date": "2025-06-30"
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `tahun_ajaran` | string | The academic year in format "YYYY/YYYY" |
| `delete_logs_data` | boolean | Whether to delete existing data for the academic year |
| `start_date` | string | Start date for filtering logs (format: YYYY-MM-DD) |
| `end_date` | string | End date for filtering logs (format: YYYY-MM-DD) |

#### Response

```json
{
  "message": "Rekap semester berhasil dibuat",
  "data": [
    {
      "id_data": 1,
      "email": "student@example.com",
      "nama_siswa": "Student Name",
      "id_kelas": 1,
      "tahun_ajaran": "2025/2026",
      "jumlah_logs": 10
    }
  ]
}
```

### 2. Create Semester Recap for a Specific Student

Creates a semester recap for a specific student and academic year.

- **URL**: `/rekap-semester-siswa`
- **Method**: `POST`
- **Auth Required**: Yes

#### Request Body

```json
{
  "email": "student@example.com",
  "tahun_ajaran": "2025/2026",
  "delete_logs_data": false,
  "start_date": "2025-01-01",
  "end_date": "2025-06-30"
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `email` | string | The email of the student |
| `tahun_ajaran` | string | The academic year in format "YYYY/YYYY" |
| `delete_logs_data` | boolean | Whether to delete existing data for the student and academic year |
| `start_date` | string | Start date for filtering logs (format: YYYY-MM-DD) |
| `end_date` | string | End date for filtering logs (format: YYYY-MM-DD) |

#### Response

```json
{
  "message": "Rekap semester untuk siswa berhasil dibuat",
  "data": {
    "id_data": 1,
    "email": "student@example.com",
    "nama_siswa": "Student Name",
    "id_kelas": 1,
    "tahun_ajaran": "2025/2026",
    "jumlah_logs": 10
  }
}
```

### 3. Create Semester Recap for All Students

Creates semester recaps for all students for a specific academic year.

- **URL**: `/rekap-semester-all-siswa`
- **Method**: `POST`
- **Auth Required**: Yes

#### Request Body

```json
{
  "tahun_ajaran": "2025/2026",
  "delete_logs_data": false,
  "start_date": "2025-01-01",
  "end_date": "2025-06-30"
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `tahun_ajaran` | string | The academic year in format "YYYY/YYYY" |
| `delete_logs_data` | boolean | Whether to delete existing data for the academic year |
| `start_date` | string | Start date for filtering logs (format: YYYY-MM-DD) |
| `end_date` | string | End date for filtering logs (format: YYYY-MM-DD) |

#### Response

```json
{
  "message": "Rekap semester untuk semua siswa berhasil dibuat",
  "total_siswa": 10,
  "total_logs_processed": 150,
  "data": [
    {
      "id_data": 1,
      "email": "student1@example.com",
      "nama_siswa": "Student One",
      "id_kelas": 1,
      "tahun_ajaran": "2025/2026",
      "jumlah_logs": 15
    },
    {
      "id_data": 2,
      "email": "student2@example.com",
      "nama_siswa": "Student Two",
      "id_kelas": 1,
      "tahun_ajaran": "2025/2026",
      "jumlah_logs": 0
    }
  ]
}
```

### 4. Get All Student Data

Retrieves all student data with optional filtering.

- **URL**: `/rekap-semester-all`
- **Method**: `GET`
- **Auth Required**: Yes

#### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `tahun_ajaran` | string | (Optional) Filter by academic year |
| `id_kelas` | string | (Optional) Filter by class ID |
| `email` | string | (Optional) Filter by student email |

#### Response

```json
{
  "message": "Data rekap semester berhasil diambil",
  "count": 2,
  "data": [
    {
      "id_data": 1,
      "email": "student1@example.com",
      "id_kelas": 1,
      "kelas": "Kelas 10A",
      "nama_siswa": "Student One",
      "progres": [
        {
          "module": 1,
          "topic": 1,
          "subtopic": 1,
          "status": "completed",
          "score": 90
        }
      ],
      "tahun_ajaran": "2025/2026"
    },
    {
      "id_data": 2,
      "email": "student2@example.com",
      "id_kelas": 1,
      "kelas": "Kelas 10A",
      "nama_siswa": "Student Two",
      "progres": [],
      "tahun_ajaran": "2025/2026"
    }
  ]
}
```

### 5. Get Student Data by ID

Retrieves a specific student's data by ID.

- **URL**: `/rekap-semester/:id_data`
- **Method**: `GET`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id_data` | integer | The ID of the student data record |

#### Response

```json
{
  "message": "Data rekap semester berhasil diambil",
  "data": {
    "id_data": 1,
    "email": "student@example.com",
    "id_kelas": 1,
    "kelas": "Kelas 10A",
    "nama_siswa": "Student Name",
    "progres": [
      {
        "module": 1,
        "topic": 1,
        "subtopic": 1,
        "status": "completed",
        "score": 90
      }
    ],
    "tahun_ajaran": "2025/2026"
  }
}
```

### 6. Delete Student Data

Deletes a specific student's data record.

- **URL**: `/rekap-semester/:id_data`
- **Method**: `DELETE`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id_data` | integer | The ID of the student data record to delete |

#### Response

```json
{
  "message": "Data rekap semester berhasil dihapus"
}
```

### 7. Edit Academic Year

Updates the academic year for existing records.

- **URL**: `/edit-tahun-ajaran`
- **Method**: `POST`
- **Auth Required**: Yes

#### Request Body

```json
{
  "tahun_ajaran_lama": "2025/225",
  "tahun_ajaran_baru": "2025/2026"
}
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `tahun_ajaran_lama` | string | The current academic year to be updated |
| `tahun_ajaran_baru` | string | The new academic year value |

#### Response

```json
{
  "message": "Tahun ajaran berhasil diubah",
  "records_updated": 10
}
```

### 8. Test Logs Retrieval (Debug Endpoint)

Test endpoint for debugging logs retrieval.

- **URL**: `/test-logs`
- **Method**: `GET`
- **Auth Required**: Yes

#### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `email` | string | (Optional) Filter by student email |
| `start_date` | string | (Optional) Start date for filtering logs (format: YYYY-MM-DD) |
| `end_date` | string | (Optional) End date for filtering logs (format: YYYY-MM-DD) |

#### Response

```json
{
  "message": "Logs retrieved successfully",
  "total_logs": 15,
  "logs": [
    {
      "id": 1,
      "email": "student@example.com",
      "module": 1,
      "topic": 1,
      "subtopic": 1,
      "status": "completed",
      "score": 90,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

## Error Responses

### Authentication Error

```json
{
  "message": "Unauthorized"
}
```

### Validation Error

```json
{
  "message": "Validation error",
  "error": "tahun_ajaran is required"
}
```

### Server Error

```json
{
  "message": "Internal server error",
  "error": "Error message details"
}
```

## Data Structures

### ProgresItem

Represents a student's progress on a specific module, topic, and subtopic.

```json
{
  "module": 1,
  "topic": 1,
  "subtopic": 1,
  "status": "completed",
  "score": 90
}
```

| Field | Type | Description |
|-------|------|-------------|
| `module` | integer | Module number |
| `topic` | integer | Topic number within the module |
| `subtopic` | integer | Subtopic number within the topic |
| `status` | string | Status of completion (e.g., "completed", "in_progress") |
| `score` | integer | Score achieved (0-100) |
