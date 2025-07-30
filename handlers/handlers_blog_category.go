package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type BLogCategoryRequest struct {
	BlogID     uint `json:"blogId"`
	CategoryID uint `json:"categoryId"`
}

func GetBlogByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	var blogCategories []models.BlogCategory
	if err := DB.Preload("Blog").Where("category_id = ?", categoryID).Find(&blogCategories).Error; err != nil {
		http.Error(w, "Error retrieving blogs", http.StatusInternalServerError)
		return
	}

	var blogs []models.Blog
	for _, bc := range blogCategories {
		blogs = append(blogs, bc.Blog)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func GetCategoryByBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID := vars["id"]

	var blogCategories []models.BlogCategory
	if err := DB.Preload("Category").Where("blog_id = ?", blogID).Find(&blogCategories).Error; err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	var categories []models.Category
	for _, bc := range blogCategories {
		categories = append(categories, bc.Category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func CreateBLogCategory(w http.ResponseWriter, r *http.Request) {
	var req BLogCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	var existing models.BlogCategory
	if err := DB.Where("blog_id = ? AND category_id = ?", req.BlogID, req.CategoryID).First(&existing).Error; err != nil {
		http.Error(w, "This Blog-Category relation already exists!", http.StatusBadRequest)
		return
	}

	var newBlogCategory = models.BlogCategory{
		BlogID:     req.BlogID,
		CategoryID: req.CategoryID,
	}

	if err := DB.Create(&newBlogCategory).Error; err != nil {
		http.Error(w, "Create failed!", http.StatusBadRequest)
		return
	}

	if err := DB.Preload("Blog").Preload("Category").
		First(&newBlogCategory, newBlogCategory.ID).Error; err != nil {
		http.Error(w, "Error loading relations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBlogCategory)
}

func DeleteBlogCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var blogCategory models.BlogCategory
	if err := DB.First(&blogCategory, id).Error; err != nil {
		http.Error(w, "Can't find this blog-category!", http.StatusBadRequest)
		return
	}

	if err := DB.Delete(&blogCategory).Error; err != nil {
		http.Error(w, "Delete failed!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogCategory)
}
