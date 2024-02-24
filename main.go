package main

import (
	"encoding/json"
	"fmt"
	"portAdmin/internal/account"
	"portAdmin/internal/cfg"
	"portAdmin/internal/util"
	"time"

	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func init() {
	// 用户信息文件是否存在，不存在则新建
	_, err := os.Stat(cfg.AccountFile)
	if err != nil {
		os.Mkdir("conf", 0700)
		is := util.InputMsg("是否需要导入配置文件(y/n)")
		if is == "y" {
			fmt.Println("请将配置文件移动到conf目录")
			os.Exit(0)
		} else {
			fmt.Println("不导入配置文件，正在新建......")
		}
		if os.IsNotExist(err) {
			err := os.WriteFile(cfg.AccountFile, []byte(`null`), 0600)
			if err != nil {
				util.EchoExit(err)
			}
			fmt.Println("创建用户信息文件成功")
		}
	}

	// 检测自签证书存在
	_, err = os.Stat(cfg.SSTPCrtFile)
	if err != nil {
		os.Mkdir("conf/cert", 0700)
		// 创建sstp服务端证书
		linuxShell(`openssl req -x509 -nodes -days 365 -newkey rsa:3072 -keyout conf/cert/sstp.key -out conf/cert/sstp.crt -subj "/C=US/ST=California/L=SanFrancisco/O=likeniorg/OU=a/CN=` + cfg.Domain + `"`)
		fmt.Println("创建SSTP证书成功")
	}

	// 检测身份验证状态字符串是否存在
	_, err = os.Stat(cfg.VerifyHashCodeFile)
	if err != nil {
		verifyCode.AccountBound = util.RandStringIdent()
		verifyCode.BindHash = util.RandStringIdent()
		verifyCode.ErrorHash = util.RandStringIdent()
		verifyCode.Identify = util.RandStringIdent()
		verifyCode.VerifyHash = util.RandStringIdent()
		verifyCode.VipLapse = util.RandStringIdent()
		verifyCode.NoBindClient = util.RandStringIdent()
		data, err := json.MarshalIndent(verifyCode, "", "	")
		util.EchoExit(err)

		err = os.WriteFile(cfg.VerifyHashCodeFile, data, 0600)
		util.EchoExit(err)
		fmt.Println("创建C/S验证码成功")

	} else {
		data, err := os.ReadFile(cfg.VerifyHashCodeFile)
		util.EchoExit(err)
		err = json.Unmarshal(data, &verifyCode)
		util.EchoExit(err)
	}

	// 创建服务端程序
	if _, err := os.Stat("server"); err != nil {
		err := NewServerFile()
		util.EchoExit(err)
		fmt.Println("创建服务端程序成功")
	}

	// 默认生成10个帐号密码及客户端
	if _, err := os.Stat("tmp/" + cfg.SoftwareName + "0.go"); err != nil {
		fmt.Println("正在创建客户端程序，请稍后")
		acc := account.Manager{}
		acc.MakeAccounts(10)
		acc.SaveAccount()
		NewWinFile(0, 10)
		fmt.Println("创建客户端程序成功")
	}
}

var verifyCode verifyHashCode

// 后续可以升级为加密通信(使用证书客户端验证服务端)
func main() {
	fmt.Println("lno加速器管理程序启动成功")
	fmt.Println("选项：")
	fmt.Println("		1. 生成服务器程序")
	fmt.Println("		2. 创建连接客户端(生成帐号后才能创建连接客户端)")
	fmt.Println("		3. 生成帐号")
	fmt.Println("		4. 将生成帐号写入到softEther中")
	fmt.Println("		5. 输出帐号密码")
	fmt.Println("		6. 帐号状态重置为未绑定状态")
	fmt.Println("		0. 退出程序")

	// 循环等待执行条件
	for {
		switch input("\n请输入操作序号") {
		// 编译服务端程序
		case 1:
			NewServerFile()

		// 编译客户端程序
		case 2:
			uc, err := getInfos()
			if uc == nil {
				fmt.Println("需要生成帐号后才能创建客户端")
				os.Exit(0)
			}
			util.EchoError(err)
			length := len(uc)
			fmt.Println("请输入对应账号索引")
			for i := 0; i < length; i++ {
				fmt.Println(strconv.Itoa(i) + " " + uc[i].Id + uc[i].Note)
			}
			startIndex := input("输入开始索引")
			endIndex := input("输入结束索引")
			err = NewWinFile(startIndex, endIndex)
			util.EchoError(err)

		//生成帐号
		case 3:
			manager := Manager{}
			manager.GetAccount()

			index := input("输入生成帐号数量")
			manager.MakeAccounts(index)
			err := manager.SaveAccount()
			util.EchoError(err)

		// 生成的帐号信息写入到softether中
		case 4:
			softAccountWriten()

		// 输出帐号信息
		case 5:
			uc, err := getInfos()
			util.EchoError(err)
			length := len(uc)
			fmt.Println("请输入对应账号索引")
			for i := 0; i < length; i++ {
				fmt.Println(strconv.Itoa(i) + " " + uc[i].Name + "		" + uc[i].Pass + "	" + uc[i].Note + "	" + uc[i].Activate + " " + uc[i].HostID)
			}

		// 帐号绑定状态重置
		case 6:
			acc := account.Manager{}
			acc.GetAccount()
			length := len(acc.Accounts)
			fmt.Println("请输入对应账号索引")
			for i := 0; i < length; i++ {
				fmt.Println(strconv.Itoa(i) + " " + acc.Accounts[i].Name + "		" + acc.Accounts[i].Pass + acc.Accounts[i].Note + "	" + acc.Accounts[i].Activate + " " + acc.Accounts[i].HostID)
			}

			startI := input("开始索引")
			endI := input("结束索引")
			if endI >= len(acc.Accounts) {
				fmt.Println("超出索引范围,重新选择操作项目")
				return
			}
			for ; startI <= endI; startI++ {
				acc.Accounts[startI].Activate = "false"
				acc.Accounts[startI].EndTime = time.Time{}
				acc.Accounts[startI].StartTime = time.Time{}
				acc.Accounts[startI].LastLoginTime = time.Time{}
				acc.Accounts[startI].HostID = ""
			}

			acc.SaveAccount()
		case 0:
			fmt.Println("退出程序成功")
			os.Exit(0)

		default:
			fmt.Println("请按序号进行操作")

		}
	}
}

// 输出条件消息，返回一个索引
func input(outMsg string) int {
	fmt.Println(outMsg)
	index := 123
	fmt.Scanln(&index)
	return index
}

// 获取用户信息
func getInfos() ([]account.Account, error) {
	u := []account.Account{}
	data, err := os.ReadFile(cfg.AccountFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return nil, err

	}
	return u, err
}

func NewServerFile() error {
	os.Mkdir("tmp", 0700)
	data, err := os.ReadFile("template/server.go")
	util.EchoError(err)
	file, err := os.OpenFile("tmp/server.go", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	file.Write(data)
	file.WriteString(`var accountBound = "` + verifyCode.AccountBound + `"` + "\n")
	file.WriteString(`var bindHash = "` + verifyCode.BindHash + `"` + "\n")
	file.WriteString(`var errorHash = "` + verifyCode.ErrorHash + `"` + "\n")
	file.WriteString(`var identify = "` + verifyCode.Identify + `"` + "\n")
	file.WriteString(`var verifyHash = "` + verifyCode.VerifyHash + `"` + "\n")
	file.WriteString(`var vipLapse = "` + verifyCode.VipLapse + `"` + "\n")
	file.WriteString(`var noBindClient = "` + verifyCode.NoBindClient + `"` + "\n")

	file.Close()
	linuxShell("go build tmp/server.go")
	// linuxShell("rm -f tmp/server.go")
	return nil
}

// 创建windows可执行文件
func NewWinFile(start, end int) error {
	uc, err := getInfos()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join("template", "client.go"))
	if err != nil {
		return err
	}
	for ; start < end; start++ {
		currGoFileName := cfg.SoftwareName + strconv.Itoa(start) + ".go"
		currGoexeName := cfg.SoftwareName + strconv.Itoa(start) + ".exe"

		newFileData := []byte{}
		copy(newFileData, data)
		file, err := os.OpenFile(filepath.Join("tmp", currGoFileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		file.Write(data)
		file.WriteString(`var id = "` + uc[start].Id + `"` + "\n")
		file.WriteString(`var username = "` + uc[start].Name + `"` + "\n")
		file.WriteString(`var password = "` + uc[start].Pass + `"` + "\n")
		file.WriteString(`var sstpCrt = ` + "`" + util.GetCertString(cfg.SSTPCrtFile) + "`" + "\n")
		file.WriteString(`var vpnName = "` + cfg.VpnName + `"` + "\n")
		file.WriteString(`var domain = "` + cfg.Domain + `"` + "\n")
		file.WriteString(`var domainAndPort = "` + cfg.DomainAndPort + `"` + "\n")
		file.WriteString(`var accountBound = "` + verifyCode.AccountBound + `"` + "\n")
		file.WriteString(`var bindHash = "` + verifyCode.BindHash + `"` + "\n")
		file.WriteString(`var errorHash = "` + verifyCode.ErrorHash + `"` + "\n")
		file.WriteString(`var identify = "` + verifyCode.Identify + `"` + "\n")
		file.WriteString(`var verifyHash = "` + verifyCode.VerifyHash + `"` + "\n")
		file.WriteString(`var vipLapse = "` + verifyCode.VipLapse + `"` + "\n")
		file.WriteString(`var NoBindClient = "` + verifyCode.NoBindClient + `"` + "\n")

		file.Close()

		linuxShell("env GOOS=windows GOARCH=amd64 go build " + filepath.Join("tmp", currGoFileName))
		linuxShell("mv " + currGoexeName + " " + filepath.Join("tmp", currGoexeName))
	}

	return nil
}

// Linux Shell
func linuxShell(dst string) error {
	cmd := exec.Command("bash", "-c", dst)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func outlinuxShell(dst string) error {
	cmd := exec.Command("bash", "-c", dst)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}
