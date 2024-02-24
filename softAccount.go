package main

import (
	"fmt"
	"os"
	"portAdmin/internal/cfg"
	"strconv"
	"time"
)

func addSoftEtherUser(hubName, username, userPassword string) error {
	time.Sleep(1 * time.Second)
	// 构造vpncmd命令
	outlinuxShell("./vpnserver/vpncmd  localhost:5555 /server  /HUB:" + cfg.HubName + "  /password:" + cfg.HubPass + " /cmd UserCreate " + username + " /group:none /realname:none  /note:none")
	outlinuxShell("./vpnserver/vpncmd  localhost:5555 /server  /HUB:" + cfg.HubName + "  /password:" + cfg.HubPass + " /cmd userpasswordset " + username + " /password:" + userPassword)
	fmt.Println()
	return nil
}

func softAccountWriten() {
	accounts, err := getInfos()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 批量添加用户
	for i, v := range accounts {
		err := addSoftEtherUser(cfg.HubName, v.Name, v.Pass)
		if err != nil {
			fmt.Println(err)
			fmt.Println("写入失败帐号序号为" + strconv.Itoa(i))
			return
		} else {
			fmt.Println("成功写入编号为" + strconv.Itoa(i) + "的帐号")
		}
	}
}
