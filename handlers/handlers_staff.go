package handlers

import (
	"blog-app/models"
	"blog-app/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateStaffRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type UpdateStaffRequest struct {
	Phone string `json:"phone"`
}

func GetAllStaffs(w http.ResponseWriter, r *http.Request) {
	var staffList []models.Staff
	DB.Preload("Account").Find(&staffList)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staffList)
}

func GetStaffById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var staff models.Staff

	if err := DB.Preload("Account").First(&staff, id).Error; err != nil {
		http.Error(w, "This staff is not exists!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func GetStaffByAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	var staff models.Staff

	if err := DB.Preload(`Account`).Where(`account_id = ?`, accountID).First(&staff).Error; err != nil {
		http.Error(w, "Error when get staff!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func GetStaffByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	var staff models.Staff

	err := DB.
		Preload("Account").
		Joins("JOIN accounts ON accounts.id = staffs.account_id").
		Where("accounts.email = ?", email).
		First(&staff).Error

	if err != nil {
		http.Error(w, "Staff not found", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(staff)
}

func SearchStaffs(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	var staffList []models.Staff

	if err := DB.Preload("Account").
		Joins("JOIN accounts ON accounts.id = staffs.account_id").
		Where("staffs.name LIKE ? OR accounts.email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Find(&staffList).Error; err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staffList)
}

func CreateStaff(w http.ResponseWriter, r *http.Request) {
	var req CreateStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	var acc models.Account
	if err := DB.Where("Email = ?", req.Email).First(&acc).Error; err == nil {
		http.Error(w, "Email exists!", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Hash failed", http.StatusInternalServerError)
		return
	}

	newAccount := models.Account{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "staff",
	}

	if err := DB.Create(&newAccount).Error; err != nil {
		http.Error(w, "Create account failed", http.StatusInternalServerError)
		return
	}

	newStaff := models.Staff{
		Phone:     req.Phone,
		AccountID: newAccount.ID,
	}

	if err := DB.Where("phone = ?", req.Phone).First(&newStaff).Error; err == nil {
		http.Error(w, "Phone is exists!", http.StatusBadRequest)
		return
	}

	if err := DB.Create(&newStaff).Error; err != nil {
		http.Error(w, "Create staff failed", http.StatusBadRequest)
		return
	}

	if err := DB.Preload("Account").First(&newStaff, newStaff.ID).Error; err != nil {
		http.Error(w, "Failed to load staff with account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&newStaff)
}

func UpdateStaff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var staff models.Staff

	if err := DB.Preload("Account").First(&staff, id).Error; err != nil {
		http.Error(w, "This staff does not exist!", http.StatusBadRequest)
		return
	}

	var req UpdateStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	staff.Phone = req.Phone

	if err := DB.Save(&staff).Error; err != nil {
		http.Error(w, "Failed to update staff!", http.StatusInternalServerError)
		return
	}

	if err := DB.Preload("Account").First(&staff, staff.ID).Error; err != nil {
		http.Error(w, "Failed to load updated staff!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func DeleteStaff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var deleteStaff models.Staff
	var deleteAccount models.Account
	if err := DB.First(&deleteStaff, id).Error; err != nil {
		http.Error(w, "This staff is not exists!", http.StatusBadRequest)
		return
	}

	deleteStaff.IsDeleted = true
	if err := DB.Save(&deleteStaff).Error; err != nil {
		http.Error(w, "Error when delete staff", http.StatusBadRequest)
		return
	}

	if err := DB.First(&deleteAccount, deleteStaff.AccountID).Error; err != nil {
		http.Error(w, "This account is not exists!", http.StatusBadRequest)
	}

	deleteStaff.IsDeleted = true
	if err := DB.Save(&deleteStaff).Error; err != nil {
		http.Error(w, "Error when delete account!", http.StatusBadRequest)
	}
	DB.Preload("Account").First(&deleteStaff)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleteStaff)
}
