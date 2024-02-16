package main

import (
	"bufio"
	"errors"
	"fmt"
	"lnovpn/internal/account"
	"lnovpn/internal/cfg"

	"log/slog"
	"net"
	"os"
	"os/exec"
	"time"
)

// 允许防火墙通行的IP
var allowFireIP = make(map[string]string)

// 日志记录
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {

	listener, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("AcceptErr", err.Error(), "ip", conn.RemoteAddr())
			continue
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 获取客户端身份
	cliIdent := msgRead(conn)
	switch cliIdent {
	// 普通客户端连接
	case identify:
		defer conn.Close()
		// 读取客户端ID
		cID := msgRead(conn)
		// 读取客户端主机ID
		hostID := msgRead(conn)

		// 获取客户端IP
		cIP, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			logger.Error("getIPFail", conn.RemoteAddr())
			return
		}

		startTime, endTime, lastLoginTime, grade, err := verifyCliID(conn, cID, hostID)
		// 验证客户端
		if err != nil {
			// 如果同时存在信息说明验证通过但是会员过期
			if err.Error() == "daoqi" {
				msgWrite(conn, vipLapse)
				msgWrite(conn, grade)
				msgWrite(conn, startTime)
				msgWrite(conn, endTime)
				msgWrite(conn, lastLoginTime)
				return
			}

			if err.Error() == "notBindClientConn" {
				msgWrite(conn, accountBound)
				return
			}
			// 验证客户端身份错误
			msgWrite(conn, errorHash)
			return
		}
		msgWrite(conn, verifyHash)
		msgWrite(conn, grade)
		msgWrite(conn, startTime)
		msgWrite(conn, endTime)
		msgWrite(conn, lastLoginTime)

		// 开启服务器防火墙
		linuxShell(`firewall-cmd --add-rich-rule='rule family="ipv4" source address="` + cIP + `" port protocol="tcp" port="443" accept'`)
		allowFireIP[cID] = cIP

		// 等待客户端关闭连接
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		// 严格来说EOF错误才是客户端关闭连接
		if err != nil {
			exitNumber := 0
			for _, v := range allowFireIP {
				if v == conn.RemoteAddr().String() {
					exitNumber += 1
				}
			}
			if exitNumber != 1 {
				allowFireIP[cID] = ""
				linuxShell(`firewall-cmd --remove-rich-rule='rule family="ipv4" source address="` + conn.RemoteAddr().String() + `" port protocol="tcp" port="443" accept'`)
			}
		}

	// 非法访问
	default:
		logger.Error("errorIdentAccess", identify, "ip", conn.RemoteAddr())
		return
	}

}

// Linux Shell
func linuxShell(dst string) error {
	cmd := exec.Command("bash", "-c", dst)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func msgRead(conn net.Conn) string {
	time.Sleep(300 * time.Millisecond)
	input := bufio.NewScanner(conn)
	input.Scan()
	return input.Text()
}

func msgWrite(conn net.Conn, msg string) error {
	time.Sleep(300 * time.Millisecond)
	write := bufio.NewWriter(conn)
	_, err := write.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	err = write.Flush()
	return err
}

// 验证客户端ID并返回消息
func verifyCliID(conn net.Conn, id, hostID string) (startTime, endTime, lastLoginTime, grade string, verifyErr error) {
	// 获取最新的帐号信息
	m := account.Manager{}
	err := m.GetAccount()
	if err != nil {
		logger.Error("AllUserAccountErr", err)
		return "", "", "", "", errors.New("error")
	}

	for i := 0; i < len(m.Accounts); i++ {
		//是否存在这个ID
		if m.Accounts[i].Id == id {
			// 是否已激活
			if m.Accounts[i].StartTime.IsZero() {
				//是否绑定该机器
				msgWrite(conn, bindHash)
				isBind := msgRead(conn)
				// 不绑定就终止连接
				if isBind == "y" {
					m.Accounts[i].StartTime = time.Now()
					m.Accounts[i].EndTime = time.Now().AddDate(0, 0, 31)
					m.Accounts[i].LastLoginTime = time.Now()
					m.Accounts[i].HostID = hostID
					m.Accounts[i].Activate = "true"
					err := m.SaveAccount()
					if err != nil {
						return "", "", "", "", errors.New("error")
					}
					return m.Accounts[i].StartTime.String(), m.Accounts[i].EndTime.String(), m.Accounts[i].LastLoginTime.String(), m.Accounts[i].Grade, nil
				}

				return

			}

			// 判断客户端ID是否为绑定的唯一值
			if hostID == m.Accounts[i].HostID {
				// 是否会员到期
				if m.Accounts[i].StartTime.Before(m.Accounts[i].EndTime) {
					return m.Accounts[i].StartTime.String(), m.Accounts[i].EndTime.String(), m.Accounts[i].LastLoginTime.String(), m.Accounts[i].Grade, nil
				} else {
					logger.Error("daoqiVip", errors.New(m.Accounts[i].Name))
					return m.Accounts[i].StartTime.String(), m.Accounts[i].EndTime.String(), m.Accounts[i].LastLoginTime.String(), m.Accounts[i].Grade, errors.New("daoqi")
				}
			} else {
				return "", "", "", "", errors.New("notBindClientConn")
			}
		}
	}
	return "", "", "", "", errors.New("请联系管理员获取最新客户端")
}
