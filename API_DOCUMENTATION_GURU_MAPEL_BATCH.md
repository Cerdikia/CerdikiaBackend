# Cerdikia Backend - Guru-Mapel Batch Operations API Documentation

This document provides detailed information about the batch operations for Teacher-Subject Relations (Guru-Mapel Relasi) in the Cerdikia Backend API.

## Table of Contents
- [Overview](#overview)
- [Batch Endpoints](#batch-endpoints)
  - [Batch Add Teacher-Subject Relations](#batch-add-teacher-subject-relations)
  - [Batch Update Teacher-Subject Relations](#batch-update-teacher-subject-relations)
  - [Batch Delete Teacher-Subject Relations](#batch-delete-teacher-subject-relations)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Overview

The Guru-Mapel Batch Operations API allows you to efficiently manage relationships between multiple teachers (guru) and their subjects (mapel) in a single request. These endpoints are particularly useful for:

- Initial setup of teacher-subject assignments at the beginning of a school term
- Bulk updates when curriculum changes affect multiple teachers
- Mass cleanup operations when reorganizing teaching assignments

All batch operations use database transactions to ensure data consistency, meaning that either all operations succeed or none do (except for batch delete, which has special handling for non-existent teachers).

## Batch Endpoints

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
    },
    {
      "id_guru": 3,
      "id_mapels": [5]
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
    ],
    "3": [
      {
        "IDGuru": 3,
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
  "error": "Gagal menambahkan relasi guru-mapel",
  "id_guru": 1
}
```

**Notes:**
- This endpoint allows you to add subject relationships for multiple teachers in a single request
- The operation is wrapped in a database transaction to ensure consistency
- If any teacher doesn't exist, the entire operation will fail
- It uses an "upsert" approach that ignores duplicate entries (if a relationship already exists, it won't create a duplicate)
- The `teachers` field must be an array of objects, each containing a teacher ID and an array of subject IDs

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
    },
    {
      "id_guru": 3,
      "id_mapels": []
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
    ],
    "3": []
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
- If the `id_mapels` array is empty for a teacher (like for teacher ID 3 in the example), all relationships for that teacher will be deleted without creating new ones
- This is useful for completely resetting a teacher's subject assignments

---

### Batch Delete Teacher-Subject Relations

Deletes all subject relationships for multiple teachers at once.

**Endpoint:** `DELETE /guru_mapel/batch`

**Authentication:** Required

**Request Body:**
```json
{
  "teacher_ids": [1, 2, 3, 4]
}
```

**Response (Success - 200):**
```json
{
  "message": "Relasi guru-mapel batch berhasil dihapus",
  "results": {
    "1": 3,
    "2": 2,
    "3": 1,
    "4": 0
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
- If a teacher doesn't exist (like teacher ID 4 in the example), it will be skipped (with a count of 0) and the operation will continue with other teachers
- The operation is wrapped in a database transaction to ensure consistency
- This is useful for cleaning up relationships for teachers who are no longer active or for completely resetting multiple teachers' assignments

## Data Models

### BatchGuruMapelRequest
```json
{
  "teachers": [               // Array of teacher-subject relation requests
    {
      "id_guru": "integer",   // Teacher ID
      "id_mapels": ["integer"] // Array of Subject IDs
    }
  ]
}
```

### BatchDeleteRequest
```json
{
  "teacher_ids": ["integer"]  // Array of Teacher IDs to delete relationships for
}
```

## Error Handling

The batch endpoints have specific error handling strategies:

1. **Validation Errors (400)**: Returned when the request body doesn't match the expected format or when required fields are missing.

2. **Not Found Errors (404)**: 
   - For Batch Add and Batch Update: If any teacher in the request doesn't exist, the entire operation fails with a 404 error.
   - For Batch Delete: Non-existent teachers are skipped and the operation continues with other teachers.

3. **Server Errors (500)**: Returned when there's a database error or other server-side issue.

All error responses include a descriptive message and, when relevant, the ID of the teacher that caused the error.

## Best Practices

1. **Limit Batch Size**: While there's no hard limit on how many teachers you can include in a batch operation, it's recommended to keep the batch size reasonable (e.g., under 100 teachers per request) to avoid timeouts and excessive memory usage.

2. **Use Transactions Wisely**: All batch operations use database transactions, which means they lock the affected rows until the operation completes. For very large batches, consider breaking them up into smaller batches to reduce lock contention.

3. **Error Handling**: Always check for errors in the response and handle them appropriately. For Batch Delete, check the results to see which teachers had their relationships deleted successfully.

4. **Idempotent Operations**: Batch Add and Batch Update are designed to be idempotent, meaning you can safely retry them if they fail without causing duplicate data.

5. **Empty ID Arrays**: For Batch Update, providing an empty `id_mapels` array for a teacher will delete all their subject relationships without creating new ones. This is a valid use case for clearing a teacher's assignments.

6. **Verification**: After performing batch operations, especially for critical data, it's a good practice to verify the results by fetching the current state using the `GET /guru/:id_guru` endpoint for a sample of the affected teachers.
