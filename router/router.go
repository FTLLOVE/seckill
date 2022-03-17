package router

import (
	"gf-seckill/app/api"
	"gf-seckill/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()

	s.SetNameToUriType(ghttp.UriTypeCamel)

	s.BindMiddlewareDefault(middleware.Middleware.LimitRate)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/secKill", api.SecKill)
	})
}
