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
	//connString := "host=localhost user=password password=password dbname=password1 sslmode=disable"
	/*SELECT s.schemaname,
	         s.relname AS tablename,
	         s.indexrelname AS indexname,
	         pg_relation_size(s.indexrelid) AS index_size,
	         s.idx_scan
	  FROM pg_catalog.pg_stat_user_indexes s
	     JOIN pg_catalog.pg_index i ON s.indexrelid = i.indexrelid
	  WHERE s.idx_scan < 10      -- has never been scanned
	    AND 0 <>ALL (i.indkey)  -- no index column is an expression
	    AND NOT i.indisunique   -- is not a UNIQUE index
	    AND NOT EXISTS          -- does not enforce a constraint
	           (SELECT 1 FROM pg_catalog.pg_constraint c
	            WHERE c.conindid = s.indexrelid)
	  ORDER BY pg_relation_size(s.indexrelid) DESC;*/
	/*SELECT s.schemaname,
	         s.relname AS tablename,
	         s.indexrelname AS indexname,
	         pg_relation_size(s.indexrelid) AS index_size,
	         s.idx_scan
	  FROM pg_catalog.pg_stat_user_indexes s
	     JOIN pg_catalog.pg_index i ON s.indexrelid = i.indexrelid
	  ORDER BY pg_relation_size(s.indexrelid) DESC;*/
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

	//router.GET("/debug/pprof/profile", pprofhandler.PprofHandler)

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
