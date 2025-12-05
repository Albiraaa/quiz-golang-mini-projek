-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username     VARCHAR(100) NOT NULL UNIQUE,
    password     VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    created_by   VARCHAR(100),
    modified_at  TIMESTAMP NOT NULL,
    modified_by  VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    created_by  VARCHAR(100),
    modified_at TIMESTAMP NOT NULL,
    modified_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  VARCHAR(500) NOT NULL,
    image_url    VARCHAR(255),
    release_year INTEGER NOT NULL,
    price        INTEGER NOT NULL,
    total_page   INTEGER NOT NULL,
    thickness    VARCHAR(20) NOT NULL,
    category_id  INTEGER NOT NULL REFERENCES categories(id),
    created_at   TIMESTAMP NOT NULL,
    created_by   VARCHAR(100),
    modified_at  TIMESTAMP NOT NULL,
    modified_by  VARCHAR(100)
);

-- +migrate Down
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;