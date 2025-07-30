package handlers

import (
	"blog-app/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/datatypes"
)

type createAcRequest struct {
	ClientID uint           `json:"client_id"`
	Content  string         `json:"content"`
	Images   datatypes.JSON `json:"images"`
}

type updateAcRequest struct {
	Content *string         `json:"content"`
	Images  *datatypes.JSON `json:"images"`
}

func GetActivityByUser(w http.ResponseWriter, r *http.Request) {
	client_id := mux.Vars(r)["id"]

	var activities []models.Activity

	if err := DB.Preload("Client").Preload("Client.Account").Where("client_id = ?", client_id).Find(&activities).Error; err != nil {
		http.Error(w, "Can't find any activities", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func GetActivityByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	var activity []models.Activity

	err := DB.Preload("Client").Preload("Client.Account").
		Joins("JOIN clients ON clients.id = activities.client_id").
		Joins("JOIN accounts ON accounts.id = clients.account_id").
		Where("accounts.email = ?", email).
		Find(&activity).Error

	if err != nil {
		http.Error(w, "Cannot fetch activity", http.StatusInternalServerError)
		return
	}

	if len(activity) == 0 {
		http.Error(w, "No activity found for this email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

func GetAllActivities(w http.ResponseWriter, r *http.Request) {
	var activities []models.Activity
	DB.Preload("Client").Find(&activities)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func SearchActivities(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	var activities []models.Activity

	if err := DB.Preload("Client").
		Joins("JOIN clients ON clients.id = activities.client_id").
		Where("client.name LIKE ? OR activities.content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Find(&activities).Error; err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func GetActivityById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var activity models.Activity
	if err := DB.First(&activity, id).Error; err != nil {
		http.Error(w, "Can't find any activity", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	var req createAcRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		http.Error(w, "Failed to parse images", http.StatusBadRequest)
		return
	}

	newActivity := models.Activity{
		ClientID: req.ClientID,
		Content:  req.Content,
		Images:   datatypes.JSON(imagesJSON),
	}

	fmt.Printf("==> client_id: %d\n", req.ClientID)
	fmt.Printf("==> content: %s\n", req.Content)
	fmt.Printf("==> images: %+v\n", req.Images)

	if err := DB.Create(&newActivity).Error; err != nil {
		http.Error(w, "Error when creating activity!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newActivity)
}

func UpdateActivity(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var activity models.Activity
	if err := DB.First(&activity, id).Error; err != nil {
		http.Error(w, "Can't find this activity!", http.StatusBadGateway)
		return
	}

	var req updateAcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadGateway)
		return
	}

	if req.Content != nil {
		activity.Content = *req.Content
	}

	if req.Images != nil {
		activity.Images = *req.Images
	}

	if err := DB.Save(&activity).Error; err != nil {
		http.Error(w, "Error when update activity!", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)

}

func DeleteActivity(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var activity models.Activity
	if err := DB.First(&activity, id).Error; err != nil {
		http.Error(w, "Can't find this activity!", http.StatusBadGateway)
		return
	}
	activity.IsDeleted = true
	if err := DB.Save(&activity).Error; err != nil {
		http.Error(w, "Error when update activity!", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}
