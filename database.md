### **Entitas Tabel di Database (PostgreSQL)**

#### 1. **Tabel `users`**
- `id`: UUID (Primary Key)
- `username`: String (Unique)
- `email`: String (Unique)
- `password_hash`: String
- `role`: Enum (admin, editor, penulis)
- `created_at`: Timestamp
- `updated_at`: Timestamp
- `2fa_enabled`: Boolean
- `2fa_secret`: String

#### 2. **Tabel `articles`**
- `id`: UUID (Primary Key)
- `title`: String
- `content`: Text
- `slug`: String (Unique)
- `status`: Enum (draft, review, published)
- `author_id`: UUID (Foreign Key ke `users`)
- `published_at`: Timestamp (Nullable)
- `created_at`: Timestamp
- `updated_at`: Timestamp

#### 3. **Tabel `article_versions`** (Untuk Versioning)
- `id`: UUID (Primary Key)
- `article_id`: UUID (Foreign Key ke `articles`)
- `title`: String
- `content`: Text
- `created_at`: Timestamp

#### 4. **Tabel `categories`**
- `id`: UUID (Primary Key)
- `name`: String (Unique)
- `description`: Text
- `created_at`: Timestamp
- `updated_at`: Timestamp

#### 5. **Tabel `tags`**
- `id`: UUID (Primary Key)
- `name`: String (Unique)
- `created_at`: Timestamp
- `updated_at`: Timestamp

#### 6. **Tabel `article_tags`** (Many-to-Many antara `articles` dan `tags`) (auto-generate using gorm)
- `article_id`: UUID (Foreign Key ke `articles`)
- `tag_id`: UUID (Foreign Key ke `tags`)

#### 6. **Tabel `article_categories`** (Many-to-Many antara `articles` dan `categories`) (auto-generate using gorm)
- `article_id`: UUID (Foreign Key ke `articles`)
- `category_id`: UUID (Foreign Key ke `tags`)

#### 7. **Tabel `images`**
- `id`: UUID (Primary Key)
- `url`: String
- `alt_text`: String
- `article_id`: UUID (Foreign Key ke `articles`)
- `caption`: String
- `uploaded_by`: UUID (Foreign Key ke `users`)
- `created_at`: Timestamp

#### 8. **Tabel `seo_metadata`**
- `id`: UUID (Primary Key)
- `article_id`: UUID (Foreign Key ke `articles`)
- `meta_title`: String
- `meta_description`: Text
- `keywords`: Text
- `created_at`: Timestamp
- `updated_at`: Timestamp

#### 9. **Tabel `comments`**
- `id`: UUID (Primary Key)
- `article_id`: UUID (Foreign Key ke `articles`)
- `user_id`: UUID (Foreign Key ke `users`)
- `content`: Text
- `created_at`: Timestamp
- `updated_at`: Timestamp

#### 10. **Tabel `analytics`**
- `id`: UUID (Primary Key)
- `article_id`: UUID (Foreign Key ke `articles`)
- `views`: Integer
- `unique_visitors`: Integer
- `average_time_spent`: Integer
- `created_at`: Timestamp

#### 11. **Tabel `backups`**
- `id`: UUID (Primary Key)
- `file_path`: String
- `created_by`: UUID (Foreign Key ke `users`)
- `created_at`: Timestamp
