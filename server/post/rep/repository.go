package rep

import (
	"DBMSForum/models"
	"fmt"
)

func InsertPosts(pos *models.Posts, id int, forum string) error {
	if len(*pos) == 0 {
		return nil
	}
	query := "INSERT INTO posts(author, forum, message, thread, parent) VALUES "
	i := 1
	values := make([]interface{}, 0, len(*pos))
	authors := ""
	for _, post := range *pos {
		query += fmt.Sprintf("('%s', '%s', '%s', %d, $%d), ",
			post.Author, forum, post.Message, id, i)
		i++
		values = append(values, post.Parent)
		authors += fmt.Sprintf("('%s', '%s'),",
			forum, post.Author)
	}

	rows, err := models.DB.Query(query[:len(query)-2]+" RETURNING created, id;", values...)
	if err != nil {
		return err
	}

	i = 0
	for rows.Next() {
		err = rows.Scan(&(*pos)[i].Created, &(*pos)[i].ID)
		(*pos)[i].Forum = forum
		(*pos)[i].Thread = id
		if err != nil {
			return err
		}
		i++
	}

	//_, err = models.DB.Exec("UPDATE forums SET Posts=(Posts + $1) WHERE slug = $2;", len(*pos), forum)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	_, err = models.DB.Exec("INSERT INTO forum_users (slug, author) VALUES" + authors[:len(authors)-1] + " ON CONFLICT DO NOTHING;")

	return err
}

func FindPosts(author int, limit, since int, sort string, desc bool) (*models.Posts, error) {
	posts := models.Posts{}

	sqlRec := "SELECT p.author, p.created, p.forum, p.id, p.message, p.parent, p.thread FROM posts AS p "
	symb := ">"
	descS := ""
	if desc {
		symb = "<"
		descS = " DESC"
	}
	and := ""
	if sort == "flat" || sort == "" {
		if since != 0 {
			and = fmt.Sprintf("AND p.id %s %d ", symb, since)
		}
		sqlRec += fmt.Sprintf("WHERE p.thread = %d %sORDER BY p.created %s, p.id %s LIMIT %d;", author, and, descS, descS, limit)
	} else if sort == "tree" {
		if since != 0 {
			and = fmt.Sprintf("AND p.PATH %s (SELECT path FROM posts WHERE id = %d) ", symb, since)
		}
		sqlRec += fmt.Sprintf("WHERE p.thread = %d %sORDER BY p.path %s LIMIT %d;", author, and, descS, limit)
	} else {
		if since != 0 {
			and = fmt.Sprintf("AND PATH[1] %s (SELECT path[1] FROM posts WHERE id = %d) ", symb, since)
		}
		//sqlRec += fmt.Sprintf("WHERE p.path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent IS NULL %sORDER BY id %s LIMIT %d)ORDER BY p.path[1] %s, p.path;",
		//	author, and, descS, limit, descS)
		sqlRec += fmt.Sprintf("JOIN (SELECT path FROM posts WHERE thread = %d AND parent = 0 %sORDER BY id %s LIMIT %d) AS prnt ON prnt.path[1] = p.path[1] ORDER BY p.path[1] %s, p.path;",
			author, and, descS, limit, descS)
	}

	rows, err := models.DB.Query(sqlRec)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p := models.Post{}

	for rows.Next() {
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return &posts, nil
}

func FindByID(id int) (*models.Post, error) {
	post := models.Post{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, is_edited, message, parent, thread FROM posts WHERE id = $1;", id).
		Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

	return &post, err
}

func UpdatePost(post *models.Post, postUpdate models.PostUpdate) error {
	if postUpdate.Message != "" && postUpdate.Message != post.Message {
		err := models.DB.QueryRow("UPDATE posts SET message=$1, is_edited=true WHERE id=$2 RETURNING is_edited, message;", postUpdate.Message, post.ID).
			Scan(&post.IsEdited, &post.Message)
		return err
	}
	return nil
}
