package user

import (
	"DBMSForum/models"
	"DBMSForum/server/user/rep"
	"database/sql"

	json "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		// log.Println(err)
		return
	}

	user.Nickname = ctx.UserValue("nickname").(string)

	users, err := rep.ConflictUsers(user.Email, user.Nickname)
	if err != nil {
		// log.Println(err)
		return
	}

	if len(*users) > 0 {
		body, err := json.Marshal(users)
		if err != nil {
			// log.Println(err)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.Write(body)
		return
	}

	err = rep.Create(user)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(body)
}

func Profile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	user, err := rep.FindByNickname(nickname)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find user"))
			return
		}
		// log.Println(err)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
	return

}

func ProfilePOST(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	userUpdate := models.UserUpdate{}
	err := json.Unmarshal(ctx.PostBody(), &userUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	if userUpdate.Email != "" && rep.CheckByEmail(userUpdate.Email) {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		ctx.Write(models.MarshalErrorSt("This email is already exist"))
		return
	}

	res, err := rep.UpdateUser(nickname, userUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find user"))
			return
		}
		// log.Println(err)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
	return
}
