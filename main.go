package main

import (
	"log"
	"net/http"

	"DB_Project_TP/config"
	"DB_Project_TP/pkg/server/handlers"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Server started!")
	models.TruncateAllTables()

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.HandleFunc("/forum/create", handlers.CreateForumHandler).Methods("POST")
	r.HandleFunc("/forum/{slug}/create", handlers.CreateThreadHandler).Methods("POST")
	r.HandleFunc("/forum/{slug}/details", handlers.ForumHandler).Methods("GET")
	r.HandleFunc("/forum/{slug}/threads", handlers.ForumsBranchsHandler).Methods("GET")
	r.HandleFunc("/forum/{slug}/users", handlers.ForumsUsersHandlers).Methods("GET")
	//
	r.HandleFunc("/post/{id}/details", handlers.PostFullHandler).Methods("GET")
	r.HandleFunc("/post/{id}/details", handlers.UpdatePostHandler).Methods("POST")
	//
	r.HandleFunc("/service/clear", handlers.ClearDB).Methods("POST")
	r.HandleFunc("/service/status", handlers.DBInfoHandler).Methods("GET")
	//
	r.HandleFunc("/thread/{slug_or_id}/create", handlers.CreatePostHandler).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/details", handlers.ThreadInfo).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/details", handlers.UpdateBranchHandler).Methods("POST")
	r.HandleFunc("/thread/{slug_or_id}/posts", handlers.SortedPostsHandler).Methods("GET")
	r.HandleFunc("/thread/{slug_or_id}/vote", handlers.BranchVoteHandler).Methods("POST")

	r.HandleFunc("/user/{nickname}/create", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/user/{nickname}/profile", handlers.ProfileHandler).Methods("GET")
	r.HandleFunc("/user/{nickname}/profile", handlers.UpdateProfileHandler).Methods("POST")

	log.Fatalln(http.ListenAndServe(":"+config.PORT, r))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}