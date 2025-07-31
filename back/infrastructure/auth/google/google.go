package google

import (
	"golang.org/x/oauth2"
)

var (
	GoogleEndpoint = oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	}
)
