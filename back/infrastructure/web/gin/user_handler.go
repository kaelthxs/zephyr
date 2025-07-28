package user_handler

import (
	"log"
	"net/http"
	"zephyr-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, uc *usecase.UserUseCase, authMiddleware gin.HandlerFunc) {
    r.POST("/register", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password_hash"`
            Birth_date string `json:"birth_date"`
            PhoneNumber string `json:"phone_number"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        err := uc.Register(req.Username, req.Email, req.Password, req.Birth_date, req.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "registered"})
    })

    r.POST("/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password_hash"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }

        token, err := uc.Login(req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": token})
    })

    authGroup := r.Group("/auth", authMiddleware)
    authGroup.GET("/profile", func(c *gin.Context) {
        userID := c.GetString("user_id")
        c.JSON(http.StatusOK, gin.H{
            "message": "Ты залогинился, брат",
            "user_id": userID,
        })
    })

    r.POST("/send-code", func(c *gin.Context) {
        var req struct {
            Phone string `json:"phone"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }
        code, err := uc.SendPhoneCode(req.Phone)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send code"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "code sent", "code": code})
        log.Printf("Код для подтверждения телефона: " + code + " Для номера ", req.Phone)
    })
    
    r.POST("/confirm-phone", func(c *gin.Context) {
        var req struct {
            Phone string `json:"phone"`
            Code  string `json:"code"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }
        err := uc.ConfirmPhone(req.Phone, req.Code)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "phone confirmed"})
    })

    r.POST("/send-email-code", func(c *gin.Context) {
        var req struct {
            Email string `json:"email"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }
        err := uc.SendEmailVerificationCode(req.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email verification code"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "verification email sent"})
    })

    r.GET("/verify-email", func(c *gin.Context) {
        code := c.Query("code")
        if code == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
            return
        }
        err := uc.ConfirmEmail(code)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired code"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "email confirmed"})
    })
    
}