package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/empirefox/bogger"
	"github.com/golang/glog"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/plugin/cors"
	"github.com/kataras/iris"
)

func main() {
	qiniu := bogger.NewQiniu(bogger.Config{
		Ak:              "ZGhoovpdfp8qNzeQmPFMjW0TWfDLkuJ47szA3pdD",
		Sk:              "oFpChuMomUSdOdxPQnGlyxH0nyWSuxC0GXoKcwD1",
		Bucket:          "dogger",
		UpLifeMinute:    1,
		MaxUpLifeMinute: 10,
		UpHost:          "http://upload.qiniu.com",
		UpHostSecure:    "https://up.qbox.me",
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
	app.Plugins.Add(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
	}))
	app.Use(logger.New())

	//	app.OnError(iris.StatusBadRequest, func(ctx *iris.Context) {
	//		ctx.Write("CUSTOM 404 NOT FOUND ERROR PAGE")
	//		ctx.Log("http status: 400 happened!")
	//	})

	s := &Server{
		Framework: app,
		qiniu:     qiniu,
	}

	s.Get("/qiniu/headtoken/:life", s.GetQiniuHeadToken)
	s.Get("/qiniu/uptoken/:key/:life", s.GetQiniuUptoken)
	s.Post("/qiniu/:prefix", s.PostQiniuList)
	s.Delete("/qiniu/:key", s.DeleteQiniu)

	return s
}

func (s *Server) GetQiniuHeadToken(ctx *iris.Context) {
	userId := 100
	life, _ := ctx.ParamInt("life")
	secure := strings.HasPrefix(ctx.RequestHeader("Origin"), "https://")
	ctx.JSON(iris.StatusOK, iris.Map{
		"Uptoken": s.qiniu.Uptoken(fmt.Sprintf("h/%d", userId), uint32(life), secure),
	})
}

func (s *Server) GetQiniuUptoken(ctx *iris.Context) {
	key, err := base64.URLEncoding.DecodeString(ctx.Param("key"))
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}
	life, _ := ctx.ParamInt("life")
	secure := strings.HasPrefix(ctx.RequestHeader("Origin"), "https://")
	ctx.JSON(iris.StatusOK, iris.Map{
		"Uptoken": s.qiniu.Uptoken(string(key), uint32(life), secure),
	})
}

func (s *Server) PostQiniuList(ctx *iris.Context) {
	prefix, err := base64.URLEncoding.DecodeString(ctx.Param("prefix"))
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	items, err := s.qiniu.List(string(prefix))
	if err != nil {
		glog.Errorln(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, items)

	glog.Warningln(string(prefix))
}

func (s *Server) DeleteQiniu(ctx *iris.Context) {
	key, err := base64.URLEncoding.DecodeString(ctx.Param("key"))
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	err = s.qiniu.Delete(string(key))
	if err != nil {
		glog.Errorln(err)
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, string(key))
	glog.Warningln(string(key))
}
