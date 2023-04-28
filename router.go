package main

import (
	"fmt"
	"github.com/Serendipity-sw/gutil"
	"github.com/gin-gonic/gin"
	"github.com/swgloomy/gutil/glog"
	"net/http"
)

func setGinRouter(r *gin.Engine) {
	g := &r.RouterGroup
	{
		g.GET("/", func(c *gin.Context) { c.File(fmt.Sprintf("%s/index.html", webTemplate)) })
		g.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, "ok") }) //确认接口服务程序是否健在
		g.GET("/ws", func(c *gin.Context) {
			//websocket.WsHandler(c.Writer, c.Request)
		})
	}
	r.NoRoute(noRoute)
}

func noRoute(c *gin.Context) {
	urlPath := c.Request.URL.Path
	if urlPath != "" {
		filePath := fmt.Sprintf("%s%s", webTemplate, urlPath)
		bo, err := gutil.PathExists(filePath)
		if err != nil {
			glog.Error("main noRoute path exists run err! urlPath: %s filePath: %s err: %+v \n", urlPath, filePath, err)
			return
		}
		if bo {
			c.File(filePath)
		}
	}
}
