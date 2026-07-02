package models

import "gorm.io/gorm"

// User adalah model yang mapping ke tabel "users" di MySQL
// GORM otomatis plural nama struct: User -> users
type User struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)

	FirebaseUID  string `gorm:"uniqueIndex;size:128;not null" json:"firebase_uid"`
	Email        string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name         string `gorm:"size:100" json:"name"`
	Role         string `gorm:"size:20;default:user" json:"role"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	LastLoginAt  *int64 `gorm:"index" json:"last_login_at,omitempty"`
}

/*
Penjelasan:

gorm.Model memberikan fields:
- ID uint (primary key, auto increment)
- CreatedAt time.Time
- UpdatedAt time.Time
- DeletedAt gorm.DeletedAt (soft delete)

Struct tag "gorm":
- uniqueIndex = buat unique index
- size:128 = varchar(128)
- not null = tidak boleh NULL
- default:user = nilai default

Struct tag "json":
- json:"firebase_uid" = nama key di response JSON
- json:"-" = tidak ditampilkan
- omitempty = tidak dikirim kalau nil/kosong
*/