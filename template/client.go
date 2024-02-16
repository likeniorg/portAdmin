package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

func env() {

	// 检测是否带有参数并执行相应操作
	if len(os.Args) == 2 {

		// 配置环境
		if os.Args[1] == "1" {
			configEnv()
		}

	}

	// 检测环境，是否能运行vpn
	if err := envChecking(); err != nil {
		fmt.Println("检测到运行环境不匹配,需要管理员权限配置运行环境")
		adminAuth("1")
		for i := 8; i > 0; i-- {
			fmt.Println(strconv.Itoa(i) + "秒后尝试重新连接")
			time.Sleep(1 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}

}

// 配置环境
func configEnv() {
	WinShell("echo '" + sstpCrt + "' > ca.cert.pem")
	WinShell("certutil -addstore ROOT ca.cert.pem")
	WinShell("del ca.cert.pem")
	WinShell(`Add-VpnConnection -Name ` + vpnName + ` -ServerAddress ` + domain + ` -TunnelType SSTP -EncryptionLevel Required  -AuthenticationMethod Chap,MSChapv2 `)
	regWriten(`SYSTEM\CurrentControlSet\Services\PolicyAgent`, `AssumeUDPEncapsulationContextOnSendRule`, 2)
	regWriten(`SYSTEM\CurrentControlSet\Services\Rasman\Parameters`, `ProhibitIpSec`, 0)
	fmt.Println("*************************************************")
	fmt.Println("*	      请重启电脑使环境生效后连接加速器  	 *")
	fmt.Println("*	 如果运行过前几个版本的加速器不需要重启电脑  *")
	fmt.Println("*************************************************")
	fmt.Println("15秒后自动关闭此窗口")
	time.Sleep(15 * time.Second)
	os.Exit(0)
}

// 检测环境
func envChecking() error {
	// 计算证书的 SHA-1 指纹
	certsha := sstpCrtSha1()
	cmd := exec.Command("powershell", `Get-ChildItem -Path Cert:\LocalMachine\Root | Where-Object { $_.Thumbprint -eq '`+certsha+`'}`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// 输出不是空就是已有证书
	if string(out) == "" {
		return errors.New("1环境不匹配")
	}

	// 支持网络地址转换
	val, err := regVal(`SYSTEM\CurrentControlSet\Services\PolicyAgent`, `AssumeUDPEncapsulationContextOnSendRule`)
	if err != nil {
		return err
	}
	if val != 2 {
		return errors.New("2环境不匹配")
	}

	// 使用ipsec协议(检测是否需要这个配置)
	val1, err := regVal(`SYSTEM\CurrentControlSet\Services\Rasman\Parameters`, `ProhibitIpSec`)
	if err != nil {
		return err
	}
	if val1 != 0 {
		return errors.New("3环境不匹配")
	}

	_, err = WinShell("Get-VpnConnection " + vpnName)
	if err != nil {
		return errors.New("没有创建VPN")
	}
	return nil
}

// 获取证书指纹
func sstpCrtSha1() string {
	// 解码 PEM 编码的证书
	block, _ := pem.Decode([]byte(sstpCrt))
	if block == nil {
		fmt.Println("Error decoding PEM block")
	}

	// 解析 X.509 证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing certificate:", err)
	}
	certSha1 := sha1.Sum(cert.Raw)
	return hex.EncodeToString(certSha1[:])
}

// 获取windwos注册表键值对
func regVal(keyPath, valueName string) (value uint64, err error) {
	// 打开注册表键
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer key.Close()

	// 读取注册表值
	val, _, err := key.GetIntegerValue(valueName)
	if err != nil {
		return 0, err
	}
	return val, err
}

// 注册表写入
func regWriten(keyPath, valueName string, value int) error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetDWordValue(valueName, uint32(value))
	if err != nil {
		return err
	}
	return nil
}

func WinShell(dst string) (string, error) {
	cmd := exec.Command("powershell", dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(out), nil
}

func ReadMsg(conn net.Conn) string {
	time.Sleep(300 * time.Millisecond)
	input := bufio.NewScanner(conn)
	input.Scan()
	return input.Text()
}
func WriteMsg(conn net.Conn, msg string) error {
	time.Sleep(300 * time.Millisecond)
	write := bufio.NewWriter(conn)
	_, err := write.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	err = write.Flush()
	return err
}

// 程序以管理员权限执行
func adminAuth(parameter string) {
	// 获取当前文件路径
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径:", err)
		return
	}

	cmd := exec.Command("powershell", "-Command", "Start-Process", executable, parameter, "-Verb", "runAs")
	err = cmd.Run()
	if err != nil {
		fmt.Println("需要以管理员权限配置加速器环境!", err)
		return
	}
}

func main() {
	env()
	fmt.Println("欢迎使用lno加速器稳定版v2.0")
	fmt.Println("正在尝试连接服务器")
	conn, err := net.Dial("tcp", domainAndPort)
	if err != nil {
		fmt.Println("连接服务器失败，请重试")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	fmt.Println("连接服务器成功")

	defer conn.Close()

	// 发送身份标识
	WriteMsg(conn, identify)
	// 发送客户端ID
	WriteMsg(conn, id)

	// 发送电脑Uid
	uc, err := user.Current()
	if err != nil {
		fmt.Println("获取客户端用户名失败")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	WriteMsg(conn, uc.Uid)
	fmt.Println("正在验证客户端身份")

	// 对比是否通过服务端验证
	hash := ReadMsg(conn)
	if hash == bindHash {
		fmt.Println("******************************************")
		fmt.Println("每个帐号只能绑定一个电脑，是否绑定当前电脑(y/n)")
		fmt.Println("******************************************")
		forbreak := true
		for forbreak {
			isBind := ""
			fmt.Scanln(&isBind)
			switch isBind {
			case "y":
				WriteMsg(conn, "y")
				fmt.Println("绑定成功,正在激活账户(与中国有时差)")
				hash = verifyHash
				forbreak = false

			case "n":
				fmt.Println("已拒绝了绑定该电脑")
				fmt.Println("15秒后自动退出")
				time.Sleep(15 * time.Second)
				os.Exit(0)
			default:
				fmt.Println("输入错误，绑定电脑输入\"y\"并回车(避免自动激活该帐号)")
			}

		}

	}
	switch hash {
	// 验证成功
	case verifyHash:
		grade := ReadMsg(conn)
		startTime := ReadMsg(conn)
		endTime := ReadMsg(conn)
		lastLoginTime := ReadMsg(conn)
		fmt.Println()
		fmt.Println("***************************************************************************************************************")
		fmt.Println("			" + grade + "欢迎登录" + "			")
		fmt.Printf("			会员开始日期：%s\n", startTime)
		fmt.Printf("			会员到期日期：%s\n", endTime)
		fmt.Printf("			最后登录时间为%s\n", lastLoginTime)
		fmt.Println("***************************************************************************************************************")

		fmt.Println("使用过程中不要关闭此程序，否则会断开加速器连接")
		cmd := exec.Command("powershell", "rasdial "+vpnName+" "+username+" "+password)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		// 设置一个信号通道
		signalChannel := make(chan os.Signal, 1)

		// 监听Ctrl+C或关闭窗口事件
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		// 启动一个 goroutine 处理信号
		go func() {
			sig := <-signalChannel
			fmt.Printf("Received signal: %v\n", sig)
			conn.Close()
			// 在这里添加你的最后操作的代码
			cmd := exec.Command("powershell", `rasdial "`+vpnName+`" /disconnect`)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error disconnecting VPN: %v\n", err)
			}

			// 退出程序
			os.Exit(0)
		}()

		select {}

	// 身份验证成功但是会员过期
	case vipLapse:
		grade := ReadMsg(conn)
		startTime := ReadMsg(conn)
		endTime := ReadMsg(conn)
		lastLoginTime := ReadMsg(conn)

		fmt.Println(grade + "欢迎登录")
		fmt.Println()
		fmt.Println("******************************************")
		fmt.Println("会员已到期，请联系管理员续费")
		fmt.Println("******************************************")
		fmt.Println()
		fmt.Printf("会员开始日期：%s，会员到期日期：%s\n", startTime, endTime)
		fmt.Printf("最后登录时间为%s\n", lastLoginTime)
		time.Sleep(100 * time.Second)
		os.Exit(1)

	case accountBound:
		fmt.Printf("该帐号已被绑定到其他电脑")
		time.Sleep(100 * time.Second)
		os.Exit(1)

	case errorHash:
		fmt.Printf("身份验证失败，请联系管理员获取新的客户端程序")
		time.Sleep(100 * time.Second)
		os.Exit(1)

	default:
		fmt.Println("交易终止，有内鬼")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

}
