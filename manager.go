package main

import (
	"encoding/json"
	"portAdmin/internal/account"
	"portAdmin/internal/cfg"
	"portAdmin/internal/util"

	"os"
	"time"
)

// 管理用户帐号
type Manager struct {
	Accounts []account.Account
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
func (m *Manager) MakeAccount() {
	acc := account.Account{
		Id:            util.RandString(),
		HostID:        "",
		Name:          util.RandString(),
		Pass:          util.RandString(),
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		LastLoginTime: time.Time{},
		Activate:      "false",
		Grade:         "会员",
		Note:          "",
	}
	m.Accounts = append(m.Accounts, acc)
}

// 生成多个帐号信息
func (m *Manager) MakeAccounts(number int) {
	for i := 0; i < number; i++ {
		m.MakeAccount()
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
