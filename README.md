
# 🧾 Microservices Architecture - Payment, Order, Inventory System

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" />
</p>

## 📦 Tentang Proyek

Ini adalah proyek **Microservices** sederhana yang terdiri dari tiga service utama:

- 🔄 **Order Service** – untuk membuat dan mengelola pesanan.
- 💰 **Payment Service** – untuk mencatat dan memproses pembayaran.
- 📦 **Inventory Service** – untuk mengelola ketersediaan stok produk.

Setiap service dibuat menggunakan bahasa **Golang**, menggunakan **PostgreSQL** sebagai database, dan berjalan secara terisolasi menggunakan **Docker Compose**.

---

## 🏗️ Arsitektur

```

\[Client/API Request]
|
v
┌─────────────┐
│  API Gateway│ (opsional - tidak dibahas di sini)
└─────────────┘
|
v
┌─────────────┐     ┌──────────────┐     ┌────────────────┐
│ Order       │ --> │ Payment      │ --> │ Inventory      │
│ Service     │     │ Service      │     │ Service        │
└─────────────┘     └──────────────┘     └────────────────┘
|
PostgreSQL DB (masing-masing service punya DB sendiri)

````

---

## 📌 Langkah Cepat Penggunaan

Ikuti langkah-langkah berikut untuk menjalankan proyek ini secara lokal:

### 1. 💾 Clone / Download Proyek

```bash
git clone https://github.com/mchann/microservices-architecture.git
cd microservices-architecture
````

Atau bisa juga langsung [📥 Download ZIP](https://github.com/mchann/microservices-architecture/archive/refs/heads/main.zip), lalu ekstrak.

---

### 2. 🐳 Jalankan Docker Compose

Pastikan Docker sudah terinstal dan aktif.

```bash
docker-compose up --build
```

Tunggu sampai semua container service (`order`, `payment`, `inventory`, dan database) berhasil berjalan.

---

### 3. ✅ Cek Status Container

```bash
docker ps
```

Pastikan container seperti berikut muncul:

* `order-service`
* `payment-service`
* `inventory-service`
* `order-db`, `payment-db`, `inventory-db` (PostgreSQL untuk masing-masing service)

---

### 4. 🌐 Akses Endpoint API

* Order Service → [http://localhost:8081](http://localhost:8081)
* Payment Service → [http://localhost:8082](http://localhost:8082)
* Inventory Service → [http://localhost:8083](http://localhost:8083)

Gunakan Postman atau `curl` untuk menguji endpoint masing-masing service.

---

### 5. 🧪 Tes Koneksi ke Database (Opsional)

Jika ingin masuk ke dalam container PostgreSQL:

```bash
docker exec -it order-db psql -U postgres
```

Ganti `order-db` dengan `payment-db` atau `inventory-db` sesuai kebutuhan.

---

### 6. 🧹 Stop & Bersihkan

Untuk menghentikan semua container:

```bash
docker-compose down
```

Jika ingin sekaligus menghapus data (volume):

```bash
docker-compose down -v
```

---

## 🗂️ Struktur Direktori

```bash
.
├── order-service/
├── payment-service/
├── inventory-service/
├── docker-compose.yml
└── README.md
```

---

## 🛠 Teknologi yang Digunakan

* [Golang](https://golang.org/) – Bahasa pemrograman utama
* [Docker](https://www.docker.com/) – Untuk containerization
* [PostgreSQL](https://www.postgresql.org/) – Sistem manajemen basis data
* [GORM](https://gorm.io/) – ORM untuk Golang

---

## 📬 Kontribusi

Pull request sangat diterima! Untuk perubahan besar, silakan buka issue terlebih dahulu untuk mendiskusikan apa yang ingin kamu ubah.

---

## 📄 License

Proyek ini dilisensikan di bawah [MIT License](LICENSE).

---

```

---

Kalau nanti kamu mau aku bantu **tambahin bagian dokumentasi API (endpoint + contoh data)** atau **koleksi Postman**, tinggal bilang aja ya.
```
