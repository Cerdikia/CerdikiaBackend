# Redemption Receipt API Documentation

This document provides information about the Redemption Receipt API endpoints, which are used to view and print receipts for gift redemptions.

## Base URL

All endpoints are relative to the base URL of your API server.

## Authentication

All endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Endpoints

### 1. View Redemption Receipt

Returns an HTML page displaying the redemption receipt with all details.

- **URL**: `/view-receipt/:code`
- **Method**: `GET`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | The unique redemption code (e.g., "ABC123XYZ456") |

#### Response

Returns an HTML page with the following information:

- Student information (name, email, class, profile image)
- Redemption details (code, date, status)
- Item details (name, quantity, price in diamonds)
- A button to print the receipt as PDF

#### Example Request

```
GET /view-receipt/ABC123XYZ456
```

#### Error Responses

**Redemption Not Found (404)**
```json
{
  "message": "Redemption with code ABC123XYZ456 not found"
}
```

**Server Error (500)**
```json
{
  "message": "Error retrieving redemption data",
  "error": "Error details"
}
```

### 2. Print Redemption Receipt

Generates and downloads a PDF version of the redemption receipt.

- **URL**: `/print-receipt/:code`
- **Method**: `GET`
- **Auth Required**: Yes

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | The unique redemption code (e.g., "ABC123XYZ456") |

#### Response

Returns a PDF file download with the following information:

- Student information (name, email, class)
- Redemption details (code, date, status)
- Item details (name, quantity, price in diamonds)
- Cerdikia branding and footer

#### Example Request

```
GET /print-receipt/ABC123XYZ456
```

#### Error Responses

**Redemption Not Found (404)**
```json
{
  "message": "Redemption with code ABC123XYZ456 not found"
}
```

**Server Error (500)**
```json
{
  "message": "Error generating PDF",
  "error": "Error details"
}
```

## Usage Examples

### Viewing a Receipt

1. After a successful redemption, note the redemption code from the response
2. Navigate to `/view-receipt/{code}` in a web browser
3. The HTML receipt will be displayed with all details

### Printing a Receipt

1. From the HTML receipt view, click the "Cetak Bukti Penukaran" button
2. Alternatively, navigate directly to `/print-receipt/{code}` in a web browser
3. The browser will download a PDF file named `receipt_{code}_{timestamp}.pdf`

## Receipt Content

### HTML Receipt Includes:

- Cerdikia logo and branding
- Redemption code (prominently displayed)
- Student photo and details
- Item details with quantity and price
- Status indicator (color-coded by status)
- Print button
- Footer with instructions

### PDF Receipt Includes:

- Cerdikia branding
- Redemption code
- Student information section
- Item details section
- Total diamond cost
- Footer with instructions and copyright notice

## Implementation Notes

- The HTML receipt is generated on-the-fly and returned as an HTML response
- The PDF receipt is temporarily stored on the server and automatically deleted after 5 minutes
- Both endpoints retrieve data from multiple tables (logs_penukaran_point, barang, siswa, kelas)
- Status values are color-coded in the HTML view (waiting: orange, completed: green, cancelled: red)
