package handlers

import (
	"blog-app/models"
	"blog-app/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateClientRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Bio      string `json:"description"`
}

type UpdateClientRequest struct {
	Bio           *string `json:"description"`
	LinkInstagram *string `json:"link_instagram"`
	LinkFacebook  *string `json:"link_facebook"`
	LinkWebsite   *string `json:"link_website"`
}

func GetAllClients(w http.ResponseWriter, r *http.Request) {
	var clientList []models.Client
	DB.Find(&clientList)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientList)
}

func GetClientById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var client models.Client

	if err := DB.Find(&client, id).Error; err != nil {
		http.Error(w, "This client is not exists!", http.StatusBadRequest)
		return
	}
}

func GetClientByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	var client models.Client

	err := DB.
		Preload("Account").
		Joins("JOIN accounts ON accounts.id = clients.account_id").
		Where("accounts.email = ?", email).
		First(&client).Error

	if err != nil {
		http.Error(w, "Client not found", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func GetClientByAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	var client models.Client

	if err := DB.Where(`account_id = ?`, accountID).First(&client).Error; err != nil {
		http.Error(w, "Error when get client!", http.StatusBadRequest)
		return
	}

	if err := DB.Preload(`Account`).First(&client).Error; err != nil {
		http.Error(w, "Error when get relation!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

func SearchClients(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	var clientList []models.Client

	if err := DB.Preload("Account").
		Joins("JOIN accounts ON accounts.id = clients.account_id").
		Where("clients.name LIKE ? OR accounts.email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Find(&clientList).Error; err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientList)
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
	}

	var acc models.Account
	if err := DB.Where("Email = ?", req.Email).First(&acc); err != nil {
		http.Error(w, "Email exists!", http.StatusBadRequest)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Hash failed", http.StatusInternalServerError)
		return
	}

	newAccount := models.Account{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := DB.Create(&newAccount).Error; err != nil {
		http.Error(w, "Create account failed", http.StatusInternalServerError)
		return
	}

	newClient := models.Client{
		Bio:       req.Bio,
		AccountID: newAccount.ID,
	}

	if err := DB.Create(&newClient).Error; err != nil {
		http.Error(w, "Create client failed", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&newClient)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var client models.Client

	if err := DB.Preload("Account").First(&client, id).Error; err != nil {
		http.Error(w, "This client does not exist!", http.StatusBadRequest)
		return
	}

	var req UpdateClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if req.Bio != nil {
		client.Bio = *req.Bio
	}

	if req.LinkInstagram != nil {
		client.LinkInstagram = *req.LinkInstagram
	}

	if req.LinkFacebook != nil {
		client.LinkFacebook = *req.LinkFacebook
	}

	if req.LinkWebsite != nil {
		client.LinkWebsite = *req.LinkWebsite
	}

	if err := DB.Save(&client).Error; err != nil {
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var deleteClient models.Client
	var deleteAccount models.Account
	if err := DB.First(&deleteClient, id).Error; err != nil {
		http.Error(w, "This client is not exists!", http.StatusBadRequest)
		return
	}

	deleteClient.IsDeleted = true
	if err := DB.Save(&deleteClient).Error; err != nil {
		http.Error(w, "Error when delete client", http.StatusBadRequest)
		return
	}

	if err := DB.First(&deleteAccount, deleteClient.AccountID).Error; err != nil {
		http.Error(w, "This account is not exists!", http.StatusBadRequest)
		return

	}

	deleteClient.IsDeleted = true
	if err := DB.Save(&deleteClient).Error; err != nil {
		http.Error(w, "Error when delete account!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleteClient)
}
