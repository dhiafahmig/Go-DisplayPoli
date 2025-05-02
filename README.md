# Go-DisplayPoli

Aplikasi DisplayPoli untuk menampilkan antrian poli rumah sakit berbasis Go.

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
git clone https://github.com/your-username/go-displaypoli.git
cd go-displaypoli
```

2. Perbarui modul Go:
```
# Jika Anda menggunakan Go Modules
go mod edit -module github.com/your-username/go-displaypoli
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

PUSHER_APP_ID=yourappid
PUSHER_APP_KEY=yourappkey
PUSHER_APP_SECRET=yourappsecret
PUSHER_APP_CLUSTER=ap1
```

4. Jalankan aplikasi:
```
go run main.go
```

## Struktur Direktori

```
app/
  handlers/    # Handler HTTP
    common.go   # Struktur data umum
    displaypoli.go    # Handler display poli
    jadwaldokter.go   # Handler jadwal dokter
    panggilpoli.go    # Handler panggil pasien
    settingdisplaypoli.go   # Handler pengaturan display
    settingpoli.go          # Handler pengaturan poli
    settingposisidokter.go  # Handler pengaturan posisi dokter
  models/      # Model data
    models.go   # Definisi model database
  services/    # Layanan aplikasi
    daylist.go  # Service untuk daftar hari
    env.go      # Service untuk variabel lingkungan
config/        # Konfigurasi
templates/     # Template HTML
  displaypoli.html   # Tampilan display poli
  panggilpoli.html   # Tampilan panggil pasien
assets/        # Aset statis (CSS, JS)
main.go        # Entry point
go.mod         # Deklarasi modul Go
.env           # Konfigurasi lingkungan (buat sesuai kebutuhan)
```

## Endpoint API

### Halaman
- `GET /display/:kd_display` - Menampilkan display poli
- `GET /settings/display` - Pengaturan display
- `GET /settings/poli` - Pengaturan poli
- `GET /settings/dokter` - Pengaturan posisi dokter
- `GET /jadwal/dokter` - Jadwal dokter
- `GET /panggilpoli/:kd_ruang_poli/:kd_display` - Panggil pasien di poli
- `GET /ws/:kd_display` - WebSocket untuk notifikasi real-time

### API
- `POST /api/panggilpoli` - Memanggil pasien
- `POST /api/log` - Update log antrian pasien
- `POST /api/log/reset/:no_rawat` - Reset log antrian pasien
- `POST /api/display` - Tambah display poli
- `PUT /api/display` - Edit display poli
- `DELETE /api/display/:kd_display` - Hapus display poli
- `POST /api/poli` - Tambah poli
- `PUT /api/poli` - Edit poli
- `DELETE /api/poli/:kd_ruang_poli` - Hapus poli
- `POST /api/dokterpoli` - Edit posisi dokter
- `GET /api/jadwal/dokter` - Cari dokter
- `POST /api/jadwal` - Tambah jadwal dokter
- `PUT /api/jadwal` - Edit jadwal dokter
- `DELETE /api/jadwal` - Hapus jadwal dokter

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

## Troubleshooting

Jika Anda mengalami masalah dengan WebSocket:
- Pastikan URL WebSocket benar (ws:// untuk HTTP dan wss:// untuk HTTPS)
- Periksa port dan firewall Anda

Jika Anda mengalami masalah dengan koneksi database:
- Verifikasi kredensial database di file `.env`
- Pastikan database MySQL berjalan

## Lisensi

MIT