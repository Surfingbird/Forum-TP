package main

import (
	"log"
	"net/http"

	"DB_Project_TP/config"
	"DB_Project_TP/pkg/server/handlers"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func main() {
	defer config.Logger.Sync()

	log.Println("Server started!")
	err := config.DB.Ping()
	if err != nil {
		log.Fatalf("Can not connect to DB: %v", err.Error())
	}

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.Use(ContentTypeMiddleware)

	forum := r.Group("/api/forum")
	forum.POST("/:root", CrunchHandler)
	forum.POST("/:root/:branch", CrunchHandler)
	forum.GET("/:slug/details", handlers.ForumHandler)
	forum.GET("/:slug/threads", handlers.ForumsBranchsHandler)
	forum.GET("/:slug/users", handlers.ForumsUsersHandlers)

	post := r.Group("/api/post")
	post.GET("/:id/details", handlers.PostFullHandler)
	post.POST("/:id/details", handlers.UpdatePostHandler)

	service := r.Group("/api/service")
	service.POST("/clear", handlers.ClearDB)
	service.GET("/status", handlers.DBInfoHandler)

	thread := r.Group("/api/thread")
	thread.POST("/:slug_or_id/create", handlers.CreatePostHandler)
	thread.GET("/:slug_or_id/details", handlers.ThreadInfo)
	thread.POST("/:slug_or_id/details", handlers.UpdateBranchHandler)
	thread.GET("/:slug_or_id/posts", handlers.SortedPostsHandler)
	thread.POST("/:slug_or_id/vote", handlers.BranchVoteHandler)

	user := r.Group("api/user")
	user.POST("/:nickname/create", handlers.CreateUserHandler)
	user.GET("/:nickname/profile", handlers.ProfileHandler)
	user.POST("/:nickname/profile", handlers.UpdateProfileHandler)

	r.GET("/delete", func(c *gin.Context) {
		models.TruncateAllTables()
	})

	log.Fatalln(http.ListenAndServe(":"+config.PORT, r))
}

func ContentTypeMiddleware(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	c.Next()
}

func CrunchHandler(c *gin.Context) {
	root := c.Param("root")
	branch := c.Param("branch")

	if root == "create" {
		handlers.CreateForumHandler(c)
	}

	if branch == "create" && root != "" {
		handlers.CreateThreadHandler(c, root)
	}
}
