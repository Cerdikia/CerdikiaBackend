# Gift Redemption API Documentation

This document provides information about the Gift Redemption API endpoints, which are used to redeem gifts using diamonds and manage redemption logs.

## Base URL

All endpoints are relative to the base URL of your API server.

## Authentication

All endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Endpoints

### 1. Redeem Gifts

Allows a user to redeem one or more gifts using their diamond balance.

- **URL**: `/redeem-gifts`
- **Method**: `POST`
- **Auth Required**: Yes

#### Request Body

```json
{
  "email": "student@example.com",
  "items": [
    {
      "id_barang": 1,
      "jumlah": 2
    },
    {
      "id_barang": 3,
      "jumlah": 1
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `email` | string | The email of the student redeeming the gifts |
| `items` | array | List of items to redeem |
| `items[].id_barang` | integer | ID of the item to redeem |
| `items[].jumlah` | integer | Quantity of the item to redeem |

#### Response

**Success (200)**
```json
{
  "message": "Redemption successful",
  "redemption_code": "ABC123XYZ456",
  "items": [
    {
      "id_barang": 1,
      "nama_barang": "Item Name",
      "jumlah": 2,
      "diamond": 100,
      "total_diamond": 200
    },
    {
      "id_barang": 3,
      "nama_barang": "Another Item",
      "jumlah": 1,
      "diamond": 50,
      "total_diamond": 50
    }
  ],
  "total_diamond_used": 250,
  "remaining_diamond": 750
}
```

#### Error Responses

**Invalid Input (400)**
```json
{
  "message": "Email is required"
}
```

**User Not Verified (403)**
```json
{
  "message": "User is not verified"
}
```

**Insufficient Diamonds (400)**
```json
{
  "message": "Insufficient diamonds",
  "required": 250,
  "available": 200
}
```

**Server Error (500)**
```json
{
  "message": "Error processing redemption",
  "error": "Error details"
}
```

### 2. Get All Redemptions

Retrieves a list of all redemption logs with optional filtering.

- **URL**: `/redemptions`
- **Method**: `GET`
- **Auth Required**: Yes

#### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `email` | string | Optional. Filter by student email |
| `status` | string | Optional. Filter by status ("menunggu", "selesai", "dibatalkan") |
| `id_barang` | integer | Optional. Filter by item ID |

#### Response

**Success (200)**
```json
{
  "redemptions": [
    {
      "id_log": 1,
      "email": "student@example.com",
      "id_barang": 1,
      "jumlah": 2,
      "tanggal_penukaran": "2025-05-13T15:30:45Z",
      "kode_penukaran": "ABC123XYZ456",
      "status_penukaran": "menunggu",
      "nama_barang": "Item Name",
      "img": "item_image_url.jpg",
      "description": "Item description",
      "diamond": 100,
      "nama_siswa": "Student Name"
    },
    {
      "id_log": 2,
      "email": "another@example.com",
      "id_barang": 3,
      "jumlah": 1,
      "tanggal_penukaran": "2025-05-12T10:15:30Z",
      "kode_penukaran": "DEF456GHI789",
      "status_penukaran": "selesai",
      "nama_barang": "Another Item",
      "img": "another_image_url.jpg",
      "description": "Another description",
      "diamond": 50,
      "nama_siswa": "Another Student"
    }
  ],
  "count": 2
}
```

**Empty Result (200)**
```json
{
  "redemptions": [],
  "count": 0
}
```

**Server Error (500)**
```json
{
  "message": "Error retrieving redemption logs",
  "error": "Error details"
}
```

### 3. Get Redemption by ID

Retrieves a specific redemption log by its ID.

- **URL**: `/redemptions/:id`
- **Method**: `GET`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The ID of the redemption log |

#### Response

**Success (200)**
```json
{
  "id_log": 1,
  "email": "student@example.com",
  "id_barang": 1,
  "jumlah": 2,
  "tanggal_penukaran": "2025-05-13T15:30:45Z",
  "kode_penukaran": "ABC123XYZ456",
  "status_penukaran": "menunggu",
  "nama_barang": "Item Name",
  "img": "item_image_url.jpg",
  "description": "Item description",
  "diamond": 100,
  "nama_siswa": "Student Name"
}
```

**Redemption Not Found (404)**
```json
{
  "message": "Redemption log with ID 1 not found"
}
```

**Server Error (500)**
```json
{
  "message": "Error retrieving redemption log",
  "error": "Error details"
}
```

### 4. Get Redemption by Code

Retrieves a specific redemption log by its unique redemption code.

- **URL**: `/redemptions/code/:code`
- **Method**: `GET`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | The unique redemption code (e.g., "ABC123XYZ456") |

#### Response

**Success (200)**
```json
{
  "id_log": 1,
  "email": "student@example.com",
  "id_barang": 1,
  "jumlah": 2,
  "tanggal_penukaran": "2025-05-13T15:30:45Z",
  "kode_penukaran": "ABC123XYZ456",
  "status_penukaran": "menunggu",
  "nama_barang": "Item Name",
  "img": "item_image_url.jpg",
  "description": "Item description",
  "diamond": 100,
  "nama_siswa": "Student Name"
}
```

**Redemption Not Found (404)**
```json
{
  "message": "Redemption log with code ABC123XYZ456 not found"
}
```

**Server Error (500)**
```json
{
  "message": "Error retrieving redemption log",
  "error": "Error details"
}
```

### 5. Update Redemption Status

Updates the status of a redemption log.

- **URL**: `/redemptions/:id/status`
- **Method**: `PUT`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The ID of the redemption log |

#### Request Body

```json
{
  "status": "selesai"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | The new status ("menunggu", "selesai", "dibatalkan") |

#### Response

**Success (200)**
```json
{
  "message": "Redemption status updated successfully",
  "id_log": 1,
  "status": "selesai"
}
```

**Redemption Not Found (404)**
```json
{
  "message": "Redemption log with ID 1 not found"
}
```

**Invalid Status (400)**
```json
{
  "message": "Invalid status. Must be one of: menunggu, selesai, dibatalkan"
}
```

**Server Error (500)**
```json
{
  "message": "Error updating redemption status",
  "error": "Error details"
}
```

### 6. Delete Redemption

Deletes a redemption log.

- **URL**: `/redemptions/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The ID of the redemption log |

#### Response

**Success (200)**
```json
{
  "message": "Redemption log deleted successfully",
  "id_log": 1
}
```

**Redemption Not Found (404)**
```json
{
  "message": "Redemption log with ID 1 not found"
}
```

**Server Error (500)**
```json
{
  "message": "Error deleting redemption log",
  "error": "Error details"
}
```

## Data Models

### Redemption Log

| Field | Type | Description |
|-------|------|-------------|
| `id_log` | integer | Unique identifier for the redemption log |
| `email` | string | Email of the student who redeemed the item |
| `id_barang` | integer | ID of the redeemed item |
| `jumlah` | integer | Quantity of the item redeemed |
| `tanggal_penukaran` | datetime | Date and time of the redemption |
| `kode_penukaran` | string | Unique redemption code for verification |
| `status_penukaran` | string | Status of the redemption ("menunggu", "selesai", "dibatalkan") |

### Item (Barang)

| Field | Type | Description |
|-------|------|-------------|
| `id_barang` | integer | Unique identifier for the item |
| `nama_barang` | string | Name of the item |
| `img` | string | URL to the item's image |
| `description` | string | Description of the item |
| `diamond` | integer | Diamond cost of the item |

## Workflow

1. **Redeeming Gifts**:
   - User submits a redemption request with their email and items to redeem
   - System verifies the user's account and diamond balance
   - If successful, a unique redemption code is generated
   - The redemption is recorded with status "menunggu" (waiting)

2. **Managing Redemptions**:
   - Administrators can view all redemptions or filter by various criteria
   - Specific redemptions can be viewed by ID or redemption code
   - Administrators can update the status to "selesai" (completed) when the item is delivered
   - Administrators can update the status to "dibatalkan" (cancelled) if the redemption cannot be fulfilled
   - Redemptions can be deleted if necessary

3. **Verification**:
   - The unique redemption code serves as proof of redemption
   - Users can view their redemption receipt using the code
   - Administrators can verify the redemption by checking the code
