# Cerdikia Backend - Guru-Mapel Relasi API Documentation

This document provides detailed information about the Teacher-Subject Relation (Guru-Mapel Relasi) endpoints in the Cerdikia Backend API.

## Table of Contents
- [Overview](#overview)
- [Endpoints](#endpoints)
  - [Get Subjects by Teacher](#get-subjects-by-teacher)
  - [Add Teacher-Subject Relations](#add-teacher-subject-relations)
  - [Update Teacher-Subject Relations](#update-teacher-subject-relations)
  - [Delete Specific Teacher-Subject Relation](#delete-specific-teacher-subject-relation)
  - [Delete All Teacher-Subject Relations](#delete-all-teacher-subject-relations)
  - [Batch Add Teacher-Subject Relations](#batch-add-teacher-subject-relations)
  - [Batch Update Teacher-Subject Relations](#batch-update-teacher-subject-relations)
  - [Batch Delete Teacher-Subject Relations](#batch-delete-teacher-subject-relations)
- [Data Models](#data-models)

## Overview

The Guru-Mapel Relasi API allows you to manage the relationships between teachers (guru) and subjects (mapel) in the Cerdikia system. This includes retrieving the subjects assigned to a specific teacher and assigning multiple subjects to a teacher.

## Endpoints

### Get Subjects by Teacher

Retrieves a teacher's information along with all subjects assigned to that teacher.

**Endpoint:** `GET /guru/:id_guru`

**Authentication:** Required

**Path Parameters:**
- `id_guru` (required): The ID of the teacher

**Response (Success - 200):**
```json
{
  "email": "teacher@example.com",
  "nama": "Teacher Name",
  "jabatan": "Math Teacher",
  "mapel": [
    {
      "id_mapel": 1,
      "mapel": "Mathematics"
    },
    {
      "id_mapel": 2,
      "mapel": "Physics"
    }
  ]
}
```

**Response (Error - 404):**
```json
{
  "error": "Guru tidak ditemukan"
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal mengambil mata pelajaran"
}
```

**Notes:**
- This endpoint first retrieves the teacher's basic information (email, name, position)
- Then it retrieves all subjects assigned to the teacher through the guru_mapel relation table
- The response combines both the teacher information and their assigned subjects

---

### Add Teacher-Subject Relations

Assigns multiple subjects to a teacher by creating relationships in the guru_mapel table.

**Endpoint:** `POST /guru_mapel`

**Authentication:** Required

**Request Body:**
```json
{
  "id_guru": 1,
  "id_mapels": [1, 2, 3]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel berhasil ditambahkan",
  "data": [
    {
      "IDGuru": 1,
      "IDMapel": 1
    },
    {
      "IDGuru": 1,
      "IDMapel": 2
    },
    {
      "IDGuru": 1,
      "IDMapel": 3
    }
  ]
}
```

**Response (Error - 400):**
```json
{
  "error": "Key: 'GuruMapelRequest.IDGuru' Error:Field validation for 'IDGuru' failed on the 'required' tag"
}
```

**Response (Error - 500):**
```json
{
  "error": "Error message details"
}
```

**Notes:**
- This endpoint creates multiple teacher-subject relationships in a single request
- It uses an "upsert" approach that ignores duplicate entries (if a relationship already exists, it won't create a duplicate)
- Both `id_guru` and `id_mapels` are required fields in the request
- The `id_mapels` field must be an array of subject IDs, even if you're only assigning one subject

---

### Update Teacher-Subject Relations

Updates all subject relationships for a teacher by replacing existing relationships with new ones.

**Endpoint:** `PUT /guru_mapel`

**Authentication:** Required

**Request Body:**
```json
{
  "id_guru": 1,
  "id_mapels": [2, 3, 4]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel berhasil diperbarui",
  "data": [
    {
      "IDGuru": 1,
      "IDMapel": 2
    },
    {
      "IDGuru": 1,
      "IDMapel": 3
    },
    {
      "IDGuru": 1,
      "IDMapel": 4
    }
  ]
}
```

**Response (Error - 400):**
```json
{
  "error": "Key: 'GuruMapelRequest.IDGuru' Error:Field validation for 'IDGuru' failed on the 'required' tag"
}
```

**Response (Error - 404):**
```json
{
  "error": "Guru tidak ditemukan"
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menghapus relasi lama"
}
```

**Notes:**
- This endpoint performs a complete replacement of all subject relationships for a teacher
- It first deletes all existing relationships and then creates new ones based on the provided list
- The operation is wrapped in a database transaction to ensure consistency
- If the teacher doesn't exist, the operation will fail
- If the `id_mapels` array is empty, all relationships will be deleted without creating new ones

---

### Delete Specific Teacher-Subject Relation

Deletes a specific relationship between a teacher and a subject.

**Endpoint:** `DELETE /guru_mapel/:id_guru/:id_mapel`

**Authentication:** Required

**Path Parameters:**
- `id_guru` (required): The ID of the teacher
- `id_mapel` (required): The ID of the subject

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel berhasil dihapus"
}
```

**Response (Error - 400):**
```json
{
  "error": "ID guru tidak valid"
}
```

**Response (Error - 404):**
```json
{
  "error": "Relasi guru-mapel tidak ditemukan"
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menghapus relasi guru-mapel"
}
```

**Notes:**
- This endpoint removes a single relationship between a specific teacher and a specific subject
- If the relationship doesn't exist, a 404 error will be returned

---

### Delete All Teacher-Subject Relations

Deletes all subject relationships for a specific teacher.

**Endpoint:** `DELETE /guru_mapel/:id_guru`

**Authentication:** Required

**Path Parameters:**
- `id_guru` (required): The ID of the teacher

**Response (Success - 200):**
```json
{
  "message": "Semua relasi guru-mapel berhasil dihapus",
  "count": 5
}
```

**Response (Error - 400):**
```json
{
  "error": "ID guru tidak valid"
}
```

**Response (Error - 404):**
```json
{
  "error": "Guru tidak ditemukan"
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menghapus relasi guru-mapel"
}
```

**Notes:**
- This endpoint removes all subject relationships for a specific teacher
- The response includes a count of how many relationships were deleted
- If the teacher doesn't exist, a 404 error will be returned

---

### Batch Add Teacher-Subject Relations

Adds subject relationships for multiple teachers at once.

**Endpoint:** `POST /guru_mapel/batch`

**Authentication:** Required

**Request Body:**
```json
{
  "teachers": [
    {
      "id_guru": 1,
      "id_mapels": [1, 2, 3]
    },
    {
      "id_guru": 2,
      "id_mapels": [2, 4]
    }
  ]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel batch berhasil ditambahkan",
  "data": {
    "1": [
      {
        "IDGuru": 1,
        "IDMapel": 1
      },
      {
        "IDGuru": 1,
        "IDMapel": 2
      },
      {
        "IDGuru": 1,
        "IDMapel": 3
      }
    ],
    "2": [
      {
        "IDGuru": 2,
        "IDMapel": 2
      },
      {
        "IDGuru": 2,
        "IDMapel": 4
      }
    ]
  }
}
```

**Response (Error - 400):**
```json
{
  "error": "Key: 'BatchGuruMapelRequest.Teachers' Error:Field validation for 'Teachers' failed on the 'required' tag"
}
```

**Response (Error - 404):**
```json
{
  "error": "Guru tidak ditemukan",
  "id_guru": 1
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menambahkan relasi guru-mapel",
  "id_guru": 1
}
```

**Notes:**
- This endpoint allows you to add subject relationships for multiple teachers in a single request
- The operation is wrapped in a database transaction to ensure consistency
- If any teacher doesn't exist, the entire operation will fail
- It uses an "upsert" approach that ignores duplicate entries

---

### Batch Update Teacher-Subject Relations

Updates subject relationships for multiple teachers at once by replacing existing relationships with new ones.

**Endpoint:** `PUT /guru_mapel/batch`

**Authentication:** Required

**Request Body:**
```json
{
  "teachers": [
    {
      "id_guru": 1,
      "id_mapels": [2, 3, 4]
    },
    {
      "id_guru": 2,
      "id_mapels": [1, 5]
    }
  ]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel batch berhasil diperbarui",
  "data": {
    "1": [
      {
        "IDGuru": 1,
        "IDMapel": 2
      },
      {
        "IDGuru": 1,
        "IDMapel": 3
      },
      {
        "IDGuru": 1,
        "IDMapel": 4
      }
    ],
    "2": [
      {
        "IDGuru": 2,
        "IDMapel": 1
      },
      {
        "IDGuru": 2,
        "IDMapel": 5
      }
    ]
  }
}
```

**Response (Error - 400):**
```json
{
  "error": "Key: 'BatchGuruMapelRequest.Teachers' Error:Field validation for 'Teachers' failed on the 'required' tag"
}
```

**Response (Error - 404):**
```json
{
  "error": "Guru tidak ditemukan",
  "id_guru": 1
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menghapus relasi lama",
  "id_guru": 1
}
```

**Notes:**
- This endpoint performs a complete replacement of all subject relationships for multiple teachers at once
- For each teacher, it first deletes all existing relationships and then creates new ones based on the provided list
- The operation is wrapped in a database transaction to ensure consistency
- If any teacher doesn't exist, the entire operation will fail
- If the `id_mapels` array is empty for a teacher, all relationships for that teacher will be deleted without creating new ones

---

### Batch Delete Teacher-Subject Relations

Deletes all subject relationships for multiple teachers at once.

**Endpoint:** `DELETE /guru_mapel/batch`

**Authentication:** Required

**Request Body:**
```json
{
  "teacher_ids": [1, 2, 3]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel batch berhasil dihapus",
  "results": {
    "1": 3,
    "2": 2,
    "3": 0
  }
}
```

**Response (Error - 400):**
```json
{
  "error": "Key: 'TeacherIDs' Error:Field validation for 'TeacherIDs' failed on the 'required' tag"
}
```

**Response (Error - 500):**
```json
{
  "error": "Gagal menghapus relasi guru-mapel",
  "id_guru": 1
}
```

**Notes:**
- This endpoint removes all subject relationships for multiple teachers in a single request
- The response includes a count of how many relationships were deleted for each teacher
- If a teacher doesn't exist, it will be skipped (with a count of 0) and the operation will continue with other teachers
- The operation is wrapped in a database transaction to ensure consistency

## Data Models

### GuruMapel (Teacher-Subject Relation)
```json
{
  "IDGuru": "integer",  // Teacher ID (primary key)
  "IDMapel": "integer"  // Subject ID (primary key)
}
```

### GuruMapelRequest
```json
{
  "id_guru": "integer",       // Teacher ID
  "id_mapels": ["integer"]    // Array of Subject IDs
}
```

### BatchGuruMapelRequest
```json
{
  "teachers": [               // Array of GuruMapelRequest objects
    {
      "id_guru": "integer",
      "id_mapels": ["integer"]
    }
  ]
}
```

### GuruMapelResponse
```json
{
  "email": "string",
  "nama": "string",
  "jabatan": "string",
  "mapel": [
    {
      "id_mapel": "integer",
      "mapel": "string"
    }
  ]
}
```

### Mapel (Subject)
```json
{
  "id_mapel": "integer",
  "mapel": "string"
}
```

## Database Schema

The Guru-Mapel relationship is implemented using a many-to-many relationship with the following tables:

1. **guru** - Contains teacher information
   - id (primary key)
   - email
   - nama
   - jabatan
   - etc.

2. **mapel** - Contains subject information
   - id_mapel (primary key)
   - mapel (subject name)
   - etc.

3. **guru_mapel** - Junction table for the many-to-many relationship
   - id_guru (foreign key to guru.id, part of composite primary key)
   - id_mapel (foreign key to mapel.id_mapel, part of composite primary key)

This schema allows a teacher to be assigned to multiple subjects, and a subject to be taught by multiple teachers.
