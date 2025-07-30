package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CategoryRequest struct {
	Name string `json:"name"`
}

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	var categoryList []models.Category
	DB.Find(&categoryList)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(categoryList)
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var category models.Category

	if err := DB.First(&category, id).Error; err != nil {
		http.Error(w, "This category is not exists!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	newCategory := models.Category{
		Name: req.Name,
	}

	if err := DB.Create(&newCategory).Error; err != nil {
		http.Error(w, "Error when create category", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCategory)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["Id"]

	var updateCategory models.Category

	if err := DB.First(&updateCategory, id).Error; err != nil {
		http.Error(w, "This id is not exists!", http.StatusBadRequest)
		return
	}

	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	updateCategory.Name = req.Name

	if err := DB.Save(&updateCategory).Error; err != nil {
		http.Error(w, "Error when update category", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateCategory)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	var deleteCategory models.Category
	if err := DB.First(&deleteCategory, id).Error; err != nil {
		http.Error(w, "This category is not exist!", http.StatusBadRequest)
		return
	}
	deleteCategory.IsDeleted = true
	if err := DB.Save(&deleteCategory).Error; err != nil {
		http.Error(w, "Error when delete category", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleteCategory)
}
