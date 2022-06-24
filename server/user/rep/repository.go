package rep

import (
	"DBMSForum/models"
	"fmt"
)

func Create(user models.User) error {
	_, err := models.DB.Exec("INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);",
		user.About, user.Email, user.Fullname, user.Nickname)

	return err
}

func CheckByEmail(email string) bool {
	var mail string
	err := models.DB.QueryRow("SELECT email FROM users WHERE email = $1;", email).Scan(&mail)
	return err == nil
}

func CheckByNickname(nickname string) bool {
	var name string
	err := models.DB.QueryRow("SELECT nickname FROM users WHERE nickname = $1;", nickname).Scan(&name)
	return err == nil
}

func ConflictUsers(email, nickname string) (*models.Users, error) {
	users := models.Users{}
	user := models.User{}

	rows, err := models.DB.Query("SELECT about, email, fullname, nickname FROM users WHERE nickname = $2 OR email = $1;", email, nickname)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}
func CheckByNicknameR(nickname string) (string, error) {
	err := models.DB.QueryRow("SELECT nickname FROM users WHERE nickname = $1;", nickname).
		Scan(&nickname)

	return nickname, err
}

func FindByNickname(nickname string) (*models.User, error) {
	user := models.User{}

	err := models.DB.QueryRow("SELECT about, email, fullname, nickname FROM users WHERE nickname = $1;", nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	return &user, err
}

func UpdateUser(nickname string, userUpdate models.UserUpdate) (*models.User, error) {
	values := make([]interface{}, 0, 3)
	i := 1
	var s string
	if userUpdate.About != "" {
		values = append(values, userUpdate.About)
		s += " about=$" + fmt.Sprint(i)
		i++
	}
	if userUpdate.Email != "" {
		values = append(values, userUpdate.Email)

		if i > 1 {
			s += ","
		}
		s += " email=$" + fmt.Sprint(i)
		i++
	}
	if userUpdate.Fullname != "" {
		values = append(values, userUpdate.Fullname)
		if i > 1 {
			s += ","
		}
		s += " fullname=$" + fmt.Sprint(i)
		i++
	}
	user := models.User{}

	user.Nickname = nickname

	if i > 1 {
		sqlRow := "UPDATE users SET" + s + " WHERE nickname = $" + fmt.Sprint(i) + " RETURNING fullname, about, email;"
		values = append(values, nickname)
		err := models.DB.QueryRow(sqlRow, values...).Scan(&user.Fullname, &user.About, &user.Email)
		return &user, err
	} else {
		err := models.DB.QueryRow("SELECT about, email, fullname FROM users WHERE nickname = $1;", nickname).
			Scan(&user.About, &user.Email, &user.Fullname)
		return &user, err
	}
}

func FindByForum(slug, since string, limit int, desc bool) (*models.Users, error) {
	users := models.Users{}

	symb := '>'
	descS := ""
	if desc {
		symb = '<'
		descS = "DESC "
	}
	sqlRek := "SELECT users.about, users.email, users.fullname, users.nickname FROM forum_users " +
		"JOIN users ON users.nickname = forum_users.author WHERE slug = $1 "

	values := make([]interface{}, 0, 2)

	values = append(values, slug)
	if since != "" {
		sqlRek += fmt.Sprintf(`AND users.nickname %c $2 COLLATE "C" `, symb)
		values = append(values, since)
	}
	sqlRek += fmt.Sprintf(`ORDER BY users.nickname COLLATE "C" %sLIMIT %d;`, descS, limit)

	rows, err := models.DB.Query(sqlRek, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	user := models.User{}
	for rows.Next() {
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}
