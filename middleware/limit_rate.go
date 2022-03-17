package middleware

import (
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/time/rate"
	"time"
)

//
// @auth: fangtongle
// @date: 2022/3/8 17:32
// @desc: limit_rate.go
//

type middleware struct {
}

var Middleware = new(middleware)

var limiter = rate.NewLimiter(rate.Every(time.Second*1), 5)

func (*middleware) LimitRate(r *ghttp.Request) {
	//if !limiter.Allow() {
	//	glog.Warning("=======被限制流量了........")
	//	r.Response.WriteStatusExit(http.StatusForbidden, http.StatusText(http.StatusForbidden))
	//}
	//
	//if len(r.Response.Header()) == 0 &&
	//	r.Response.Status == 0 &&
	//	r.Response.BufferLength() == 0 {
	//	r.Response.WriteJsonExit(g.Map{
	//		"code":    http.StatusNotFound,
	//		"message": "路由找不到",
	//		"data": struct {
	//		}{},
	//	})
	//}

	r.Middleware.Next()
}
