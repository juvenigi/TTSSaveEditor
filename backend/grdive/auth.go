package grdive

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

var AuthConfig oauth2.Config
var AuthToken oauth2.Token

func init() {
	var ClientId = os.Getenv("TTS_GDIRVE_CLIENT")
	var ClientSecret = os.Getenv("TTS_GDRIVE_SECRET")
	if ClientId == "" || ClientSecret == "" {
		panic("envvars not set!")
	}

	AuthConfig = oauth2.Config{
		ClientID:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  "http://localhost:3000/callback",
		Scopes:       []string{drive.DriveScope},
		Endpoint:     google.Endpoint,
	}
}

func SetTokenFromCallback(code string) error {
	myToken, err := AuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	AuthToken = *myToken
	return nil
}

func GetDriveService(ctx context.Context) (*drive.Service, error) {
	var src = AuthConfig.TokenSource(ctx, &AuthToken)
	return drive.NewService(ctx, option.WithTokenSource(src))
}
