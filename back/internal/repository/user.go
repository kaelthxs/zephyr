package repository

import "zephyr-backend/internal/domain"

// UserRepository defines the required persistence operations for a user.
// Implementations live in the infrastructure layer and should be free
// of any framework code. Use cases depend on this interface to remain
// decoupled from storage details.
type UserRepository interface {
    // CreateUser persists a new user in the storage. Fields such as
    // username, email, hashed password and other optional metadata are
    // provided individually to keep the entity ignorant of external
    // concerns (e.g. password hashing happens in the service layer).
    CreateUser(
        username, email, passwordHash, birthDate, phoneNumber,
        firstName, lastName, gender, yandexID, oauthProvider string,
    ) error

    // GetByEmail retrieves a user by email. Returns gorm.ErrRecordNotFound
    // if no user exists with the given email.
    GetByEmail(email string) (*domain.User, error)

    // GetByPhone retrieves a user by phone number. Returns
    // gorm.ErrRecordNotFound if no user exists with the given phone.
    GetByPhone(phone string) (*domain.User, error)

    // SetPhoneVerified updates a user record to mark the phone number as
    // verified. Returns an error on failure.
    SetPhoneVerified(phone string) error

    // SetEmailVerified updates a user record to mark the email as verified.
    SetEmailVerified(email string) error
}
