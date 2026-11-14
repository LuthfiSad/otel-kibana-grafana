# Massive Enterprise Observability Platform (Go + PostgreSQL + Full Observability)

Platform aplikasi backend enterprise berskala besar (*massive scale*) yang dilengkapi dengan **Full Observability Stack** (OpenTelemetry, Elasticsearch, Kibana, Prometheus, Grafana, dan Jaeger) serta database relasional multi-tabel (`testing-otel`) di PostgreSQL eksternal.

---

## 🏗️ Arsitektur & Komponen Stack

1. **Aplikasi Go (Gin Framework + OpenTelemetry SDK)**
   - Port: `8080`
   - Melakukan otomatisasi migrasi schema & *automated massive seeder* saat pertama kali startup.
   - Mengirimkan *Distributed Tracing* ke OpenTelemetry Collector (`otel-collector:4317`).
   - Menyediakan endpoint RESTful API terstruktur.

2. **Database PostgreSQL Eksternal (`local-postgres`)**
   - Berjalan mandiri di luar docker-compose (`postgres:15-alpine`, port `5432`, user `postgres`, db `testing-otel`).
   - Terhubung dengan Docker network `otel-kibana-grafana-test_default`.

3. **Full Observability Stack**
   - **OpenTelemetry Collector**: Menerima trace dan meneruskannya ke Jaeger.
   - **Jaeger UI**: Distributed Tracing visualization (`http://localhost:16686`).
   - **Prometheus**: Metrics collection & scraping (`http://localhost:9090`).
   - **Grafana**: Dashboard visualisasi metrik (`http://localhost:3000`).
   - **Elasticsearch & Kibana**: Log aggregation & analytics (`http://localhost:5601`).

---

## 📊 Skema Database Relasional Enterprise (10+ Tabel)
1. `tenants`: Multi-tenant management.
2. `users`: Manajemen pengguna sistem (500 data seed).
3. `categories`: Kategori produk (10 data seed).
4. `products`: Katalog produk enterprise (300 data seed).
5. `warehouses`: Gudang logistik (2 data seed).
6. `inventory_stocks`: Stok produk per gudang.
7. `orders`: Transaksi pesanan (1000 data seed).
8. `order_items`: Detail item pesanan.
9. `payments`: Gerbang pembayaran & status transaksi (1000 data seed).
10. `shipments`: Pengiriman logistik & nomor resi (1000 data seed).
11. `audit_logs`: Catatan audit keamanan sistem.

---

## 🚀 Daftar Endpoint RESTful API (`/api/v1`)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/health` | Pemeriksaan kesehatan sistem & timestamp. |
| `GET` | `/api/v1/tenants` | Menampilkan daftar tenant aktif. |
| `GET` | `/api/v1/users` | Menampilkan daftar pengguna (pagination/limit). |
| `GET` | `/api/v1/products` | Menampilkan katalog produk enterprise. |
| `GET` | `/api/v1/warehouses` | Menampilkan daftar gudang logistik. |
| `GET` | `/api/v1/orders` | Menampilkan transaksi pesanan pelanggan. |
| `GET` | `/api/v1/payments` | Menampilkan data pembayaran transaksi. |
| `GET` | `/api/v1/shipments` | Menampilkan status pengiriman logistik & resi. |
| `GET` | `/api/v1/audit-logs` | Menampilkan audit log aktivitas sistem. |
| `GET` | `/api/v1/analytics/summary` | Ringkasan statistik (total user, produk, order, & total revenue). |

---

## 🔍 Panduan Lengkap Cara Melihat & Menggunakan Observability

Berikut adalah langkah-langkah detail cara memantau dan melihat hasil observability dari aplikasi enterprise ini:

### 1. Distributed Tracing (Jaeger)
* **URL:** `http://localhost:16686`
* **Cara Melihat Tracing:**
  1. Buka browser dan akses `http://localhost:16686`.
  2. Pada panel sebelah kiri, pilih Service: `massive-enterprise-app`.
  3. Klik **Find Traces**.
  4. Anda dapat melihat waktu eksekusi setiap request HTTP (seperti `/health` atau `/api/v1/analytics/summary`), durasi pemanggilan database, hingga hierarki *spans* secara real-time.

### 2. Metrics & Performance Monitoring (Prometheus & Grafana)
* **Prometheus URL:** `http://localhost:9090`
* **Grafana URL:** `http://localhost:3000` (Login: `admin` / `admin`)
* **Cara Melihat Metrik:**
  1. Buka Grafana di `http://localhost:3000` dan login menggunakan `admin`/`admin`.
  2. Tambahkan Data Source Prometheus (`http://prometheus:9090`).
  3. Buat atau impor dashboard untuk memantau metrik penggunaan CPU, memori, HTTP request rate (`gin_requests_total`), dan latensi dari kontainer aplikasi.

### 3. Log Aggregation & Analytics (Elasticsearch & Kibana)
* **Kibana URL:** `http://localhost:5601`
* **Cara Melihat Log:**
  1. Buka Kibana di `http://localhost:5601`.
  2. Masuk ke menu **Discover** atau **Index Management**.
  3. Anda dapat memantau log aplikasi secara terpusat yang dikirimkan oleh OpenTelemetry Collector ke Elasticsearch.

---

## 🛠️ Cara Menjalankan & Verifikasi

1. **Jalankan PostgreSQL Eksternal:**
   ```bash
   docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=*** POSTGRES_DB=testing-otel postgres:15-alpine
   docker network connect otel-kibana-grafana-test_default local-postgres
   ```

2. **Jalankan Stack dengan Docker Compose:**
   ```bash
   docker compose up --build -d
   ```

3. **Verifikasi API & Analytics:**
   ```bash
   curl -s http://localhost:8080/health
   curl -s http://localhost:8080/api/v1/analytics/summary
   ```
