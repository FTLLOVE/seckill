package main

import (
	"context"
	"fmt"
	"gf-seckill/app/dao"
	"gf-seckill/app/model"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"testing"
	"time"

	myRedis "gf-seckill/redis"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

//
// @auth: fangtongle
// @date: 2022/3/11 17:31
// @desc: lua_test.go
//

type Product struct {
	ProductName string
	ProductNum  int
	Status      int
	StartAt     gtime.Time
	EndAt       gtime.Time
}

func TestRedisOfLua(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
		DB:      0,
	})
	//t.Log(client.HMSet(context.Background(), "product_01", "product_name", "满100减50", "product_num", 100, "status", 1, "create_at", "2022-03-01", "update_at", "2022-05-01").Err())
	result, err := client.HGetAll(context.Background(), "product_01").Result()
	if err != nil {
		return
	}
	var p = new(model.Product)
	err = gconv.Struct(result, &p)
	if err != nil {
		return
	}
	t.Logf("%#v \n", p)
}

var wg sync.WaitGroup

func TestSecKill(t *testing.T) {

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			var (
				err    error
				lock   *redislock.Lock
				ctx    = context.Background()
				locker = redislock.New(myRedis.Client)
			)
			defer wg.Done()
			defer lock.Release(ctx)

			product := &Product{}
			productVal, _ := myRedis.Client.Get(ctx, "product_01").Result()
			_ = gconv.Struct(productVal, &product)
			expireTime := product.EndAt.FormatNew("Y-m-d").Sub(gtime.Now().FormatNew("Y-m-d")).Seconds()
			userKey := fmt.Sprintf("product_id_%d_user_id_%d", 1, 1)

			lock, err = locker.Obtain(ctx, userKey, time.Second*time.Duration(expireTime), nil)
			if err == redislock.ErrNotObtained {
				fmt.Println("locker.Obtain Could not obtain lock!")
			} else if err != nil {
				log.Fatalln("locker.Obtain", err)
			}
			_, _ = dao.UserProductRef.Insert(&model.UserProductRef{
				UserId:    1,
				ProductId: 1,
			})
		}()
	}
	wg.Wait()
}

func TestRedisExist(t *testing.T) {
	var ctx = context.Background()
	userKey := "user_id"
	t.Log(myRedis.Client.Exists(ctx, userKey).Val())

}
