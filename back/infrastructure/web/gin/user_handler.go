package handler

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"

    "zephyr-backend/config"
    "zephyr-backend/internal/usecase"
    yandexOAuth "zephyr-backend/infrastructure/auth/yandex"
)

type AuthHandler struct {
    UC *usecase.UserUseCase
}

func NewAuthHandler(uc *usecase.UserUseCase) *AuthHandler {
    return &AuthHandler{UC: uc}
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var req struct {
        UserID       string `json:"user_id"`
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    token, refresh, err := h.UC.RefreshToken(req.UserID, req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  token,
        "refresh_token": refresh,
    })
}

func (h *AuthHandler) Logout(c *gin.Context) {
    var req struct {
        UserID       string `json:"user_id"`
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    err := h.UC.Logout(req.UserID, req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func RegisterRoutes(r *gin.Engine, uc *usecase.UserUseCase, authMiddleware gin.HandlerFunc, cfg *config.Config) {
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

        accessToken, refreshToken, err := uc.Login(req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "access_token": accessToken,
            "refresh_token": refreshToken,
        })
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

    oauthConfig := &oauth2.Config{
        ClientID:     cfg.YandexClientID,
        ClientSecret: cfg.YandexClientSecret,
        RedirectURL:  cfg.YandexRedirectURI,
        Scopes:       []string{"login:email", "login:info"},
        Endpoint:     yandexOAuth.YandexEndpoint,
    }

    r.GET("/auth/yandex/login", func(c *gin.Context) {
        url := oauthConfig.AuthCodeURL("state-zephyr")
        c.Redirect(http.StatusTemporaryRedirect, url)
    })

    r.GET("/auth/yandex/callback", func(c *gin.Context) {
        code := c.Query("code")
        if code == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
            return
        }
        token, err := oauthConfig.Exchange(context.Background(), code)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
            return
        }
        client := oauthConfig.Client(context.Background(), token)
        resp, err := client.Get("https://login.yandex.ru/info?format=json")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
            return
        }
        defer resp.Body.Close()
        var userInfo struct {
            Email     string `json:"default_email"`
            Login     string `json:"login"`
            FirstName string `json:"first_name"`
            LastName  string `json:"last_name"`
            Birthday  string `json:"birthday"`
            SexRaw    string `json:"sex"`
            ID        string `json:"id"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "decode failed"})
            return
        }
        gender := "не выбран"
        switch userInfo.SexRaw {
        case "1":
            gender = "мужской"
        case "2":
            gender = "женский"
        }
        tokenStr, err := uc.LoginOrRegisterWithYandex(
            userInfo.Email,
            userInfo.Login,
            userInfo.FirstName,
            userInfo.LastName,
            userInfo.Birthday,
            gender,
            userInfo.ID,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "login/register failed"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": tokenStr})
    })

    authHandler := NewAuthHandler(uc)
    r.POST("/auth/refresh", authHandler.RefreshToken)
    r.POST("/auth/logout", authHandler.Logout)
}