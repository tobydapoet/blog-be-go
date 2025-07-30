package handlers

import (
	"blog-app/middlewares"
	"blog-app/models"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = os.Getenv("CORS_ORIGIN")
	}
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"picture"`
		GoogleID  string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	var acc models.Account

	if err := DB.Where("email = ?", userInfo.Email).First(&acc).Error; err != nil {
		acc = models.Account{
			Email:     userInfo.Email,
			Password:  "",
			Role:      "client",
			Name:      userInfo.Name,
			IsDeleted: false,
			GoogleID:  userInfo.GoogleID,
			AvatarURL: userInfo.AvatarURL,
		}

		if err := DB.Create(&acc).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		client := models.Client{
			AccountID: acc.ID,
		}

		if err := DB.Create(&client).Error; err != nil {
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

	}

	access_token, err := middlewares.GenerateAccessToken(acc.ID, acc.Role, acc.Email, acc.AvatarURL, acc.Name)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refresh_token, err := middlewares.GenerateRefreshToken(acc.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	redirectURL := state
	if redirectURL == "" {
		redirectURL = os.Getenv("CORS_ORIGIN")
	}

	http.Redirect(w, r, redirectURL+"?access_token="+access_token+"&"+"refresh_token="+refresh_token, http.StatusFound)
}
