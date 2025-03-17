# BAAK API

An unofficial API for BAAK.

## Disclaimer

Project ini ga ada hubungannya sama Universitas Gunadarma atau BAAK ya. Ini cuma buat belajar aja, jangan dipake buat yang aneh-aneh. Developer ga bertanggung jawab kalo ada yang nyalahgunain ini.

## Fitur

- Pencarian Jadwal Kuliah
- Kalender Akademik
- Informasi Kelas Baru
- Jadwal UTS
- Informasi Mahasiswa Baru
- Rate limiting
- Dukungan CORS
- Monitoring kesehatan
- Format error yang terstandarisasi

## Endpoint API

### Health Check

```
GET /health
```

Mengembalikan status kesehatan API.

### Jadwal Kuliah

```
GET /jadwal/{kelas}
```

Mendapatkan informasi jadwal untuk kelas tertentu.

Parameter:

- `kelas` (path parameter): Kode kelas (minimal 3 karakter)

### Kalender Akademik

```
GET /kalender
```

Mendapatkan informasi kalender akademik.

### Informasi Kelas Baru

```
GET /kelasbaru/{kelas}
```

Mendapatkan informasi tentang kelas baru.

Parameter:

- `kelas` (path parameter): Kode kelas

### Jadwal UTS

```
GET /uts/{kelas}
```

Mendapatkan jadwal UTS (Ujian Tengah Semester) untuk kelas tertentu.

Parameter:

- `kelas` (path parameter): Kode kelas

### Informasi Mahasiswa Baru

```
GET /mahasiswabaru/{npm}
```

Mendapatkan informasi untuk mahasiswa baru.

Parameter:

- `npm` (path parameter): Nomor Pokok Mahasiswa

## Format Response

Semua response mengikuti format ini:

```json
{
  "success": true,
  "data": {
    // Data response di sini
  }
}
```

Response error:

```json
{
  "success": false,
  "error": "Pesan error di sini"
}
```

## Rate Limiting

API ini menggunakan rate limiting untuk mencegah penyalahgunaan. Secara default, mengizinkan 60 request per menit per alamat IP.

## Konfigurasi

API bisa dikonfigurasi menggunakan environment variables:

- `PORT`: Port server (default: ":8080")
- `BASE_URL`: URL dasar website BAAK (default: "https://baak.gunadarma.ac.id")
- `RATE_LIMIT_PER_MIN`: Batas rate per menit (default: 60)
- `ALLOWED_ORIGINS`: Daftar origin CORS yang diizinkan, dipisahkan dengan koma (default: "\*")

## Development

### Prasyarat

- Go 1.16 atau lebih tinggi
- Git

### Setup

1. Clone repository:

```bash
git clone https://github.com/yourusername/baak-api.git
cd baak-api
```

2. Install dependencies:

```bash
go mod download
```

3. Jalankan server:

```bash
go run api/index.go
```

## To-Do

- [x] Jadwal
- [x] Kalender Akademik
- [x] Mahasiswa Baru
- [x] Mahasiswa Kelas 2 Baru
- [x] UTS
- [ ] UU
- [ ] UAS

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
