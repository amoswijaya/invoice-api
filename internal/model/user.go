package model

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    FullName string `gorm:"not null" json:"fullName"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `gorm:"not null" json:"-"`
}