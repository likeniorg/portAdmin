package main

type verifyHashCode struct {
	// 身份标识
	Identify string

	// 绑定身份Hash
	BindHash string

	// 验证成功Hash
	VerifyHash string

	// vip过期
	VipLapse string

	// 帐号被绑定
	AccountBound string

	// 身份验证失败
	ErrorHash string

	// 没有绑定客户端
	NoBindClient string
}
