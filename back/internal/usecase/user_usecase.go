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

// init seeds the random number generator used for generating verification codes.
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
    // hash the provided plain text password. If hashing fails we return the error
    hashed, err := uc.auth.HashPassword(password)
    if err != nil {
        return err
    }
    // persist a new user with default values for optional fields
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
    // verify the provided password against the stored hash
    if !uc.auth.CheckPassword(password, user.Password) {
        return "", errors.New("invalid credentials")
    }
    // generate a JWT token for the authenticated user
    return uc.auth.GenerateToken(user.ID.String())
}


func (uc *UserUseCase) SendPhoneCode(phone string) (string, error) {
    // generate a 4‑digit numeric code
    code := fmt.Sprintf("%04d", rand.Intn(10000))
    // store the code in cache with a 5 minute expiry
    if err := uc.cache.Set("phone_code:"+phone, code, 5*time.Minute); err != nil {
        return "", err
    }
    // send the code via the SMS client
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
    // mark the phone as verified in the persistent store
    if err := uc.repo.SetPhoneVerified(phone); err != nil {
        return err
    }
    // remove the code from cache
    return uc.cache.Delete("phone_code:" + phone)
}

func (uc *UserUseCase) SendEmailVerificationCode(email string) error {
    code := generateRandomCode()
    // associate the email with the generated code in cache
    if err := uc.cache.SetEmailCode(email, code, time.Hour); err != nil {
        return err
    }
    // build a verification link using the base URL
    link := fmt.Sprintf("%s/verify-email?code=%s", uc.baseURL, code)
    // send the email via the mailer implementation
    return uc.mailer.SendVerificationEmail(email, link)
}


func (uc *UserUseCase) ConfirmEmail(code string) error {
    email, err := uc.cache.GetEmailByCode(code)
    if err != nil {
        return err
    }
    // mark the email as verified in the persistent store
    if err := uc.repo.SetEmailVerified(email); err != nil {
        return err
    }
    // remove both forward and reverse mappings from cache; ignore errors
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
    // if the user does not exist, register a new account using the
    // information returned from Yandex. Some defaults are applied
    // when data is missing.
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
        // reload the user after creation
        user, err = uc.repo.GetByEmail(email)
        if err != nil {
            return "", err
        }
    }
    // generate a JWT token for the (existing or newly created) user
    return uc.auth.GenerateToken(user.ID.String())
}





