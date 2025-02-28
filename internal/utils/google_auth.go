package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// Variabel global untuk menyimpan konfigurasi OAuth2 Google
var GoogleOauthConfig *oauth2.Config

// ✅ Fungsi eksplisit untuk inisialisasi OAuth2 (panggil ini di main.go)
func SetupGoogleOAuth() {
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     env.GetString("GOOGLE_CLIENT_ID", ""),
		ClientSecret: env.GetString("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  env.GetString("GOOGLE_REDIRECT_URI", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	fmt.Println("✅ Google OAuth Config Initialized")
}

// 🔄 Tukar authorization code dengan token akses
func ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	if GoogleOauthConfig == nil {
		return nil, fmt.Errorf("GoogleOauthConfig is not initialized. Call SetupGoogleOAuth() first")
	}

	ctx := context.Background()
	token, err := GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		fmt.Println("Error exchanging code:", err)
		return nil, err
	}
	return token, nil
}

// 🔎 Ambil informasi user dari Google API
func FetchGoogleUserInfo(accessToken string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
