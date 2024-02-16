package cfg

// 服务器信息
var (
	// 域名
	// Domain string = "localhost"
	Domain string = ""

	// 服务器端口
	// ServerPort string = ":54321"
	ServerPort string = ""

	// 域名和端口
	DomainAndPort string = Domain + ServerPort

	// 用户帐号文件路径
	AccountFile string = "conf/uInfo.data"

	// SSTP证书文件路径，除非使用自己生成的证书(移动到相同路径且名字相同)，否则不要改动
	SSTPCrtFile string = "conf/cert/sstp.crt"
	SSTPKeyFile string = "conf/cert/sstp.key"
)

// 客户端信息
var (
	// 加速器软件名
	// SoftwareName string = "lnoPortAdmin"
	SoftwareName string = ""

	// vpn名字
	VpnName string = "lno"
)

// 通用信息
var (
	// 服务器客户端验证code
	VerifyHashCodeFile string = "conf/verifyCode.data"
)

// vpnInfo
var (
	// Hub密码
	HubPass = ""
	// Hub名称
	HubName = "vpn"
)
