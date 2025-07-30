package handlers

import (
	"blog-app/middlewares"
	"encoding/json"
	"net/http"
)

type jwtRes struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(uint)
	role := r.Context().Value(middlewares.RoleKey).(string)

	res := jwtRes{
		ID:   userID,
		Role: role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
