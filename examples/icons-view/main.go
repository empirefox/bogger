package main

import (
	"github.com/empirefox/bogger"
	"github.com/golang/glog"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
)

func main() {
	qiniu := bogger.NewQiniu(bogger.Config{
		Ak:           "ZGhoovpdfp8qNzeQmPFMjW0TWfDLkuJ47szA3pdD",
		Sk:           "oFpChuMomUSdOdxPQnGlyxH0nyWSuxC0GXoKcwD1",
		Bucket:       "dogger",
		UpLifeMinute: 1,
		UpHost:       "https://up.qbox.me",
	})
	s := NewServer(qiniu)
	s.Listen(":9999")
}

type Server struct {
	*iris.Framework
	qiniu *bogger.Qiniu
}

func NewServer(qiniu *bogger.Qiniu) *Server {
	app := iris.New(iris.Configuration{
		IsDevelopment: true,
	})
	app.Use(logger.New())
	app.Use(cors.Default())

	//	app.OnError(iris.StatusBadRequest, func(ctx *iris.Context) {
	//		ctx.Write("CUSTOM 404 NOT FOUND ERROR PAGE")
	//		ctx.Log("http status: 400 happened!")
	//	})

	s := &Server{
		Framework: app,
		qiniu:     qiniu,
	}

	s.Get("/uptoken", s.GetUptoken)
	s.Post("/list", s.PostList)
	s.Post("/delete", s.PostDelete)

	return s
}

func (s *Server) GetUptoken(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, iris.Map{
		"Uptoken": s.qiniu.Uptoken(),
	})
}

func (s *Server) PostList(ctx *iris.Context) {
	var data struct{ Prefix string }
	if err := ctx.ReadJSON(&data); err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	items, err := s.qiniu.List(data.Prefix)
	if err != nil {
		glog.Errorln(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, items)
}

func (s *Server) PostDelete(ctx *iris.Context) {
	var data struct{ Key string }
	if err := ctx.ReadJSON(&data); err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	err := s.qiniu.Delete(data.Key)
	if err != nil {
		glog.Errorln(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
}
