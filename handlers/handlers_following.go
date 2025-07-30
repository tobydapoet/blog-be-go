package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FollowReq struct {
	ClientID    uint `json:"client_id"`
	FollowingID uint `json:"following_id"`
}

func CreateFollow(w http.ResponseWriter, r *http.Request) {
	var follow models.Following
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existing models.Following
	if err := DB.Where("client_id = ? AND following_id = ?", follow.ClientID, follow.FollowingID).First(&existing).Error; err == nil {
		http.Error(w, "You already follow this user", http.StatusConflict)
		return
	}

	if err := DB.Create(&follow).Error; err != nil {
		http.Error(w, "Failed to follow", http.StatusInternalServerError)
		return
	}

	DB.Preload("Client").Preload("FollowedUser").First(&follow, follow.ID)
	json.NewEncoder(w).Encode(follow)
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	var follow models.Following
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := DB.Where("client_id = ? AND following_id = ?", follow.ClientID, follow.FollowingID).Delete(&models.Following{}).Error; err != nil {
		http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Unfollow successful",
	})
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	clientID, _ := strconv.Atoi(params["id"])

	var follows []models.Following
	DB.Preload("Client").Where("following_id = ?", clientID).Find(&follows)

	json.NewEncoder(w).Encode(follows)
}

func GetFollowings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	clientID, _ := strconv.Atoi(params["id"])

	var follows []models.Following
	DB.Preload("FollowedUser").Where("client_id = ?", clientID).Find(&follows)

	json.NewEncoder(w).Encode(follows)
}
