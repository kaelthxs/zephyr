package usecase

import (
	"errors" // и этого
	"fmt"
	"log"
	"math/rand" // вот этого не хватает
	"time"
	"zephyr-backend/infrastructure/cache"
	"zephyr-backend/infrastructure/sms"
	"zephyr-backend/internal/port/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
    hashed, _ := uc.auth.HashPassword(password)
    return uc.repo.CreateUser(
        username,
        email,
        hashed,
        birth_date,
        phone_number,
        "", // firstName
        "", // lastName
        "мужской", // gender
        "", // yandexID
        "local", // oauthProvider
    )
    
}

func (uc *UserUseCase) Login(email, password string) (string, error) {
    fmt.Println("⚡️ Вызвался usecase.Login")

    user, err := uc.repo.GetByEmail(email)
    if err != nil || !uc.auth.CheckPassword(password, user.Password) {
        fmt.Println("⛔️ Ошибка в логине:", err)
        return "", err
    }

    token, err := uc.auth.GenerateToken(user.ID.String())
    fmt.Println("🎫 Токен из auth:", token)

    return token, err
}


func (uc *UserUseCase) SendPhoneCode(phone string) (string, error) {
    code := fmt.Sprintf("%04d", rand.Intn(10000))
    err := uc.cache.Set("phone_code:"+phone, code, 5*time.Minute)
    if err != nil {
        return "", err
    }

    err = uc.sms.SendSms(phone, "Ваш код подтверждения: "+code)
    if err != nil {
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
    err = uc.repo.SetPhoneVerified(phone)
    if err != nil {
        return err
    }
    return uc.cache.Delete("phone_code:" + phone)
}

func (uc *UserUseCase) SendEmailVerificationCode(email string) error {
    code := generateRandomCode()
    err := uc.cache.SetEmailCode(email, code, time.Hour)
    if err != nil {
        log.Printf("Redis set error: %v", err)
        return err
    }

    link := fmt.Sprintf("%s/verify-email?code=%s", uc.baseURL, code)

    err = uc.mailer.SendVerificationEmail(email, link)
    if err != nil {
        log.Printf("SMTP send error: %v", err)
        return err
    }

    return nil
}


func (uc *UserUseCase) ConfirmEmail(code string) error {
    email, err := uc.cache.GetEmailByCode(code)
    if err != nil {
        return err
    }

    err = uc.repo.SetEmailVerified(email)
    if err != nil {
        return err
    }

    if err := uc.cache.DeleteEmailCode(email); err != nil {
        log.Printf("Warning: failed to delete email code from cache for %s: %v", email, err)
    }

    return nil
}

func (uc *UserUseCase) LoginOrRegisterWithYandex(
	email, login, firstName, lastName, birthday, gender, yandexID string,
) (string, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user = nil
	} else if err != nil {
		return "", err
	}

	// если пользователь не найден — регистрируем
	if user == nil {
		if birthday == "" {
			birthday = "1900-01-01"
		}

		err = uc.repo.CreateUser(
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
		)
		if err != nil {
			return "", err
		}

		user, _ = uc.repo.GetByEmail(email)
	}

	// токен
	token, err := uc.auth.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}





