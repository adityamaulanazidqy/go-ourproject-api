# Go OurProject API

🚧 **Proyek ini masih dalam tahap pengembangan aktif. Fitur dan struktur dapat berubah sewaktu-waktu.** 🚧

Go OurProject API adalah RESTful API yang dibangun menggunakan [Fiber](https://gofiber.io/) (framework web untuk Golang), dengan arsitektur modular, dukungan untuk Redis, dan fitur pengiriman OTP melalui email. Proyek ini bertujuan untuk menjadi backend yang efisien dan terstruktur untuk kebutuhan aplikasi modern.

## Fitur

- Struktur proyek terorganisir dengan arsitektur MVC.
- Otentikasi dan pengiriman OTP melalui email.
- Middleware untuk logging dan Redis client.
- Koneksi ke database menggunakan konfigurasi modular.
- Penanganan rute yang bersih dan terstruktur.

## Teknologi yang Digunakan

- Go (Golang)
- [Fiber](https://gofiber.io/)
- Redis
- PostgreSQL (diasumsikan dari `ConnDB`)
- Logrus untuk logging
- JWT (kemungkinan digunakan dalam otentikasi, berdasarkan struktur umum)

## Struktur Proyek

```
.
├── config/                 # Konfigurasi aplikasi (DB, Redis, Logger)
├── controllers/           # Logika controller seperti OTP Email
├── helpers/               # Fungsi bantu
├── middlewares/          # Middleware untuk otentikasi, logging, dll
├── models/                # Model database
├── repositories/          # Akses ke data layer
├── routes/                # Daftar rute dan endpoint
├── test/                  # File testing
├── main.go                # Entry point aplikasi
├── go.mod                 # File dependensi Go
└── .env                   # Variabel lingkungan
```

## Cara Menjalankan

1. **Clone repositori:**
   ```bash
   git clone https://github.com/adityamaulanazidqy/go-ourproject-api.git
   cd go-ourproject-api
   ```

2. **Buat file `.env` berdasarkan konfigurasi yang dibutuhkan:**
   Contoh variabel lingkungan:
   ```env
   DB_USER=youruser
   DB_PASS=yourpass
   DB_NAME=yourdb
   REDIS_ADDR=localhost:6379
   EMAIL_API_KEY=yourkey
   ```

3. **Jalankan aplikasi:**
   ```bash
   go run main.go
   ```

   Aplikasi akan berjalan di port `8673`.

## Rencana Pengembangan

- [ ] Tambah unit test & integrasi test.
- [ ] Dokumentasi API (Swagger/OpenAPI).
- [ ] CI/CD pipeline (GitHub Actions).
- [ ] Manajemen user dan otorisasi lanjutan.

## Kontribusi

Pull request sangat diterima! Untuk perubahan besar, harap buka isu terlebih dahulu untuk mendiskusikan apa yang ingin Anda ubah.

## Lisensi

Proyek ini menggunakan lisensi MIT. Silakan lihat `LICENSE` untuk informasi lebih lanjut.
