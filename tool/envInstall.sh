#!/bin/bash
# 安装go环境变量
./tool/goenv.sh

# 编译
go build .

# 执行默认操作
./tool/lnovpn

# 商业版专用搭建vpn脚本，社区版不提供
# ./vpnInstall.sh