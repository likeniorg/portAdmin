package main

import (
	"time"
)

// 用户信息
type Account struct {
	// 帐号ID
	Id string
	// 客户端主机ID
	HostID string
	// 用户名
	Name string
	//密码
	Pass string
	// 用户注册时间
	StartTime time.Time
	// 用户服务结束时间
	EndTime time.Time
	// 用户上次登录时间
	LastLoginTime time.Time
	// 是否激活
	Activate string
	// 用户等级
	Grade string
	// 备注
	Note string
}
