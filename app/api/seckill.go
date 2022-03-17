package api

import (
	"context"
	"fmt"
	"gf-seckill/app/dao"
	"gf-seckill/app/model"
	"gf-seckill/redis"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"net/http"
)

type apiSecKill struct{}

var SecKill = new(apiSecKill)

// SecKillDB 数据库秒杀
func (*apiSecKill) SecKillDB(r *ghttp.Request) {
	userId := r.GetInt("user_id")
	if userId == 0 {
		FailService(r, "用户不能为空")
	}

	product, err := dao.Product.FindOne(1)
	if err != nil {
		Fail(r, err)
	}

	if product.Status != 1 {
		FailService(r, "无当前活动")
	}

	if gtime.Now().Before(product.StartAt) {
		FailService(r, "活动未开始")
	}

	if gtime.Now().After(product.EndAt) {
		FailService(r, "活动已结束")
	}

	if product.ProductNum < 1 {
		FailService(r, "券码已经卖完")
	}

	userProductRef, err := dao.UserProductRef.FindOne(g.Map{"user_id": userId})
	if err != nil {
		Fail(r, err)
	}

	if userProductRef != nil {
		FailService(r, "当前用户已经抢到")
	}

	if _, err = g.DB().Exec("update product set product_num = product_num -1 where id = ?", product.Id); err != nil {
		Fail(r, err)
	}

	if _, err = dao.UserProductRef.Insert(&model.UserProductRef{
		UserId:    userId,
		ProductId: product.Id,
	}); err != nil {
		Fail(r, err)
	}
	Success(r)
}

// SecKillRedis Redis秒杀
func (*apiSecKill) SecKillRedis(r *ghttp.Request) {
	var ctx = context.Background()
	userId := r.GetInt("user_id")
	if userId == 0 {
		FailService(r, "用户不能为空")
	}
	result, err := redis.Client.HGetAll(ctx, "product_01").Result()
	if err != nil {
		glog.Errorf("client.HGetAll:err=%s \n", err.Error())
		Fail(r, err)
	}

	product := new(model.Product)

	if err = gconv.Struct(result, &product); err != nil {
		glog.Errorf("gconv.Struct:err=%s \n", err.Error())
		Fail(r, err)
	}

	if product.Status != 1 {
		FailService(r, "无当前活动")
	}

	if gtime.Now().Before(product.StartAt) {
		FailService(r, "活动未开始")
	}

	if gtime.Now().After(product.EndAt) {
		FailService(r, "活动已结束")
	}

	if product.ProductNum <= 0 {
		FailService(r, "券码已经卖完")
	}

	// 使用lua脚本来保证查询库存和减少库存的时是原子性操作
	secKillScript := `
		if redis.call("hexists", KEYS[1], KEYS[2]) == 1 then
			local stock = tonumber(redis.call("hget", KEYS[1], KEYS[2]))
			if stock > 0 then
				redis.call("hincrby", KEYS[1], KEYS[2], -1)
				return stock
			end
				return 0
		end
	`
	script := goRedis.NewScript(secKillScript)

	productNum, err := script.Run(ctx, redis.Client, []string{"product_01", "product_num"}).Int()
	if err != nil {
		glog.Errorf("script.Run:err=%s \n", err.Error())
		Fail(r, err)
	}

	// 这边再次获取库存是否满足足够
	if productNum <= 0 {
		FailService(r, "券码已经卖完")
	}

	// 使用lua脚本来实现一人一单
	userKillScript := `
		if redis.call("exists", KEYS[1]) == 1 then
			return 1
		else
			if redis.call("set", KEYS[1], "1", "ex", KEYS[2]) == "OK" then
				return 0
			else
				return 0
			end
		end
	`
	userScript := goRedis.NewScript(userKillScript)
	key1 := fmt.Sprintf("product_%d_user_%d", product.Id, userId)
	isExist, err := userScript.Run(ctx, redis.Client, []string{key1, "1000"}).Int()
	if isExist == 1 {
		glog.Warningf("用户[%d]已经抢过了... \n", userId)
		FailService(r, "当前用户已经抢过了")
	}

	if _, err = dao.UserProductRef.Insert(&model.UserProductRef{
		UserId:    userId,
		ProductId: 1,
	}); err != nil {
		glog.Errorf("dao.UserProductRef:err=%s \n", err.Error())
		Fail(r, err)
	}

	r.Response.WriteJsonExit("ok")
}

type Response struct {
	Code int
	Msg  string
	Data interface{}
}

func Success(r *ghttp.Request) {
	r.Response.WriteJsonExit(Response{
		Code: http.StatusOK,
		Msg:  http.StatusText(http.StatusOK),
		Data: nil,
	})
}

func Fail(r *ghttp.Request, err error) {
	r.Response.WriteJsonExit(Response{
		Code: http.StatusInternalServerError,
		Msg:  err.Error(),
		Data: nil,
	})
}

func FailService(r *ghttp.Request, msg string) {
	r.Response.WriteJsonExit(Response{
		Code: http.StatusInternalServerError,
		Msg:  msg,
		Data: nil,
	})
}
