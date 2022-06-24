package rep

import (
	"DBMSForum/models"
)

func InsertForum(forum models.Forum) error {
	_, err := models.DB.Exec("INSERT INTO forums(username, slug, title) VALUES ($1, $2, $3);",
		forum.User, forum.Slug, forum.Title)
	return err
}

func CheckForum(slug string) bool {
	var slug2 string
	err := models.DB.QueryRow("SELECT slug FROM forums WHERE slug = $1;", slug).Scan(&slug2)
	return err == nil
}

func FindForum(slug string) (*models.Forum, error) {
	f := models.Forum{}
	err := models.DB.QueryRow("SELECT username, posts, threads, slug, title FROM forums WHERE slug = $1;", slug).
		Scan(&f.User, &f.Posts, &f.Threads, &f.Slug, &f.Title)
	return &f, err
}

func ClearForum() error {
	var err error
	_, err = models.DB.Exec("TRUNCATE forums, posts, threads, users, votes, forum_users CASCADE;")
	return err
}

func StatusForum() models.Status {
	status := models.Status{}
	models.DB.QueryRow("SELECT COUNT(*) FROM forums;").Scan(&status.Forum)
	models.DB.QueryRow("SELECT COUNT(*) FROM posts;").Scan(&status.Post)
	models.DB.QueryRow("SELECT COUNT(*) FROM threads;").Scan(&status.Thread)
	models.DB.QueryRow("SELECT COUNT(*) FROM users;").Scan(&status.User)
	return status
}
