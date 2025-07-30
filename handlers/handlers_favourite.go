package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type FavouriteRequest struct {
	ClientID           uint   `json:"client_id"`
	BlogID             uint   `json:"blogId"`
	FavouritetableID   uint   `json:"favouriteTableId"`
	FavouritetableType string `json:"favouriteTableType"`
}

func GetFavouritesByClient(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["id"]
	favType := mux.Vars(r)["type"]

	var favourites []models.Favourite
	err := DB.Where("client_id = ? AND favouritetable_type = ?", clientID, favType).Find(&favourites).Error
	if err != nil {
		http.Error(w, "Lỗi truy vấn favourites", http.StatusInternalServerError)
		return
	}

	var ids []uint
	for _, fav := range favourites {
		ids = append(ids, fav.FavouritetableID)
	}

	w.Header().Set("Content-Type", "application/json")

	if len(ids) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	switch favType {
	case "blog":
		var blogs []models.Blog
		err = DB.Where("id IN ?", ids).Find(&blogs).Error
		if err != nil {
			http.Error(w, "Lỗi truy vấn blogs", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(blogs)

	case "activity":
		var activities []models.Activity
		err = DB.Where("id IN ?", ids).Find(&activities).Error
		if err != nil {
			http.Error(w, "Lỗi truy vấn activities", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(activities)

	case "comment":
		var comments []models.Comment
		err := DB.Where("id IN ?", ids).Find(&comments).Error
		if err != nil {
			http.Error(w, "Lỗi truy vấn comments", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(comments)

	default:
		http.Error(w, "Loại favourite không hợp lệ", http.StatusBadRequest)
	}
}

func GetClientsByFavourite(w http.ResponseWriter, r *http.Request) {
	targetID := mux.Vars(r)["id"]
	targetType := mux.Vars(r)["type"]

	var favourites []models.Favourite

	err := DB.Preload("Client").
		Where("favouritetable_id = ? AND favouritetable_type = ?", targetID, targetType).
		Find(&favourites).Error
	if err != nil {
		http.Error(w, "Lỗi khi truy vấn favourites", http.StatusInternalServerError)
		return
	}

	var clients []models.Client
	for _, fav := range favourites {
		clients = append(clients, fav.Client)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func CreateFavourite(w http.ResponseWriter, r *http.Request) {
	var req FavouriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	var existing models.Favourite
	if err := DB.Where("blog_id = ?, client_id = ?", req.BlogID, req.ClientID).First(&existing).Error; err != nil {
		http.Error(w, "Already like blog with this client", http.StatusBadRequest)
		return
	}

	newFavourite := models.Favourite{
		ClientID:           req.ClientID,
		FavouritetableID:   req.FavouritetableID,
		FavouritetableType: req.FavouritetableType,
	}

	if err := DB.Create(&newFavourite).Error; err != nil {
		http.Error(w, "Can't like this blog!", http.StatusBadRequest)
		return
	}

	if err := DB.Preload("Client").Preload("Blog").First(&newFavourite, newFavourite.ID).Error; err != nil {
		http.Error(w, "Error loading relations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newFavourite)
}

func DeleteFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var favourite models.Favourite
	if err := DB.First(&favourite, id).Error; err != nil {
		http.Error(w, "Can't find this favourite!", http.StatusBadRequest)
		return
	}

	if err := DB.Delete(&favourite).Error; err != nil {
		http.Error(w, "Delete failed!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favourite)
}
