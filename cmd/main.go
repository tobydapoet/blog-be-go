package main

import (
	"blog-app/config"
	"blog-app/handlers"
	"blog-app/middlewares"
	"blog-app/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	fmt.Println("Connect success!")
	handlers.InitDatabase(config.DB)

	config.DB.AutoMigrate(
		&models.Account{},
		&models.Staff{},
		&models.Client{},
		&models.Blog{},
		&models.Category{},
		&models.BlogCategory{},
		&models.Comment{},
		&models.Activity{},
		&models.Favourite{},
		&models.BlackListToken{},
		&models.Following{},
	)

	router := mux.NewRouter().StrictSlash(true)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CORS_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	//Uploads
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	router.HandleFunc("/upload/{folder}", handlers.UploadImage).Methods("POST")

	//Autho
	router.HandleFunc("/oauth2/google/login", handlers.HandleGoogleLogin).Methods("GET")
	router.HandleFunc("/oauth2/callback", handlers.HandleGoogleCallback).Methods("GET")

	//Account
	router.Handle("/account", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.GetAllAccounts)))).Methods("GET")
	router.HandleFunc("/account/login", handlers.Login).Methods("POST")
	router.HandleFunc("/account/register", handlers.Register).Methods("POST")
	router.Handle("/account/protected", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.ProtectedHandler))).Methods("GET")
	router.Handle("/account/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetAccountById))).Methods("GET")
	router.Handle("/account/email/{email}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetAccountByEmail))).Methods("GET")
	router.Handle("/account/logout", middlewares.RefreshJWTMiddleware(http.HandlerFunc(handlers.Logout))).Methods("POST")
	router.Handle("/account/refresh", middlewares.RefreshJWTMiddleware(http.HandlerFunc(handlers.RefreshToken))).Methods("POST")
	router.Handle("/account/update/{id}", middlewares.RefreshJWTMiddleware(http.HandlerFunc(handlers.UpdateAccount))).Methods("POST")

	//Staff
	router.Handle("/staff", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.GetAllStaffs)))).Methods("GET")
	router.Handle("/staff/search", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.SearchStaffs)))).Methods("GET")
	router.Handle("/staff/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetStaffById))).Methods("GET")
	router.Handle("/staff/email/{email}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetStaffByEmail))).Methods("GET")
	router.Handle("/staff/user/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetStaffByAccount))).Methods("GET")
	router.Handle("/staff/create", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.CreateStaff)))).Methods("POST")
	router.Handle("/staff/update/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.UpdateStaff)))).Methods("PUT")
	router.Handle("/staff/delete/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin")(http.HandlerFunc(handlers.DeleteStaff)))).Methods("PUT")

	//Client
	router.Handle("/client", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.GetAllClients)))).Methods("GET")
	router.Handle("/client/search", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.SearchClients)))).Methods("GET")
	router.Handle("/client/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetClientById))).Methods("GET")
	router.Handle("/client/email/{email}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetClientByEmail))).Methods("GET")
	router.Handle("/client/user/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetClientByAccount))).Methods("GET")
	router.Handle("/client/create", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.CreateClient)))).Methods("POST")
	router.Handle("/client/update/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.UpdateClient)))).Methods("PUT")
	router.Handle("/client/delete/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.DeleteClient)))).Methods("PUT")

	//Category
	router.Handle("/category", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetAllCategory))).Methods("GET")
	router.Handle("/category/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetCategoryById))).Methods("GET")
	router.Handle("/category/create", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.CreateCategory)))).Methods("POST")
	router.Handle("/category/update/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.UpdateCategory)))).Methods("PUT")
	router.Handle("/category/delete/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.DeleteCategory)))).Methods("PUT")

	//Blog-category
	router.Handle("/blog_category/blog/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetCategoryByBlog))).Methods("GET")
	router.Handle("/blog_category/category/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetBlogByCategory))).Methods("GET")
	router.Handle("/blog_category/create", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.CreateBLogCategory)))).Methods("POST")
	router.Handle("/blog_category/delete/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.DeleteBlogCategory)))).Methods("DELETE")

	//Blog
	router.Handle("/blog", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetAllBlogs))).Methods("GET")
	router.Handle("/blog/search", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.SearchBlogs))).Methods("GET")
	router.Handle("/blog/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetBlogById))).Methods("GET")
	router.Handle("/blog/user/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetBlogByUser))).Methods("GET")
	router.Handle("/blog/email/{email}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetBlogByEmail))).Methods("GET")
	router.Handle("/blog/create", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff", "client")(http.HandlerFunc(handlers.CreateBlog)))).Methods("POST")
	router.Handle("/blog/update/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff", "client")(http.HandlerFunc(handlers.UpdateBlog)))).Methods("PUT")
	router.Handle("/blog/approve/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.ApproveBlog)))).Methods("PUT")
	router.Handle("/blog/cancel/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff")(http.HandlerFunc(handlers.CancelBlog)))).Methods("PUT")
	router.Handle("/blog/delete/{id}", middlewares.AccessJWTMiddleware(middlewares.RequireRoles("admin", "staff", "client")(http.HandlerFunc(handlers.DeleteBlog)))).Methods("PUT")

	//Activity
	router.Handle("/activity", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetAllActivities))).Methods("GET")
	router.Handle("/activity/search", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.SearchActivities))).Methods("GET")
	router.Handle("/activity/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetActivityById))).Methods("GET")
	router.Handle("/activity/email/{email}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetActivityByEmail))).Methods("GET")
	router.Handle("/activity/user/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetActivityByUser))).Methods("GET")
	router.Handle("/activity/create", middlewares.AccessJWTMiddleware(http.HandlerFunc(http.HandlerFunc(handlers.CreateActivity)))).Methods("POST")
	router.Handle("/activity/update/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(http.HandlerFunc(handlers.UpdateActivity)))).Methods("PUT")
	router.Handle("/activity/delete/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(http.HandlerFunc(handlers.DeleteActivity)))).Methods("PUT")

	//Comment
	router.Handle("/comment/{type}/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetCommentsByType))).Methods("GET")
	router.Handle("/comment/create", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.CreateComment))).Methods("POST")
	router.Handle("/comment/update/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.UpdateComment))).Methods("PUT")
	router.Handle("/comment/delete/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.DeleteComment))).Methods("PUT")

	//Favourite
	router.Handle("/favourite/{type}/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetClientsByFavourite))).Methods("GET")
	router.Handle("/favourite/client/{id}/{type}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetFavouritesByClient))).Methods("GET")
	router.Handle("/favourite/create", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.CreateFavourite))).Methods("POST")
	router.Handle("/favourite/delete/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.DeleteFavourite))).Methods("DELETE")

	//Following
	router.Handle("/follow/follower/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetFollowers))).Methods("GET")
	router.Handle("/follow/following/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.GetFollowings))).Methods("GET")
	router.Handle("/follow/create", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.CreateFollow))).Methods("POST")
	router.Handle("/follow/delete/{id}", middlewares.AccessJWTMiddleware(http.HandlerFunc(handlers.Unfollow))).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
