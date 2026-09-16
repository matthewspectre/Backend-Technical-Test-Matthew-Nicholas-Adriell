## Project Overview
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

   - jika setelah run, masukan launch json nya berikut :
   `{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API (PostgreSQL :8080)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api",
            "buildFlags": "-p=1 -buildvcs=false",
            "env": {
                "APP_PORT": "8080",
                "DB_USER": "postgres",
                "DB_PASS": "lopezrica123",
                "DB_HOST": "localhost",
                "DB_PORT": "5432",
                "DB_NAME": "be_2evindo",
                "JWT_SECRET": "be-evindo-development-secret"
            }
        },
        {
            "name": "Launch API (PostgreSQL :8081)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api",
            "buildFlags": "-p=1 -buildvcs=false",
            "env": {
                "APP_PORT": "8081",
                "DB_USER": "postgres",
                "DB_PASS": "lopezrica123",
                "DB_HOST": "localhost",
                "DB_PORT": "5432",
                "DB_NAME": "be_2evindo",
                "JWT_SECRET": "be-evindo-development-secret"
            }
        }
    ]
}
`
<img width="640" height="755" alt="image" src="https://github.com/user-attachments/assets/ee527d11-92bd-4b8f-a59f-ef716840c3bc" />

     
   - server berjalan pada http://localhost:8080
     <img width="1384" height="235" alt="image" src="https://github.com/user-attachments/assets/707d2562-b5c3-47f1-a20c-d76dda5846b7" />

# Cara menjalankan url API melalui postman 
link postman : https://technical-meteorologist-7314860-s-team.postman.co/workspace/My-Workspace~a8b21373-9507-4f2c-87da-6f0b6006a7b7/collection/41742420-5b875c99-943c-4f2f-85ba-bf59685accb9?action=share&source=copy-link&creator=41742420

Modul **user** :
- create user : POST / http://localhost:8080/auth/register
PAYLOAD :{
  "name": "matthew",
  "email": "matthew@gmail.com",
  "password": "password123",
  "role": "APPROVER"
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

- unactive product : PATCH / http://localhost:8080/products/16/deactivate
  <img width="785" height="392" alt="image" src="https://github.com/user-attachments/assets/5cc0eede-0f05-4e4b-b603-74a2d47549eb" />
-------------------------------------------------------------------------

Modul **supplier** 
- create supplier : POST / http://localhost:8080/suppliers/
  PAYLOAD : {
  "company_name": "PT Mencari Tenar",
  "name": "Brando Jawato",
  "email": "brdjwt21@gmail.com",
  "phone": "0812345623123789",
  "address": "Jakarta Selatan, JL bpk dara NO.43"
}
<img width="608" height="523" alt="image" src="https://github.com/user-attachments/assets/778e1572-064c-4190-bfc9-958735389859" />

- get all supplier : GET / http://localhost:8080/suppliers/
  <img width="791" height="801" alt="image" src="https://github.com/user-attachments/assets/281050bf-c816-467c-9ebe-e31eba64f401" />

- get supplier by id : GET / http://localhost:8080/suppliers/:id
  <img width="620" height="519" alt="image" src="https://github.com/user-attachments/assets/8c329466-e50f-40c9-baae-a1263630b124" />

- update supplier : PATCH / http://localhost:8080/suppliers/:id
  <img width="434" height="554" alt="image" src="https://github.com/user-attachments/assets/7f007e1c-5415-4ba7-89f8-1e878fc89f35" />

- deactivate supplier : PATCH / http://localhost:8080/suppliers/:id/deactivate
  <img width="528" height="382" alt="image" src="https://github.com/user-attachments/assets/7b134254-4f58-4ef1-a9b0-11e5c8167a15" />
-------------------------------------------------------------------------

Modul **warehouse** 
- create warehouse : POST / http://localhost:8080/warehouses/
  PAYLOAD : {
  "code": "WH-JKT-002",
  "name": "Warehouse Kopdes",
  "location": "Jakarta Utara NO.18"
}
<img width="398" height="471" alt="image" src="https://github.com/user-attachments/assets/85cf64fe-35a2-4026-a89b-374cd943c1af" />

- get all warehouse : GET /http://localhost:8080/warehouses
  <img width="534" height="735" alt="image" src="https://github.com/user-attachments/assets/78f3b19b-691e-4699-ad9e-ccf0dcacc1f4" />

- get by id : GET / http://localhost:8080/warehouses/:id
  <img width="601" height="471" alt="image" src="https://github.com/user-attachments/assets/8879be90-78b6-45e1-ae47-54329febfd3e" />

- update : PATCH / http://localhost:8080/warehouses/:id
  PAYLOAD : {
  "code": "WH-JKT-001",
  "name": "Warehouse Jakarta Selatan",
  "location": "Jakarta Selatan"
}
<img width="462" height="469" alt="image" src="https://github.com/user-attachments/assets/a93f2c2d-7a1e-4001-ae74-a250bbb5053f" />
-------------------------------------------------------------------------

Modul inventory
- create inventory : POST /http://localhost:8080/inventories/
  <img width="485" height="483" alt="image" src="https://github.com/user-attachments/assets/a550884e-5eee-4e67-9a97-9156ebab3977" />

- get all : GET / http://localhost:8080/inventories
  <img width="561" height="901" alt="image" src="https://github.com/user-attachments/assets/d9837b24-9207-4ee4-b263-eea137c6e4ea" />

- get by id : GET / http://localhost:8080/inventories/:id
  <img width="513" height="517" alt="image" src="https://github.com/user-attachments/assets/33e8437b-dd01-469d-86d1-d1db9205e39e" />

- update : PATCH : http://localhost:8080/inventories/:id
  PAYLOAD : {
  "product_id": 1,
  "warehouse_id": 1,
  "stock": 15
}
<img width="492" height="533" alt="image" src="https://github.com/user-attachments/assets/eddf7f2a-83ff-4fd9-b1df-3c4a0deb7731" />
-------------------------------------------------------------------------

modul **purchase request **
- create PR : POST /
  PAYLOAD : {
  "request_number": "PR-2026-000004",
  "warehouse_id": 1,
  "items": [
    {
      "product_id": 3,
      "quantity": 5
    },
    {
      "product_id": 1,
      "quantity": 10
    }
  ]
  atau jika hanya 1 produk :
  {
  "request_number": "PR-2026-000001",
  "warehouse_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 5
    }
  ]
}
<img width="565" height="886" alt="image" src="https://github.com/user-attachments/assets/2c46fb4c-fca8-46f2-acc1-05a141bf4ac0" />

- get all : GET / http://localhost:8080/purchase-requests/
  <img width="679" height="907" alt="image" src="https://github.com/user-attachments/assets/75a4ec6a-c2c4-470b-b85c-0af7276577b0" />

- get by id : GET / http://localhost:8080/purchase-requests/:id
  <img width="678" height="688" alt="image" src="https://github.com/user-attachments/assets/b2163564-ae82-4472-b8f3-8d8bae35b992" />

- update : PATCH / http://localhost:8080/purchase-requests/:id
  PAYLOAD jumlah ingin berapa produk bisa disesuaikan bisa 1 atau lebih : {
  "warehouse_id": 1,
  "status": "SUBMITTED",
  "items": [
    {
      "product_id": 3,
      "quantity": 15
    },
    {
      "product_id": 1,
      "quantity": 6
    }
  ]
}
<img width="552" height="913" alt="image" src="https://github.com/user-attachments/assets/4f165c3a-425c-4971-976c-15a602c471b6" />

- approve PR (hanya token yang role nya APPROVER yang bisa approve, jika tidak akan muncul error) : PATCH : http://localhost:8080/purchase-requests/:id/approve
  <img width="574" height="844" alt="image" src="https://github.com/user-attachments/assets/cc0854b9-5125-40b5-be7a-7f3cc6603bb0" />

  contoh url dengan filter status tertentu : [http://localhost:8080/purchase-requests/?status=REJECTED ](http://localhost:8080/purchase-requests/?status=APPROVED)

- reject PR : PATCH / http://localhost:8080/purchase-requests/:id/reject
  <img width="641" height="872" alt="image" src="https://github.com/user-attachments/assets/2416c56c-d77a-48d4-9ec6-3b357937164d" />
-------------------------------------------------------------------------

modul **purchase order**
- create PO : POST / http://localhost:8080/purchase-requests/:purchase_request_id/purchase-order
  PAYLOAD : {
  "supplier_id": 1
}
<img width="809" height="847" alt="image" src="https://github.com/user-attachments/assets/38cee35d-72b1-40a4-9324-9401b5a9e391" />

- get all : GET / http://localhost:8080/purchase-orders/
  <img width="745" height="888" alt="image" src="https://github.com/user-attachments/assets/807cb429-bb48-4fda-99e9-60696e471773" />

- get by id : GET / http://localhost:8080/purchase-orders/:id
<img width="621" height="851" alt="image" src="https://github.com/user-attachments/assets/0cc204d5-5d8c-41fd-8d93-fb97c5097b48" />

- update status menjadi ordered : PATCH / http://localhost:8080/purchase-orders/:id/status
  PAYLAOD : {
  "status": "ORDERED"
}
<img width="578" height="866" alt="image" src="https://github.com/user-attachments/assets/97c07b12-742f-4048-9cee-c02f0523d5b4" />

- get by purchase request id : GET / http://localhost:8080/purchase-orders/purchase-request/:purchase_request_id
<img width="674" height="699" alt="image" src="https://github.com/user-attachments/assets/16e86bcb-c472-4213-8b2d-b7efce6ed078" />
-------------------------------------------------------------------------

modul **goods receipt**
- create GR : POST /  http://localhost:8080/purchase-orders/:purchase_order_id/goods-receipts
  PAYLOAD (items bisa 1 atau lebih) : {
  "items": [
    {
      "product_id": 1,
      "received_quantity": 3
    }
  ]
}
<img width="623" height="530" alt="image" src="https://github.com/user-attachments/assets/5b8b437e-f7fa-4257-a6e6-0af76f929157" />

- get all : GET / http://localhost:8080/goods-receipts/
  <img width="523" height="903" alt="image" src="https://github.com/user-attachments/assets/79c28265-9b25-4c64-9e0b-76e60ff4d73d" />

- get by id : GET / http://localhost:8080/goods-receipts/:id
  <img width="574" height="620" alt="image" src="https://github.com/user-attachments/assets/378c0c2e-6c25-401c-9b9d-8ad208741b85" />
-------------------------------------------------------------------------

modul **inventory movement**
get all : GET / http://localhost:8080/inventory-movements/
<img width="837" height="909" alt="image" src="https://github.com/user-attachments/assets/6d236e95-34ad-4237-a7ba-da38795717e9" />

parameter filter yang tesedia : product_id, warehouse_id, reference
contoh : http://localhost:8080/inventory-movements/?reference=GR-2026-1789445397350649200
<img width="777" height="569" alt="image" src="https://github.com/user-attachments/assets/b778b7da-6017-4c29-a39c-9cf3f814d80c" />
-------------------------------------------------------------------------

##Automated test :
```text
masukan command berikut di terminal vscode

purchase request :
- Test create purchase request tanpa item 
go test -v -count=1 -run="TestCreate_WithoutItems" ./internal/usecase/purchase_request

- Test Approve purchase request yang masih belum ter submit
go test -v -count=1 -run="TestUpdate_NotSubmitted" ./internal/usecase/purchase_request


purchase order :
- Test create order yang purchase request nya belum approve
go test -v -count=1 -run="TestCreatePR_NotApprove" ./internal/usecase/purchase_order

- Test create order yang purchase request nya sudah memiliki purchase order (duplikat)
go test -v -count=1 -run="AlreadyExists" ./internal/usecase/purchase_order


Goods Receipt :
- test jika produk yang diterima tidak lebih dari 0 (<0)
go test -v -count=1 ./internal/usecase/goods_receipt

- goods receipt menambahkan stock produk 
go test -v -count=1 ./internal/usecase/goods_receipt -run TestCreate_IncreasesInventoryStockByReceivedQuantity

- create saat purchase order sudah terpenuhi ( status = RECEIVED)
go test -v -count=1 -run="TestCreate_WhenPurchaseOrderAlreadyReceived_ShouldReturnError" ./internal/usecase/goods_receipt
```
<img width="931" height="274" alt="image" src="https://github.com/user-attachments/assets/55110d7e-3a57-4651-8c33-cee010f20bd4" />





















