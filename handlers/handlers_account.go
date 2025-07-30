package handlers

import (
	"blog-app/middlewares"
	"blog-app/models"
	"blog-app/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UpdateAccountReq struct {
	Name      *string `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	var Accounts []models.Account

	if err := DB.Preload("Staff.Account").Preload("Client.Account").Find(&Accounts).Error; err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Accounts)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Login failed!")
		return
	}
	var acc models.Account
	if err := DB.Where("email = ? ", req.Email).First(&acc).Error; err != nil {
		http.Error(w, "Wrong email!", http.StatusBadRequest)
		return
	}

	if !CheckPasswordHash(req.Password, acc.Password) {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
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

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"refresh_token": refresh_token,
		"access_token":  access_token,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	var account models.Account
	if err := DB.Where("email = ?", req.Email).First(&account).Error; err == nil {
		http.Error(w, "Email is exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Hash failed", http.StatusInternalServerError)
		return
	}

	newAccount := models.Account{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	if err := DB.Create(&newAccount).Error; err != nil {
		http.Error(w, "Register failed (account)", http.StatusInternalServerError)
		return
	}

	newClient := models.Client{
		AccountID: newAccount.ID,
	}

	if err := DB.Create(&newClient).Error; err != nil {
		http.Error(w, "Register failed", http.StatusInternalServerError)
		return
	}

	var createdAccount models.Account
	if err := DB.Preload("Client").First(&createdAccount, newAccount.ID).Error; err != nil {
		http.Error(w, "Reload failed", http.StatusInternalServerError)
		return
	}

	access_token, err := middlewares.GenerateAccessToken(createdAccount.ID, createdAccount.Role, createdAccount.Email, createdAccount.AvatarURL, createdAccount.Name)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refresh_token, err := middlewares.GenerateRefreshToken(createdAccount.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	decoded := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(refresh_token, decoded)
	if err != nil {
		http.Error(w, "Failed to parse token", http.StatusInternalServerError)
		return
	}
	exp := time.Unix(int64(decoded["exp"].(float64)), 0)

	tokenRecord := models.BlackListToken{
		Token:     refresh_token,
		Type:      "refresh",
		AccountID: createdAccount.ID,
		ExpiresAt: exp,
	}
	if err := DB.Create(&tokenRecord).Error; err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"refresh_token": refresh_token,
		"access_token":  access_token,
	})
}

func GetAccountById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var acc models.Account
	if err := DB.Find(&acc, id).Error; err != nil {
		http.Error(w, "This account is not exists!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func GetAccountByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	var acc models.Account
	if err := DB.Where("email = ?", email).First(&acc).Error; err != nil {
		http.Error(w, "This account does not exist!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Thiếu Bearer token", http.StatusUnauthorized)
		return
	}
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := middlewares.ParseToken(refreshToken)
	if err != nil {
		http.Error(w, "Token không hợp lệ", http.StatusUnauthorized)
		return
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.After(time.Now()) {
		DB.Create(&models.BlackListToken{
			Token:     refreshToken,
			AccountID: claims.AccountID,
			ExpiresAt: claims.ExpiresAt.Time,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout success!",
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Can't find tolen", http.StatusUnauthorized)
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := middlewares.ParseToken(refreshToken)
	if err != nil || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Token is invalid or expire", http.StatusUnauthorized)
		return
	}

	var blacklisted models.BlackListToken
	if err := DB.Where("token = ?", refreshToken).First(&blacklisted).Error; err == nil {
		http.Error(w, "Token recall!", http.StatusUnauthorized)
		return
	}

	var acc models.Account

	if err := DB.Find(&acc, claims.AccountID).Error; err != nil {
		http.Error(w, "This account is not exists!", http.StatusBadRequest)
		return
	}

	newToken, err := middlewares.GenerateAccessToken(acc.ID, acc.Role, acc.Email, acc.AvatarURL, acc.Name)
	if err != nil {
		http.Error(w, "Không thể tạo token mới", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": newToken,
	})
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var account models.Account

	if err := DB.First(&account, id).Error; err != nil {
		http.Error(w, "This account does not exists!", http.StatusBadRequest)
		return
	}

	var req UpdateAccountReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusBadGateway)
		return
	}

	if req.Name != nil {
		account.Name = *req.Name
	}

	if req.AvatarURL != nil {
		account.AvatarURL = *req.AvatarURL
	}

	if err := DB.Save(&account).Error; err != nil {
		http.Error(w, "Failed to update Account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
