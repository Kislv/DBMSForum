package post

import (
	"DBMSForum/models"
	forumRep "DBMSForum/server/forum/rep"
	postRep "DBMSForum/server/post/rep"
	threadRep "DBMSForum/server/thread/rep"
	userRep "DBMSForum/server/user/rep"
	"database/sql"
	"strconv"
	"strings"

	json "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)

	var thread int
	var forum string
	var err error
	id, errInt := strconv.Atoi(slugOrID)

	if errInt != nil {
		thread, forum, err = threadRep.FindThreadIDForum(slugOrID)
	} else {
		thread, forum, err = threadRep.FindThreadByIDIDForum(id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find post thread"))
			return
		}
		// log.Println("1",err)
		return
	}

	posts := models.Posts{}
	err = json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		// log.Println("2",err)
		ctx.SetStatusCode(fasthttp.StatusCreated)
		ctx.Write([]byte("[]"))
		return
	}

	err = postRep.InsertPosts(&posts, thread, forum)

	if err != nil {
		if strings.Contains(err.Error(), "bad parent thread") {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(models.MarshalErrorSt("Parent post was created in another thread"))
			return
		}
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.MarshalErrorSt("Can't find post author by nickname"))
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(body)
}

func ThreadPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	query := ctx.QueryArgs()
	limit, err := strconv.Atoi(string(query.Peek("limit")))
	if err != nil {
		limit = 100
	}

	since, err := strconv.Atoi(string(query.Peek("since")))
	if err != nil {
		since = 0
	}
	sort := string(query.Peek("sort"))

	desc := query.GetBool("desc")
	//if err != nil {
	//	desc = false
	//}

	slugOrID := ctx.UserValue("slug_or_id").(string)

	var thread int
	id, errInt := strconv.Atoi(slugOrID)
	if errInt != nil {
		thread, err = threadRep.FindThreadID(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByIDID(id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find thread"))
			return
		}
		// log.Println(err)
		return
	}

	posts, err := postRep.FindPosts(thread, limit, since, sort, desc)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}

func Details(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		// log.Println(err)
		return
	}

	related := string(ctx.QueryArgs().Peek("related"))

	post, err := postRep.FindByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find post by id"))
			return
		}
		// log.Println(err)
		return
	}

	postFull := models.PostFull{}
	if strings.Contains(related, "user") {
		user, err := userRep.FindByNickname(post.Author)
		if err != nil {
			// log.Println(err)
			return
		}
		postFull.Author = user
	}

	if strings.Contains(related, "forum") {
		forum, err := forumRep.FindForum(post.Forum)
		if err != nil {
			// log.Println(err)
			return
		}
		postFull.Forum = forum
	}

	if strings.Contains(related, "thread") {
		thread, err := threadRep.FindThreadByID(post.Thread)
		if err != nil {
			// log.Println(err)
			return
		}
		postFull.Thread = thread
	}

	postFull.Post = post

	body, err := json.Marshal(postFull)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
	return

}

func DetailsPOST(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		// log.Println(err)
		return
	}

	post, err := postRep.FindByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find post by id"))
			return
		}
		// log.Println(err)
		return
	}

	postUpdate := models.PostUpdate{}
	err = json.Unmarshal(ctx.PostBody(), &postUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	err = postRep.UpdatePost(post, postUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(post)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}
