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

## 2. Struktur Folder

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

