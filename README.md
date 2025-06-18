# Go-DisplayPoli

BACKEND untuk menampilkan antrian poli rumah sakit berbasis Go.

## Fitur

- Menampilkan antrian pasien di poli
- Mengatur display poli
- Mengatur posisi dokter
- Mengatur jadwal dokter
- Memanggil pasien dengan notifikasi real-time

## Teknologi

- Go
- Gin Web Framework
- GORM ORM
- WebSocket untuk komunikasi real-time
- MySQL Database

## Persyaratan

- Go 1.19+
- MySQL 5.7+

## Instalasi

1. Clone repositori:
```
git clone https://github.com/dhiafahmig/go-displaypoli.git
cd go-displaypoli
```

2. Perbarui modul Go:
```
# Jika Anda menggunakan Go Modules
go mod edit -module github.com/dhiafahmig/go-displaypoli
go mod tidy
```

3. Buat file .env dengan mengisi konfigurasi database Anda:
```
APP_NAME=DisplayPoli
APP_URL=http://localhost:8080

DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=dbname
DB_USERNAME=username
DB_PASSWORD=password

```

4. Jalankan aplikasi:
```
go run main.go
```

## Database

Aplikasi ini menggunakan skema database yang sama dengan versi PHP, dengan tabel-tabel sebagai berikut:

- `bw_display_poli` - Display poli
- `bw_ruang_poli` - Ruang poli
- `bw_ruangpoli_dokter` - Relasi dokter dan ruang poli
- `bw_jadwal_dokter` - Jadwal praktek dokter
- `bw_log_antrian_poli` - Log antrian poli
- `dokter` - Data dokter
- `poliklinik` - Data poliklinik
- `pasien` - Data pasien
- `reg_periksa` - Registrasi pemeriksaan
- `penjab` - Penjamin

## Panduan Migrasi dari PHP

Jika Anda sebelumnya menggunakan versi PHP dari aplikasi ini, Anda dapat melakukan migrasi dengan langkah-langkah berikut:

1. Pastikan database yang digunakan sama dan telah diperbarui
2. Sesuaikan konfigurasi pada file `.env` dengan database Anda
3. Jalankan aplikasi Go-DisplayPoli dan verifikasi fungsionalitasnya
