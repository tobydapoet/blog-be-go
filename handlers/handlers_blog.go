package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateBlogRequest struct {
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ClientID  uint   `json:"client_id"`
}

type UpdateBlogRequest struct {
	Thumbnail *string `json:"thumbnail"`
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	ClientID  *uint   `json:"client_id"`
}

func GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	var blogList []models.Blog
	DB.Find(&blogList)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogList)
}

func GetBlogById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var blog models.Blog

	if err := DB.Find(&blog, id).Error; err != nil {
		http.Error(w, "This blog is not exists!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func GetBlogByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	var blog []models.Blog

	err := DB.Preload("Client").Preload("Client.Account").
		Joins("JOIN clients ON clients.id = blogs.client_id").
		Joins("JOIN accounts ON accounts.id = clients.account_id").
		Where("accounts.email = ?", email).
		Find(&blog).Error

	if err != nil {
		http.Error(w, "Cannot fetch blog", http.StatusInternalServerError)
		return
	}

	if len(blog) == 0 {
		http.Error(w, "No blog found for this email", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func GetBlogByUser(w http.ResponseWriter, r *http.Request) {
	client_id := mux.Vars(r)["id"]

	var blogs []models.Blog

	if err := DB.Where("client_id = ?", client_id).Find(&blogs).Error; err != nil {
		http.Error(w, "Can't find any blog!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func SearchBlogs(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	var blogList []models.Blog

	if err := DB.Preload("Client").
		Joins("JOIN clients ON clients.id = blog.client_id").
		Where("client.name LIKE ? OR blog.title LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Find(&blogList).Error; err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogList)
}

func CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req CreateBlogRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	newBlog := models.Blog{
		Thumbnail: req.Thumbnail,
		Title:     req.Title,
		Content:   req.Content,
		ClientID:  req.ClientID,
	}

	if err := DB.Create(&newBlog).Error; err != nil {
		http.Error(w, "Error when creating blog!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBlog)
}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Blog
	if err := DB.First(&blog, id).Error; err != nil {
		http.Error(w, "Blog not found!", http.StatusNotFound)
		return
	}

	var req UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}

	if req.Thumbnail != nil {
		blog.Thumbnail = *req.Thumbnail
	}
	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.ClientID != nil {
		blog.ClientID = *req.ClientID
	}

	if err := DB.Save(&blog).Error; err != nil {
		http.Error(w, "Failed to update blog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func ApproveBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Blog
	if err := DB.First(&blog, id).Error; err != nil {
		http.Error(w, "Blog not found!", http.StatusNotFound)
		return
	}

	if blog.Status != "wait approve" {
		http.Error(w, "Blog is not in a state to approve", http.StatusBadRequest)
		return
	}

	blog.Status = "approve"
	if err := DB.Save(&blog).Error; err != nil {
		http.Error(w, "Error when delete blog!", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func CancelBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Blog
	if err := DB.First(&blog, id).Error; err != nil {
		http.Error(w, "Blog not found!", http.StatusNotFound)
		return
	}

	if blog.Status != "wait approve" {
		http.Error(w, "Blog is not in a state to deny", http.StatusBadRequest)
		return
	}

	blog.Status = "denied"
	if err := DB.Save(&blog).Error; err != nil {
		http.Error(w, "Error when delete blog!", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blog models.Blog
	if err := DB.First(&blog, id).Error; err != nil {
		http.Error(w, "Blog not found!", http.StatusNotFound)
		return
	}

	blog.IsDeleted = false
	if err := DB.Save(&blog).Error; err != nil {
		http.Error(w, "Error when delete blog!", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}
