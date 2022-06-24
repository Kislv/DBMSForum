package thread

import (
	"DBMSForum/models"
	forumRep "DBMSForum/server/forum/rep"
	threadRep "DBMSForum/server/thread/rep"
	"database/sql"
	"strconv"
	"strings"

	json "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)

	thread := models.Thread{}
	err := json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		//log.Println(err)
		return
	}

	//if !userRep.CheckByNickname(thread.Author) {
	//	 ctx.SetStatusCode(fasthttp.StatusNotFound)
	//	ctx.Write(models.MarshalErrorSt("Can't find user"))
	//	return
	//}

	forum, err := forumRep.FindForum(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find thread forum"))
			return
		}
		//log.Println(err)
		return
	}

	thread.Forum = slug
	if thread.Slug != "" {
		thread2, err := threadRep.FindThread(thread.Slug)

		if err != sql.ErrNoRows {
			if err != nil {
				//log.Println(err)
				return
			}
			body, err := json.Marshal(thread2)
			if err != nil {
				//log.Println(err)
				return
			}

			ctx.SetStatusCode(fasthttp.StatusConflict)
			ctx.Write(body)
			return
		}
	}

	thread.Forum = forum.Slug
	err = threadRep.InsertThread(&thread)
	if err != nil {
		//log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.MarshalErrorSt("Can't find user"))
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		//log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(body)
}

func Vote(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	vote := models.Vote{}
	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		// log.Println(err)
		return
	}

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

	vote.Thread = thread
	err = threadRep.InsertVote(vote)
	if err != nil {
		// log.Println(err)
		if strings.Contains(err.Error(), "duplicate key") {
			err = threadRep.UpdateVote(vote)
		}
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.MarshalErrorSt("Can't find post author by nickname"))
			return
		}
	}

	threadUpdate, err := threadRep.FindThreadByID(thread)
	if err != nil {
		// log.Println(err)
		return
	}

	if errInt != nil && threadUpdate.Votes < 0 {
		threadUpdate.Votes *= -1
	}

	body, err := json.Marshal(threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}

func Details(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)

	var thread *models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	var err error
	if errInt != nil {
		thread, err = threadRep.FindThread(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByID(id)
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
	body, err := json.Marshal(thread)
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

	slugOrID := ctx.UserValue("slug_or_id").(string)

	var thread *models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	var err error
	if errInt != nil {
		thread, err = threadRep.FindThread(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByID(id)
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

	threadUpdate := models.ThreadUpdate{}
	err = json.Unmarshal(ctx.PostBody(), &threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	err = threadRep.UpdateThread(thread, threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		// log.Println(err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(body)
}
