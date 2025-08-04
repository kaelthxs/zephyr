package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DBUser             string
	DBPassword         string
	DBName             string
	DBHost             string
	DBPort             string
	Port               string
	JWTSecret          string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	DBUrl              string
	AppBaseURL         string
	YandexClientID     string
	YandexClientSecret string
	YandexRedirectURI  string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	redisDBStr := os.Getenv("redis_db")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		Port:       os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),

		RedisAddr:     os.Getenv("redis_addr"),
		RedisPassword: os.Getenv("redis_password"),
		RedisDB:       redisDB,
		DBUrl:         os.Getenv("DBURL"),

		AppBaseURL: os.Getenv("APP_BASE_URL"),

		YandexClientID:     os.Getenv("YANDEX_CLIENT_ID"),
		YandexClientSecret: os.Getenv("YANDEX_CLIENT_SECRET"),
		YandexRedirectURI:  os.Getenv("YANDEX_REDIRECT_URI"),

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURI:  os.Getenv("GOOGLE_REDIRECT_URI"),
	}, nil
}
