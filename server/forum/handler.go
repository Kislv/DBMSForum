package forum

import (
	"database/sql"

	json "github.com/mailru/easyjson"

	//"encoding/json"
	"DBMSForum/models"
	forumRep "DBMSForum/server/forum/rep"
	threadRep "DBMSForum/server/thread/rep"
	userRep "DBMSForum/server/user/rep"
	"strconv"

	"github.com/valyala/fasthttp"
)

func ClearHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	err := forumRep.ClearForum()
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("null"))
}

func StatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	status := forumRep.StatusForum()
	body, err := json.Marshal(status)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}

func Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	forum := models.Forum{}
	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		return
	}

	user, err := userRep.CheckByNicknameR(forum.User)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find user"))
			return
		}
		// log.Println(err)
		return
	}

	forum.User = user

	conflictForum, err := forumRep.FindForum(forum.Slug)
	if err == sql.ErrNoRows {
		err = forumRep.InsertForum(forum)
		if err != nil {
			// log.Println(err)
			return
		}

		body, err := json.Marshal(forum)
		if err != nil {
			// log.Println(err)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		ctx.Write(body)

		return
	}
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(conflictForum)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusConflict)
	ctx.Write(body)
}

func Details(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)

	forum, err := forumRep.FindForum(slug)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find forum"))
			return
		}
		// log.Println(err)
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}

func Users(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)

	query := ctx.QueryArgs()
	limit, err := strconv.Atoi(string(query.Peek("limit")))
	if err != nil {
		limit = 100
	}

	since := string(query.Peek("since"))

	desc := query.GetBool("desc")
	//if err != nil {
	//	desc = false
	//}

	if !forumRep.CheckForum(slug) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.MarshalErrorSt("Can't find forum"))
		return
	}

	users, err := userRep.FindByForum(slug, since, limit, desc)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}

func Threads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	query := ctx.QueryArgs()
	limit, err := strconv.Atoi(string(query.Peek("limit")))
	if err != nil {
		limit = 100
	}

	since := string(query.Peek("since"))

	desc := query.GetBool("desc")
	//if err != nil {
	//	desc = false
	//}

	slug := ctx.UserValue("slug").(string)

	if !forumRep.CheckForum(slug) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.MarshalErrorSt("Can't find forum"))
		return
	}

	threads, err := threadRep.FindThreads(slug, since, limit, desc)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}
