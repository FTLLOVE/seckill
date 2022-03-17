package main

import (
	"fmt"
	"gf-seckill/app/dao"
	"gf-seckill/app/model"
	"github.com/gogf/gf/frame/g"
	"testing"
)

//
// @auth: fangtongle
// @date: 2022/3/3 16:05
// @desc: product_test.go
//
func TestProductImportRedis(t *testing.T) {
	product, _ := dao.Product.FindOne(1)
	t.Log(g.Redis().Do("hset", "product_01", "product_name", product.ProductName, "product_num", product.ProductNum, "status", product.Status, "start_at", product.StartAt.String(), "end_at", product.EndAt.String()))
}

func TestImportUserToDB(t *testing.T) {
	userList := make([]*model.User, 0)
	for i := 0; i < 2000; i++ {
		user := &model.User{
			Username:  fmt.Sprintf("用户%04d", i),
			Password:  "123456",
			Telephone: fmt.Sprintf("1380000%04d", i),
		}
		userList = append(userList, user)
	}
	t.Log(dao.User.Insert(userList))
}

func TestHGetForStruct(t *testing.T) {
	doVar, err := g.Redis().DoVar("HGETALL", "product_01")
	if err != nil {
		return
	}
	product := &model.Product{}
	err = doVar.Struct(&product)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(product)
}

func TestLimitRate(t *testing.T) {

}
