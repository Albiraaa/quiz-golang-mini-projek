# Golang Book API

Project ini adalah REST API sederhana untuk manajemen **buku** dan **kategori** menggunakan:

API ini dibuat sesuai ketentuan brief:
- CRUD **Category**
- CRUD **Book**
- Relasi Category–Book
- Validasi `release_year` 1980–2024
- Hitung otomatis **thickness** (`tipis` / `tebal`) dari `total_page`
- JWT untuk proteksi endpoint

---

## 1. Tech Stack

- Go 1.22+  
- Gin  
- PostgreSQL  
- lib/pq (Postgres driver)  
- sql-migrate  
- golang-jwt/jwt v5  

---


## 2. CRUD Postman

Menggunakan postman untuk melakukan CRUD 

Sistem sudah dideploy di railway

(https://quiz-golang-mini-projek-production.up.railway.app/api/buku)
catatan penting membutuhkan token Bearer Authorization untuk mengakses
Body Post Buku
    {
  "title": "",
  "description": "",
  "image_url": "",
  "release_year": ,
  "price": ,
  "total_page": ,
  "category_id": 
}

(https://quiz-golang-mini-projek-production.up.railway.app/api/kategori)
catatan penting membutuhkan token Bearer Authorization untuk mengakses
Body Post Kategori
{
  "name": ""
}

(https://quiz-golang-mini-projek-production.up.railway.app/api/users/login)
Ketika sudah login maka akan diberikan token
jika menggunakan postman harap di inputkan di menu Authorization (Bearer Token)
jika di web maka harus diinputkan manual melalui js
(https://quiz-golang-mini-projek-production.up.railway.app/api/users/register)
Body Post register and login
{
  "username": "",
  "password": ""
}

Untuk Get sama semua (membutuhkan token juga)
https://quiz-golang-mini-projek-production.up.railway.app/api/kategori
https://quiz-golang-mini-projek-production.up.railway.app/api/buku

===========================

## 3. Cara melalui WEB
bisa buka file test.html yang sudah saya buat isinya adalah html sederhana dengan sedikit js

untuk base url api bisa dibiarkan otomatis atau jika tidak ada bisa menambahkan:
https://quiz-golang-mini-projek-production.up.railway.app

lanjut ke form login masukin aja ini
admin   :user
123456  :password

token akan otomatis tersimpan 

jika ingin register bisa melakukannya di postman ada di penjelasan sebelumnya

selanjutnya form get untuk melihat semua buku dan kategori
serta get from ID 

## 4. Struktur Folder

Kurang lebih seperti ini:

```bash
quiz-golang-books/
├── go.mod
├── main.go
├── config/
│   └── db.go
├── models/
│   ├── book.go
│   ├── category.go
│   └── user.go
├── middleware/
│   └── auth.go
├── handlers/
│   ├── auth_handler.go
│   ├── category_handler.go
│   └── book_handler.go
└── migrations/
    └── 001_init.sql
└── test.html 

