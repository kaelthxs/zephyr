package usecase

import (
    "errors"
    "fmt"
    "math/rand"
    "time"

    "zephyr-backend/infrastructure/cache"
    "zephyr-backend/infrastructure/sms"
    "zephyr-backend/internal/repository"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

type Mailer interface {
    SendVerificationEmail(email, link string) error
}

type UserUseCase struct {
    repo  repository.UserRepository
    auth  repository.AuthService
    cache *cache.RedisClient
    sms   *sms.SmsClient
    mailer Mailer
    baseURL string
}

func NewUserUseCase(
    repo repository.UserRepository,
    auth repository.AuthService,
    cache *cache.RedisClient,
    sms *sms.SmsClient,
    mailer Mailer,
    baseURL string,
) *UserUseCase {
    return &UserUseCase{
        repo: repo,
        auth: auth,
        cache: cache,
        sms: sms,
        mailer: mailer,
        baseURL: baseURL,
    }
}


func generateRandomCode() string {
    return uuid.NewString()
}

func (uc *UserUseCase) Register(username, email, password, birth_date, phone_number string) error {
    hashed, err := uc.auth.HashPassword(password)
    if err != nil {
        return err
    }
    return uc.repo.CreateUser(
        username,
        email,
        hashed,
        birth_date,
        phone_number,
        "",
        "",
        "мужской",
        "",
        "local",
    )
}

func (uc *UserUseCase) Login(email, password string) (string, error) {
    user, err := uc.repo.GetByEmail(email)
    if err != nil {
        return "", err
    }
    if !uc.auth.CheckPassword(password, user.Password) {
        return "", errors.New("invalid credentials")
    }
    return uc.auth.GenerateToken(user.ID.String())
}


func (uc *UserUseCase) SendPhoneCode(phone string) (string, error) {
    code := fmt.Sprintf("%04d", rand.Intn(10000))
    if err := uc.cache.Set("phone_code:"+phone, code, 5*time.Minute); err != nil {
        return "", err
    }
    if err := uc.sms.SendSms(phone, "Ваш код подтверждения: "+code); err != nil {
        return "", err
    }
    return code, nil
}


func (uc *UserUseCase) ConfirmPhone(phone, code string) error {
    savedCode, err := uc.cache.Get("phone_code:" + phone)
    if err != nil {
        return err
    }
    if savedCode != code {
        return errors.New("wrong code")
    }
    if err := uc.repo.SetPhoneVerified(phone); err != nil {
        return err
    }
    return uc.cache.Delete("phone_code:" + phone)
}

func (uc *UserUseCase) SendEmailVerificationCode(email string) error {
    code := generateRandomCode()
    if err := uc.cache.SetEmailCode(email, code, time.Hour); err != nil {
        return err
    }
    link := fmt.Sprintf("%s/verify-email?code=%s", uc.baseURL, code)
    return uc.mailer.SendVerificationEmail(email, link)
}


func (uc *UserUseCase) ConfirmEmail(code string) error {
    email, err := uc.cache.GetEmailByCode(code)
    if err != nil {
        return err
    }
    if err := uc.repo.SetEmailVerified(email); err != nil {
        return err
    }
    _ = uc.cache.DeleteEmailCode(email)
    return nil
}

func (uc *UserUseCase) LoginOrRegisterWithYandex(
	email, login, firstName, lastName, birthday, gender, yandexID string,
) (string, error) {
    user, err := uc.repo.GetByEmail(email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            user = nil
        } else {
            return "", err
        }
    }
    if user == nil {
        if birthday == "" {
            birthday = "1900-01-01"
        }
        if err := uc.repo.CreateUser(
            login,
            email,
            "oauth_yandex_placeholder",
            birthday,
            "0000000000",
            firstName,
            lastName,
            gender,
            yandexID,
            "yandex",
        ); err != nil {
            return "", err
        }
        user, err = uc.repo.GetByEmail(email)
        if err != nil {
            return "", err
        }
    }
    return uc.auth.GenerateToken(user.ID.String())
}





