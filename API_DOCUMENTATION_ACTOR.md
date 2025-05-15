# Cerdikia Backend - Actor API Documentation

This document provides detailed information about the Actor-related endpoints in the Cerdikia Backend API.

## Table of Contents
- [User Management](#user-management)
  - [Register User](#register-user)
  - [Get All Users](#get-all-users)
  - [Get User by Email](#get-user-by-email)
  - [Get Data Actor by Role](#get-data-actor-by-role)
  - [Update User Data](#update-user-data)
  - [Delete User](#delete-user)
  - [Change User Role](#change-user-role)
  - [Update Profile Image](#update-profile-image)
- [User Verification](#user-verification)
  - [Get Verification Status](#get-verification-status)
  - [Get All Verification Statuses](#get-all-verification-statuses)
  - [Update Verification Status](#update-verification-status)
- [User Points](#user-points)
  - [Get User Points](#get-user-points)
  - [Update User Points](#update-user-points)
- [User Energy](#user-energy)
  - [Get User Energy](#get-user-energy)
  - [Create User Energy](#create-user-energy)

---

## User Management

### Register User

Registers a new user with a specific role.

**Endpoint:** `POST /register/:role`

**Supported Roles:**
- `POST /register/siswa` - Register a student
- `POST /register/guru` - Register a teacher
- `POST /register/kepalaSekolah` - Register a school principal
- `POST /register/admin` - Register an administrator

**Request Body (for siswa):**
```json
{
  "email": "student@example.com",
  "nama": "Student Name",
  "id_kelas": 3
}
```

**Request Body (for guru/kepalaSekolah):**
```json
{
  "email": "teacher@example.com",
  "nama": "Teacher Name",
  "jabatan": "Math Teacher" // Optional for guru, will be set to "kepala sekolah" for kepalaSekolah
}
```

**Request Body (for admin):**
```json
{
  "email": "admin@example.com",
  "nama": "Admin Name",
  "keterangan": "System Administrator" // Optional
}
```

**Response (Success - 201):**
```json
{
  "message": "User dengan email user@example.com berhasil dibuat"
}
```

**Response (Error - 400/500):**
```json
{
  "error": "Error message details"
}
```

**Notes:**
- When registering a student (`siswa`), the system automatically:
  - Sets a default profile image if not provided
  - Creates initial points record
  - Creates an account verification record with "waiting" status
  - Initializes user energy with 5 points

---

### Get All Users

Retrieves a list of all users in the system.

**Endpoint:** `GET /getAllUsers`

**Authentication:** Required

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": [
    {
      "email": "user1@example.com",
      "nama": "User One",
      "role": "siswa",
      "date_created": "2023-01-01T00:00:00Z",
      "image_profile": "http://example.com/uploads/profile1.png"
    },
    {
      "email": "user2@example.com",
      "nama": "User Two",
      "role": "guru",
      "jabatan": "Math Teacher",
      "date_created": "2023-01-02T00:00:00Z",
      "image_profile": "http://example.com/uploads/profile2.png"
    }
  ]
}
```

---

### Get User by Email

Retrieves user information by email.

**Endpoint:** `GET /getDataUser?email=user@example.com`

**Authentication:** Required

**Query Parameters:**
- `email` (required): The email of the user to retrieve

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "user@example.com",
    "nama": "User Name",
    "role": "siswa",
    "id_kelas": 3,
    "date_created": "2023-01-01T00:00:00Z",
    "image_profile": "http://example.com/uploads/profile.png"
  }
}
```

**Response (Error - 404):**
```json
{
  "message": "User not found",
  "data": null
}
```

---

### Get Data Actor by Role

Retrieves users by their role, optionally filtered by email.

**Endpoint:** `GET /getDataActor/:role/:email?`

**Authentication:** Required

**Path Parameters:**
- `role` (required): The role to filter by (siswa, guru, kepalaSekolah, admin)
- `email` (optional): If provided, returns a specific user with this email

**Response (Success - 200, all users of a role):**
```json
{
  "message": "Success",
  "data": [
    {
      "email": "student1@example.com",
      "nama": "Student One",
      "id_kelas": 3,
      "date_created": "2023-01-01T00:00:00Z",
      "image_profile": "http://example.com/uploads/profile1.png"
    },
    {
      "email": "student2@example.com",
      "nama": "Student Two",
      "id_kelas": 4,
      "date_created": "2023-01-02T00:00:00Z",
      "image_profile": "http://example.com/uploads/profile2.png"
    }
  ]
}
```

**Response (Success - 200, specific user):**
```json
{
  "message": "Success",
  "data": {
    "email": "student1@example.com",
    "nama": "Student One",
    "id_kelas": 3,
    "date_created": "2023-01-01T00:00:00Z",
    "image_profile": "http://example.com/uploads/profile1.png"
  }
}
```

**Response (Error - 400):**
```json
{
  "message": "no parameter found",
  "data": null
}
```

---

### Update User Data

Updates user information based on their role.

**Endpoint:** `PUT /editDataUser/:role`

**Authentication:** Required

**Path Parameters:**
- `role` (required): The role of the user to update (siswa, guru, admin)

**Request Body (for siswa):**
```json
{
  "email": "student@example.com",
  "nama": "Updated Student Name",
  "id_kelas": 4
}
```

**Request Body (for guru):**
```json
{
  "email": "teacher@example.com",
  "nama": "Updated Teacher Name",
  "jabatan": "Science Teacher"
}
```

**Request Body (for admin):**
```json
{
  "email": "admin@example.com",
  "nama": "Updated Admin Name",
  "keterangan": "Senior Administrator"
}
```

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "student@example.com",
    "nama": "Updated Student Name",
    "id_kelas": 4,
    "date_created": "2023-01-01T00:00:00Z",
    "image_profile": "http://example.com/uploads/profile.png"
  }
}
```

**Response (Error - 400):**
```json
{
  "message": "Error message details",
  "data": null
}
```

---

### Delete User

Deletes a user from the system.

**Endpoint:** `DELETE /deleteDataUser`

**Authentication:** Required

**Request Body:**
```json
{
  "email": "user@example.com",
  "role": "siswa"
}
```

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": null
}
```

**Response (Error - 400/404):**
```json
{
  "message": "Error message details",
  "data": null
}
```

---

### Change User Role

Changes a user's role by moving their data between tables.

**Endpoint:** `POST /changeUserRole`

**Authentication:** Required

**Request Body:**
```json
{
  "email": "user@example.com",
  "old_role": "guru",
  "new_role": "siswa"
}
```

**Response (Success - 200):**
```json
{
  "message": "Role changed successfully",
  "data": {
    "email": "user@example.com",
    "nama": "User Name",
    "role": "siswa",
    "date_created": "2023-01-01T00:00:00Z",
    "image_profile": "http://example.com/uploads/profile.png"
  }
}
```

**Response (Error - 400):**
```json
{
  "message": "Error message details",
  "data": null
}
```

**Notes:**
- When changing a role to `siswa`, the system automatically:
  - Creates initial points record
  - Creates an account verification record with "waiting" status
  - Initializes user energy with 5 points
- When changing between `guru` and `kepalaSekolah`, only the `jabatan` field is updated
- Valid roles are: `siswa`, `guru`, `kepalaSekolah`, and `admin`
- The new role must be different from the old role

---

### Update Profile Image

Updates a user's profile image.

**Endpoint:** `PATCH /patchImageProfile/:role/:email`

**Authentication:** Required

**Path Parameters:**
- `role` (required): The role of the user (siswa, guru, admin)
- `email` (required): The email of the user

**Request Body:**
```json
{
  "image_profile": "http://example.com/uploads/new_profile.png"
}
```

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "user@example.com",
    "nama": "User Name",
    "image_profile": "http://example.com/uploads/new_profile.png"
  }
}
```

**Response (Error - 400/404):**
```json
{
  "message": "Error message details",
  "data": null
}
```

---

## User Verification

### Get Verification Status

Retrieves the verification status of a specific user.

**Endpoint:** `GET /verified?email=user@example.com`

**Authentication:** Required

**Query Parameters:**
- `email` (required): The email of the user to check

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "user@example.com",
    "verified_status": "approved",
    "date_created": "2023-01-01T00:00:00Z"
  }
}
```

---

### Get All Verification Statuses

Retrieves verification statuses for all users.

**Endpoint:** `GET /verifiedes`

**Authentication:** Required

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": [
    {
      "email": "user1@example.com",
      "nama": "User One",
      "id_kelas": 3,
      "kelas_name": "Class 3A",
      "verified_status": "waiting",
      "date_created": "2023-01-01T00:00:00Z"
    },
    {
      "email": "user2@example.com",
      "nama": "User Two",
      "id_kelas": 4,
      "kelas_name": "Class 4B",
      "verified_status": "approved",
      "date_created": "2023-01-02T00:00:00Z"
    }
  ]
}
```

---

### Update Verification Status

Updates the verification status for multiple users.

**Endpoint:** `PATCH /verifiedes`

**Authentication:** Required

**Request Body:**
```json
{
  "verifications": [
    {
      "email": "user1@example.com",
      "verified_status": "approved"
    },
    {
      "email": "user2@example.com",
      "verified_status": "rejected"
    }
  ]
}
```

**Response (Success - 200):**
```json
{
  "message": "Verification status updated successfully",
  "data": {
    "success_count": 2,
    "failed_count": 0,
    "failed_emails": []
  }
}
```

**Response (Partial Success - 200):**
```json
{
  "message": "Some verification updates failed",
  "data": {
    "success_count": 1,
    "failed_count": 1,
    "failed_emails": ["user2@example.com"]
  }
}
```

---

## User Points

### Get User Points

Retrieves the points for a specific user.

**Endpoint:** `GET /getPoint?email=user@example.com`

**Authentication:** Required

**Query Parameters:**
- `email` (required): The email of the user

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "user@example.com",
    "point": 100
  }
}
```

**Response (Error - 404):**
```json
{
  "message": "User points not found",
  "data": null
}
```

---

### Update User Points

Updates the points for a specific user.

**Endpoint:** `PUT /updatePoint`

**Authentication:** Required

**Request Body:**
```json
{
  "email": "user@example.com",
  "point": 150
}
```

**Response (Success - 200):**
```json
{
  "message": "Success",
  "data": {
    "email": "user@example.com",
    "point": 150
  }
}
```

**Response (Error - 400/404):**
```json
{
  "message": "Error message details",
  "data": null
}
```

---

## User Energy

### Get User Energy

Retrieves the energy for a specific user.

**Endpoint:** `GET /energy?email=user@example.com`

**Authentication:** Required

**Query Parameters:**
- `email` (required): The email of the user

**Response (Success - 200):**
```json
{
  "email": "user@example.com",
  "energy": 5,
  "last_updated": "2023-01-01T00:00:00Z"
}
```

**Response (Error - 404):**
```json
{
  "error": "User energy not found"
}
```

---

### Create User Energy

Creates initial energy for a user.

**Endpoint:** `POST /energy`

**Authentication:** Required

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response (Success - 200):**
```json
{
  "message": "User energy initialized successfully"
}
```

**Response (Error - 400/500):**
```json
{
  "error": "Failed to create user energy"
}
```

---

## Data Models

### Siswa (Student)
```json
{
  "email": "string",
  "nama": "string",
  "id_kelas": "integer",
  "date_created": "timestamp",
  "image_profile": "string"
}
```

### Guru (Teacher)
```json
{
  "email": "string",
  "nama": "string",
  "jabatan": "string",
  "date_created": "timestamp",
  "image_profile": "string"
}
```

### Admin
```json
{
  "email": "string",
  "nama": "string",
  "keterangan": "string",
  "date_created": "timestamp",
  "image_profile": "string"
}
```

### UserVerified
```json
{
  "email": "string",
  "verified_status": "string", // "waiting", "approved", or "rejected"
  "date_created": "timestamp"
}
```

### UserPoint
```json
{
  "email": "string",
  "point": "integer"
}
```

### UserEnergy
```json
{
  "email": "string",
  "energy": "integer",
  "last_updated": "timestamp"
}
```
