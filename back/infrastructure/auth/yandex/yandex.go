package yandex

import (
	"golang.org/x/oauth2"
)

var (
	YandexEndpoint = oauth2.Endpoint{
		AuthURL:  "https://oauth.yandex.ru/authorize",
		TokenURL: "https://oauth.yandex.ru/token",
	}
)
