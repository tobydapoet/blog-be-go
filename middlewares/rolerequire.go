package middlewares

import "net/http"

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(RoleKey)
			if roleVal == nil {
				http.Error(w, "Missing role in token", http.StatusForbidden)
				return
			}

			role := roleVal.(string)

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Access denied - insufficient permissions", http.StatusForbidden)
		})
	}
}
