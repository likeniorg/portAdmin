package account

import (
	"encoding/json"
	"lnovpn/internal/cfg"
	"lnovpn/internal/util"

	"os"
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
	// 备注
	Note string
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
}

// 管理用户帐号
type Manager struct {
	Accounts []Account
}

// 获取全部用户帐号信息
func (m *Manager) GetAccount() error {
	data, err := os.ReadFile(cfg.AccountFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &m.Accounts)
	if err != nil {
		return err
	}
	return nil
}

// 生成一个帐号信息
func (m *Manager) makeAccount() {
	Account := Account{
		Id:            util.RandString(),
		HostID:        "",
		Name:          util.RandString(),
		Pass:          util.RandString(),
		Note:          "",
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		LastLoginTime: time.Time{},
		Activate:      "false",
		Grade:         "会员",
	}
	m.Accounts = append(m.Accounts, Account)
}

// 生成多个帐号信息,切记需要SaveAccount保存信息
func (m *Manager) MakeAccounts(number int) {
	for i := 0; i < number; i++ {
		m.makeAccount()
	}
}

// 生成账号后需要保存用户信息
func (m *Manager) SaveAccount() error {
	data, err := json.MarshalIndent(m.Accounts, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.AccountFile, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
