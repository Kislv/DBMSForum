package models

import (
	"database/sql"
	"github.com/jackc/pgtype"
	json "github.com/mailru/easyjson"
	"time"
)

var DB *sql.DB

type Error struct {
	Message string `json:"message"`
}

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

//easyjson:json
type Users []User

type UserUpdate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type Thread struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

//easyjson:json
type Threads []Thread

type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Post struct {
	ID       int              `json:"id"`
	Parent   int64            `json:"parent"`
	Author   string           `json:"author"`
	Message  string           `json:"message"`
	IsEdited bool             `json:"isEdited"`
	Forum    string           `json:"forum"`
	Thread   int              `json:"thread,"`
	Created  time.Time        `json:"created"`
	Path     pgtype.Int8Array `json:"-"`
}

//easyjson:json
type Posts []Post

type PostUpdate struct {
	Message string `json:"message"`
}

type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"-"`
}
//
//type JsonNullInt64 struct {
//	sql.NullInt64
//}
//
//func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
//	if v.Valid {
//		return json.Marshal(v.Int64)
//	} else {
//		return json.Marshal(nil)
//	}
//}
//
//func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
//	var x *int64
//	if err := json.Unmarshal(data, &x); err != nil {
//		return err
//	}
//	if x != nil {
//		v.Valid = true
//		v.Int64 = *x
//	} else {
//		v.Valid = false
//	}
//	return nil
//}

func MarshalErrorSt(message string) []byte {
	jsonError, err := json.Marshal(Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}
