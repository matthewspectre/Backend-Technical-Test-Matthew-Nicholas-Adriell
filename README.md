Project Overview
Aplikasi ini dirancang untuk mengelola proses inventory dan procurement, yang mencakup pengelolaan produk, supplier, warehouse, Purchase Request (PR), Purchase Order (PO), 
serta proses penerimaan barang (Goods Receipt).


## 🛠️ Tech Stack

| Technology                 | Description                                                                          |
| -------------------------- | ------------------------------------------------------------------------------------ |
| 🐹 **Golang**              | Bahasa pemrograman utama untuk backend                                               |
| 🌐 **Gin**                 | Framework untuk membangun RESTful API                                                |
| 🐘 **PostgreSQL**          | Relational database untuk menyimpan data aplikasi                                    |
| 🔗 **GORM**                | ORM untuk mengelola interaksi dengan database                                        |
| 🔐 **JWT**                 | Authentication dan authorization berbasis token                                      |
| 🏗️ **Clean Architecture** | Arsitektur untuk memisahkan business logic, use case, repository, dan delivery layer |
| 📡 **REST API**            | Interface komunikasi antara client dan backend                                       |
| 🧪 **Go Testing**          | Automated testing untuk memastikan business logic berjalan sesuai kebutuhan          |
| 📬 **Postman**             | Pengujian dan validasi REST API                                                      |
| 🐙 **Git & GitHub**        | Version control dan repository management                                            |

### 🏛️ Architecture
Project ini menerapkan **Clean Architecture** dengan struktur layer sebagai berikut:

```text
Clean Architecture Golang:
1. Entity
    2. Repository
        3. Model
            4. mysql
                5. usecase
                    6. handler
                        7. router
                            8. main
Struktur layer:

1. **Entity** — Mendefinisikan struktur dan aturan dasar data/domain.
2. **Repository** — Mendefinisikan interface untuk akses dan pengelolaan data.
3. **Model** — Merepresentasikan struktur data yang digunakan untuk database.
4. **PostgreSQL** — Mengimplementasikan akses database melalui repository.
5. **Usecase** — Menangani business logic dan alur proses aplikasi.
6. **Handler** — Menangani request dan response HTTP/API.
7. **Router** — Mendefinisikan endpoint dan menghubungkan request ke handler.
8. **Main** — Melakukan inisialisasi aplikasi, dependency, database, router, dan menjalankan server.

Alur aplikasi:
Entity → Repository → Model → Database → Usecase → Handler → Router → Main
```

## Database Design

```text
1. Master Data
-------------------------------------------------------------------------
**users**
Menyimpan data pengguna sistem.
Role yang tersedia:

USER
APPROVER
User digunakan untuk mencatat siapa yang membuat purchase request dan siapa yang menerima barang.

**product**
Menyimpan data barang yang tersedia.

Kolom penting:

- id: primary key
- sku: kode unik barang
- name: nama barang
- unit: satuan barang
- is_active: status aktif barang
- created_at dan updated_at: informasi waktu data dibuat dan diubah
SKU dibuat UNIQUE agar satu barang tidak memiliki kode yang sama.

**supplier**
Menyimpan data pemasok barang.

Relasinya digunakan oleh tabel purchase_orders melalui supplier_id.

**warehouse**
Menyimpan data gudang yang ada.

Relasinya digunakan oleh:

- inventory
- purchase_requests
- purchase_orders
- goods_receipts
- inventory_movements
-------------------------------------------------------------------------

2. Purchase Request
**purchase_requests** menyimpan permintaan barang sebelum purchase order.
Status yang diperbolehkan:
- DRAFT
- SUBMITTED
- APPROVED
- REJECTED

Setiap purchase request memiliki banyak barang melalui tabel purchase_request_items.

**purchase_request_items**
Tabel ini menjadi detail barang yang diminta.
Relasi :
- satu purchase request memiliki banyak item
- satu product dapat muncul di banyak purchase request

**quantity** harus lebih besar dari nol.
-------------------------------------------------------------------------

3. Purchase Order
**purchase_orders** menyimpan pesanan pembelian yang dibuat berdasarkan purchase request.
relasi :
- purchase_request_id: request sumber
  purchase_request_id diberi constraint UNIQUE, sehingga satu purchase request hanya dapat memiliki satu purchase order.
- supplier_id: supplier yang dipilih
- warehouse_id: gudang tujuan

Status purchase order:
- DRAFT
- ORDERED
- PARTIALLY_RECEIVED
- RECEIVED
- CANCELLED

**purchase_order_items**
Menyimpan detail barang dalam purchase order.

memiliki informasi kuantitas:
- ordered_quantity: jumlah yang dipesan
- received_quantity: jumlah yang sudah diterima

Constraint:
- ordered_quantity > 0
- received_quantity >= 0
- kombinasi purchase_order_id dan product_id harus unik, sehinga satu produk tidak boleh muncul dua kali dalam purchase order yang sama.
-------------------------------------------------------------------------

4. Goods Receipt
**goods_receipts** mencatat penerimaan barang dari supplier.

Status goods receipt:
- POSTED
- CANCELLED

**goods_receipt_items**
Menyimpan detail barang yang diterima.

Relasi:
- goods receipt
- purchase order item
- product

Constraint:
- UNIQUE (goods_receipt_id, product_id) mencegah produk yang sama dicatat dua kali dalam satu penerimaan.
-------------------------------------------------------------------------

5. Inventory
**inventory** menyimpan stok barang berdasarkan kombinasi barang dan gudang.

Relasi:
- product_id
- warehouse_id

Constraint:
- UNIQUE (product_id, warehouse_id), sehingga satu produk hanya memiliki satu record stok untuk setiap gudang.
-------------------------------------------------------------------------

6. Inventory Movement
**inventory_movements** berfungsi sebagai histori perubahan stok.

- Informasi yang dicatat:
- gudang
- produk
- jenis pergerakan
- jumlah
- referensi transaksi
- waktu perubahan

 jenis movement dilakukan melalui PURCHASE_RECEIPT


```
<img width="775" height="795" alt="image" src="https://github.com/user-attachments/assets/03e1d564-8100-4c51-bdac-79984230a906" />

```text
**Alur Bisnis Database :**
1. User membuat purchase_request.
2. Purchase request memiliki detail list produk pada purchase_request_items.
3. Setelah disetujui, dibuat satu purchase_order.
4. Purchase order memiliki detail pada purchase_order_items.
5. Supplier mengirim barang.
6. Barang dicatat melalui goods_receipts dan goods_receipt_items.
7. Update stok pada inventory sesuai barang yang diterima saat goods_receipt dibuat.
8. Perubahan stok dicatat pada inventory_movements.
```

# Cara menjalankan project dari awal
1. download Dbeaver melalui link berikut : https://dbeaver.io/download/
2. download PostgreSQL melalui link berikut : https://www.postgresql.org/download/
3. download Golang melalui link berikut : https://go.dev/doc/install
4. Buat database PostgreSQL
   - Buka DBeaver dan buat koneksi ke PostgreSQL.
   - Gunakan konfigurasi database berikut:
     Host: `localhost`
     Port: `5432`
     Database: `be_2evindo`
     Username: `postgres`
     Password: password PostgreSQL pribadi
   - dropdown postgres dan pada folder 'Databases' klik kanan lalu pilih 'Create New Database dengan nama " be_2evindo'
   - Klik kanan database `be_2evindo`, lalu pilih **SQL Editor > New SQL Script**.
   - Buka file migration dari folder `migrations/test`.
   - Jalankan / copy query sesuai urutan berikut, mulai dari `001` sampai `012` (bisa menggunakan ctrl + enter pada dbeaver):
     `001_product_test.sql`
     `002_create_users.sql`
     `003_create_suppliers.sql`
     `004_create_warehouses.sql`
     `005_create_inventories.sql`
     `006_create_purchase_requests.sql`
     `007_create_purchase_orders.sql`
     `008_create_purchase_order_items.sql`
     `009_update_purchase_order_status_constraint.sql`
     `010_create_goods_receipts.sql`
     `011_create_goods_receipt_items.sql`
     `012_create_inventory_movements.sql`
     <img width="1237" height="776" alt="image" src="https://github.com/user-attachments/assets/36b83314-3330-4225-8255-8b5e9c8466bc" />
   - Setelah clone repository, untuk menjalankan file Go, klik 'Run' pada pojok kanan, lalu start debugging (tpmbol alternatif = F5)
     <img width="1887" height="562" alt="image" src="https://github.com/user-attachments/assets/bac5b2a4-ea93-475e-97b7-0d7b896ecf25" />
     
   - server berjalan pada http://localhost:8080
     <img width="1384" height="235" alt="image" src="https://github.com/user-attachments/assets/707d2562-b5c3-47f1-a20c-d76dda5846b7" />

# Cara menjalankan url API melalui postman 
Modul **user** :
- create user : POST / http://localhost:8080/auth/register
PAYLOAD : {
  "name": "matthew",
  "email": "matthew@gmail.com",
  "password": "password123",
  "role": "APPROVER"
}
<img width="382" height="413" alt="image" src="https://github.com/user-attachments/assets/e5f83162-30c5-4e03-8248-5913d073889d" />

- login : POST / http://localhost:8080/login
  PAYLOAD : {
  "email": "ikhsan@gmail.com",
  "password": "password123"
}
<img width="1387" height="466" alt="image" src="https://github.com/user-attachments/assets/0d8c1a75-6783-4aae-91e9-77bf479ad5f5" />
Token digunakan untuk autentikasi, masukan token di postman pada bagian authorization, dan pilih bearer token
-------------------------------------------------------------------------

Modul **product**
- create product : POST / http://localhost:8080/products/
  PAYLOAD : {
  "sku": "HKU-001",
  "name": "Throthle Body HKU 230CC",
  "unit": "pcs"
}
<img width="361" height="471" alt="image" src="https://github.com/user-attachments/assets/99d8475b-f76b-4f7b-acad-7c12cba061ad" />

- get all product : GET / http://localhost:8080/products/
  <img width="788" height="882" alt="image" src="https://github.com/user-attachments/assets/1cc927a0-21c8-4ef1-bf3a-6fb5b2282ade" />

- get product by ID : GET / http://localhost:8080/products/:id
  <img width="702" height="468" alt="image" src="https://github.com/user-attachments/assets/75c55cd2-b7e5-496d-9248-60bc1ffc6fc5" />

- update product : PATCH / http://localhost:8080/products/:id
  PAYLOAD : {
  "sku": "HKU-008",
  "name": "Knalpot Racing HKU",
  "unit": "pcs"
}
<img width="443" height="478" alt="image" src="https://github.com/user-attachments/assets/e0876865-e7a9-4b9c-a0ac-9e3f725f5c52" />




