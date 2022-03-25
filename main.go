package main

import (
	"context"
	"dailyPic/config"
	"dailyPic/login"
	"dailyPic/submitForm"
	"fmt"
)


func main()  {
	run(context.Background())
}

func run(ctx context.Context){
	//获取配置信息
	info := config.GetConfig("./config.yaml")
	fmt.Println(info.Users)
	//登录获取cookies
	l := login.New(info.Users.Username,info.Users.Password)
	err := l.Login()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = l.Client.GetUserInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	//提交表单
	err = submitForm.SubmitForm(l, info.Users)
	if err != nil {
		fmt.Println(err)
		return
	}
}