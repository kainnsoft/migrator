package v1

import (
	"fmt"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/kainnsoft/migrator/internal/entity"
)

type httprouter struct {
	handler IHttpHandlers
	l       *zap.Logger
}

func NewRouter(mux *router.Router, handler IHttpHandlers, log *zap.Logger) {
	rout := httprouter{
		handler: handler,
		l:       log,
	}

	mux.GET("/ping", rout.Ping)
	mux.GET("/mysql-count", rout.MySqlCount)
	mux.GET("/pg-count", rout.PGCount)
	mux.GET("/mysql-list/{limit}", rout.MySqlList)
	mux.GET("/insert/{limit}", rout.Insert)
}

func (h httprouter) Ping(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("pong")
}

func (h httprouter) MySqlCount(ctx *fasthttp.RequestCtx) {
	var (
		count int
		err   error
	)
	if count, err = h.handler.GetCountAMoneyMoves(ctx); err != nil {
		fmt.Fprintf(ctx, "internal error: %s", err.Error())
	}
	fmt.Fprintf(ctx, "a_many_moves count is: %d", count)
}

func (h httprouter) MySqlList(ctx *fasthttp.RequestCtx) {
	var (
		limit int
		err   error
		res   []*entity.AMoneyMoves
	)
	if limit, err = strconv.Atoi(string(ctx.Request.URI().LastPathSegment())); err != nil {
		fmt.Fprintf(ctx, "wrong limit: %s", err.Error())
		return
	}
	if limit < 0 {
		fmt.Fprintf(ctx, "wrong limit: %d", limit)
		return
	}
	if res, err = h.handler.ListAMoneyMoves(ctx, limit); err != nil {
		fmt.Fprintf(ctx, "internal error: %s", err.Error())
	}
	for i, v := range res {
		fmt.Fprintf(ctx, "Item N %d is: %v\n", i, *v)
	}
}

func (h httprouter) PGCount(ctx *fasthttp.RequestCtx) {
	var (
		count int
		err   error
	)
	if count, err = h.handler.GetCountPayments(ctx); err != nil {
		fmt.Fprintf(ctx, "internal error: %s", err.Error())
	}
	fmt.Fprintf(ctx, "payments count is: %d", count)
}

func (h httprouter) Insert(ctx *fasthttp.RequestCtx) {
	var (
		limit int
		err   error
	)
	if limit, err = strconv.Atoi(string(ctx.Request.URI().LastPathSegment())); err != nil {
		fmt.Fprintf(ctx, "wrong limit: %s", err.Error())
		return
	}
	if limit < 0 {
		fmt.Fprintf(ctx, "wrong limit: %d", limit)
		return
	}
	if err = h.handler.DoTransferFromAMoneyMovesToAccountPayments(ctx, limit); err != nil {
		fmt.Fprintf(ctx, "transfer error: %s", err.Error())
		return
	}
	fmt.Fprintf(ctx, "transfer success")
}
