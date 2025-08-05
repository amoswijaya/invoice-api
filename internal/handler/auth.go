package handler

import (
	"invoice-api/helper"
	"invoice-api/internal/db"
	"invoice-api/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
    FullName string `json:"fullName" binding:"required"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }

    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    user := model.User{
        FullName: req.FullName,
        Email:    req.Email,
        Password: string(hash),
    }

    if err := db.DB.Create(&user).Error; err != nil {
        c.JSON(409, gin.H{"error": "email already exists"}); return
    }

    c.JSON(201, gin.H{"message": "user created"})
}

type loginReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var req loginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }

    var user model.User
    if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"error": "invalid credentials"}); return
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
        c.JSON(401, gin.H{"error": "invalid credentials"}); return
    }

    token, _ := helper.GenerateToken(user.ID)
    c.JSON(200, gin.H{
        "message": "login success",
        "token":   token,
    })
}