package models

import "gorm.io/gorm"

type Role string

const (
	RoleOwner  Role = "Owner"
	RoleClient Role = "Client"
	RoleAdmin  Role = "Admin"
)

type User struct {
	gorm.Model
	FullName string `json:"full_name"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-"`
	Role     Role   `json:"role" gorm:"type:varchar(20);default:'client'"`
}
