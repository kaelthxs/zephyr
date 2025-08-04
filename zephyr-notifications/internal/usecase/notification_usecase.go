package usecase

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Mailer interface {
	SendVerificationEmail(email, link string) error
}

func (uc *NotificationUseCase) SendPhoneCode(phone string) (string, error) {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	if err := uc.cache.Set("phone_code:"+phone, code, 5*time.Minute); err != nil {
		return "", err
	}
	if err := uc.sms.SendSms(phone, "Ваш код подтверждения: "+code); err != nil {
		return "", err
	}
	return code, nil
}

func (uc *NotificationUseCase) ConfirmPhone(phone, code string) error {
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

func (uc *NotificationUseCase) SendEmailVerificationCode(email string) error {
	code := generateRandomCode()
	if err := uc.cache.SetEmailCode(email, code, time.Hour); err != nil {
		return err
	}
	link := fmt.Sprintf("%s/verify-email?code=%s", uc.baseURL, code)
	return uc.mailer.SendVerificationEmail(email, link)
}

func (uc *NotificationUseCase) ConfirmEmail(code string) error {
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
