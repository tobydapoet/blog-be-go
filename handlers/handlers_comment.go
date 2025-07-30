package handlers

import (
	"blog-app/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateCommentRequest struct {
	CommenttableID   uint   `json:"commentTableId"`
	CommenttableType string `json:"commentTableType"`
	ClientID         uint   `json:"client_id"`
	Content          string `json:"content"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

func GetCommentsByType(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["id"]
	commentType := mux.Vars(r)["type"]

	var comments []models.Comment
	err := DB.Preload("Client").Where("client_id = ? AND commenttable_type = ?", clientID, commentType).Find(&comments).Error
	if err != nil {
		http.Error(w, "Lỗi khi truy vấn acomment", http.StatusInternalServerError)
		return
	}

	var ids []uint
	for _, comment := range comments {
		ids = append(ids, comment.CommenttableID)
	}
	w.Header().Set("Content-Type", "application/json")

	if len(ids) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	switch commentType {
	case "blog":
		{
			var blogs []models.Blog
			err = DB.Where("id IN ?", ids).Find(&blogs).Error
			if err != nil {
				http.Error(w, "invalid blogs", http.StatusInternalServerError)
				return
			}
			blogMap := make(map[uint]models.Blog)
			for _, b := range blogs {
				blogMap[b.ID] = b
			}

			for i := range comments {
				comments[i].Commenttable = blogMap[comments[i].CommenttableID]
			}
		}

	case "activity":
		{
			var activities []models.Activity
			err := DB.Where("id IN ?", ids).Find(&activities).Error
			if err != nil {
				http.Error(w, "invalid activities", http.StatusInternalServerError)
				return
			}

			accMap := make(map[uint]models.Activity)
			for _, acc := range activities {
				accMap[acc.ID] = acc
			}

			for i := range comments {
				comments[i].Commenttable = accMap[comments[i].CommenttableID]
			}
		}

	case "comment":
		{
			var comments []models.Comment
			err := DB.Where("id IN ?", ids).Find(&comments).Error
			if err != nil {
				http.Error(w, "invalid comments", http.StatusInternalServerError)
				return
			}
			cmtMap := make(map[uint]models.Comment)
			for _, cmt := range comments {
				cmtMap[cmt.ID] = cmt
			}

			for i := range comments {
				comments[i].Commenttable = cmtMap[comments[i].CommenttableID]
			}
		}

	default:
		http.Error(w, "Loại comment không hợp lệ", http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(comments)

}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	var newComment = models.Comment{
		CommenttableID:   req.CommenttableID,
		CommenttableType: req.CommenttableType,
		ClientID:         req.ClientID,
		Content:          req.Content,
	}

	if err := DB.Create(&newComment).Error; err != nil {
		http.Error(w, "Error when create comment", http.StatusBadRequest)
		return
	}

	if err := DB.Preload("Blog").Preload("Client").First(&newComment, newComment.ID).Error; err != nil {
		http.Error(w, "Error loading relations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var updateComment models.Comment

	if err := DB.First(&updateComment, id).Error; err != nil {
		http.Error(w, "This id is not exists!", http.StatusBadRequest)
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	updateComment.Content = req.Content

	if err := DB.Save(&updateComment).Error; err != nil {
		http.Error(w, "Error when update comment", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateComment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var comment models.Comment

	if err := DB.First(&comment, id).Error; err != nil {
		http.Error(w, "This id is not exists!", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid data!", http.StatusBadRequest)
		return
	}

	comment.IsDeleted = true

	if err := DB.Save(comment).Error; err != nil {
		http.Error(w, "Error when delete comment!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteComment)
}
