package rep

import (
	"DBMSForum/models"
	"fmt"
)

// !!!!!
func InsertThread(thread *models.Thread) error {
	err := models.DB.QueryRow("INSERT INTO threads(author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, votes;",
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title).
		Scan(&thread.ID, &thread.Votes)

	return err
}

func FindThreadIDForum(slug string) (int, string, error) {
	var tid int
	var ff string
	err := models.DB.QueryRow("SELECT id, forum FROM threads WHERE slug = $1;", slug).
		Scan(&tid, &ff)

	return tid, ff, err
}

func FindThreadByIDIDForum(id int) (int, string, error) {
	var tid int
	var ff string
	err := models.DB.QueryRow("SELECT id, forum FROM threads WHERE id = $1;", id).
		Scan(&tid, &ff)

	return tid, ff, err
}

func FindThreadID(slug string) (int, error) {
	var tid int
	err := models.DB.QueryRow("SELECT id FROM threads WHERE slug = $1;", slug).
		Scan(&tid)

	return tid, err
}

func FindThreadByIDID(id int) (int, error) {
	var tid int
	err := models.DB.QueryRow("SELECT id FROM threads WHERE id = $1;", id).
		Scan(&tid)

	return tid, err
}

func FindThread(slug string) (*models.Thread, error) {
	th := models.Thread{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE slug = $1;", slug).
		Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)

	return &th, err
}

func FindThreadByID(id int) (*models.Thread, error) {
	th := models.Thread{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE id = $1;", id).
		Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)

	return &th, err
}

func FindThreads(forum, since string, limit int, desc bool) (*models.Threads, error) {
	threads := models.Threads{}
	values := make([]interface{}, 0, 2)
	symb := '>'
	descS := ""
	if desc {
		symb = '<'
		descS = "DESC "
	}

	sqlRow := "SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE forum = $1 "
	values = append(values, forum)
	if since != "" {
		sqlRow += fmt.Sprintf("AND created %c= $2 ", symb)
		values = append(values, since)
	}
	sqlRow += fmt.Sprintf("ORDER BY created %sLIMIT %d;", descS, limit)
	rows, err := models.DB.Query(sqlRow, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	th := models.Thread{}

	for rows.Next() {
		err = rows.Scan(&th.ID, &th.Author, &th.Created, &th.Forum, &th.Message, &th.Slug, &th.Title, &th.Votes)
		if err != nil {
			return nil, err
		}
		threads = append(threads, th)
	}
	return &threads, nil
}

func UpdateThread(thread *models.Thread, threadUpdate models.ThreadUpdate) error {
	values := make([]interface{}, 0, 3)
	i := 1
	var s string
	if threadUpdate.Message != "" {
		values = append(values, threadUpdate.Message)
		s += " message=$" + fmt.Sprint(i)
		i++
	}
	if threadUpdate.Title != "" {
		values = append(values, threadUpdate.Title)

		if i > 1 {
			s += ","
		}
		s += " title=$" + fmt.Sprint(i)
		i++
	}

	if i > 1 {
		sqlRow := "UPDATE threads SET" + s + " WHERE slug = $" + fmt.Sprint(i) + " RETURNING message, title;"
		values = append(values, thread.Slug)
		err := models.DB.QueryRow(sqlRow, values...).
			Scan(&thread.Message, &thread.Title)

		return err
	}
	//} else {
	//	//err := models.DB.QueryRow("SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE slug = $1;", thread.Slug).
	//	//	Scan(&thread.ID, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	//
	//	//return err
	//}
	return nil
}

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec("INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, NULLIF($3, 0));", vote.Nickname, vote.Voice, vote.Thread)

	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec("UPDATE votes SET voice=$1 WHERE thread=$2 AND nickname = $3;", vote.Voice, vote.Thread, vote.Nickname)

	return err
}
