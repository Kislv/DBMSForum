package main

import (
	"DBMSForum/models"
	"DBMSForum/server/forum"
	"DBMSForum/server/post"
	"DBMSForum/server/thread"
	"DBMSForum/server/user"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func DBConnection() *sql.DB {
	connString := "host=localhost user=vitya password=password dbname=forum sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	models.DB = DBConnection()

	router := router.New()


	router.POST("/api/forum/create", forum.Create)
	router.GET("/api/forum/{slug}/details", forum.Details)
	router.POST("/api/forum/{slug}/create", thread.Create)
	router.GET("/api/forum/{slug}/users", forum.Users)
	router.GET("/api/forum/{slug}/threads", forum.Threads)

	router.GET("/api/post/{id:[0-9]+}/details", post.Details)
	router.POST("/api/post/{id:[0-9]+}/details", post.DetailsPOST)

	router.POST("/api/service/clear", forum.ClearHandler)
	router.GET("/api/service/status", forum.StatusHandler)

	router.POST("/api/thread/{slug_or_id}/create", post.Create)
	router.GET("/api/thread/{slug_or_id}/details", thread.Details)
	router.POST("/api/thread/{slug_or_id}/details", thread.DetailsPOST)
	router.GET("/api/thread/{slug_or_id}/posts", post.ThreadPosts)
	router.POST("/api/thread/{slug_or_id}/vote", thread.Vote)

	router.POST("/api/user/{nickname}/create", user.Create)
	router.GET("/api/user/{nickname}/profile", user.Profile)
	router.POST("/api/user/{nickname}/profile", user.ProfilePOST)

	fmt.Println("Starting server at: 5000")
	fasthttp.ListenAndServe(":5000", router.Handler)
}
