package main

import (
	_ "gf-seckill/boot"
	_ "gf-seckill/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
