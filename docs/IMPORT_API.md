# API Documentation: Import Excel Endpoint

## Endpoint

`POST /import`

## Description
Upload an Excel file (.xlsx) containing questions (soal), modules, subjects (mapel), classes (kelas), and optionally images. The API will parse the file and return a structured JSON response without saving any data to the database.

## Request
- **Content-Type:** `multipart/form-data`
- **Form Field:** `file` (the Excel file to upload)

### Example Request (cURL)
```bash
curl -X POST http://localhost/import \
  -F "file=@/path/to/your/file.xlsx"
```

## Response
- **Content-Type:** `application/json`
- **Status:** `200 OK` (on success)

### Response Structure
The response is a JSON array grouped by subject (`mapel`) and class (`kelas`). Each group contains modules, and each module contains an array of questions. If a question or option contains an image, it will be returned as an HTML `<img>` tag referencing the static `/data/images/` path.

#### Example Response
```json
[
  {
    "mapel": "Matematika",
    "kelas": "6",
    "module": [
      {
        "judul_module": "Modul 1",
        "deskripsi_module": "Deskripsi Modul 1",
        "soal": [
          {
            "soal": "<img src=\"/data/images/soal1.png\">Soal dengan gambar</img>",
            "jenis": "PG",
            "opsi_a": "Pilihan A",
            "opsi_b": "Pilihan B",
            "opsi_c": "<img src=\"/data/images/opsiC.png\">Pilihan C</img>",
            "opsi_d": "Pilihan D",
            "jawaban": "a"
          }
        ]
      }
    ]
  }
]
```

## Field Explanation
- `mapel`: Subject name
- `kelas`: Class/grade
- `module`: Array of modules for the subject/class
  - `judul_module`: Module title
  - `deskripsi_module`: Module description
  - `soal`: Array of questions
    - `soal`: Question text (may include HTML for images)
    - `jenis`: Question type (e.g., `PG` for multiple choice)
    - `opsi_a` ... `opsi_d`: Option texts (may include HTML for images)
    - `jawaban`: The correct answer (e.g., `a`)

## Error Responses
- `400 Bad Request`: No file uploaded or file is not a valid Excel file.
- `500 Internal Server Error`: Failed to process the file.

### Example Error Response
```json
{
  "error": "File tidak valid atau gagal diproses."
}
```

## Notes
- Images in Excel will be extracted and referenced as `/data/images/<filename>` in the response.
- No data will be saved to any database.
- The endpoint only returns the parsed data as JSON.
