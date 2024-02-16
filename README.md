# portAdmin简介
项目别名lno加速器，用于确认客户端是否符合登录条件及建立加速器连接  
适用范围: 电脑访问外网验证、非公开端口验证访问

## 设计理念
避免服务器大范围开放防火墙白名单、出差变更IP后不在白名单无法访问的弊端

## 协议选择
* PPTP速度快但安全性低，在移动运营商网络中无法登录，不采用
* l2tp/ipsec兼顾速度及安全性，但是服务端口采用UDP协议，部分公司防火墙会屏蔽该端口，导致无法连接，不采用
* ikev2安全性极高，但是需要安装证书以及其端口经常会被GFW屏蔽，不采用
* SSTP采用ssl加密，端口流量走443，不易被屏蔽，使用速度良好，故采用

## 使用效果
程序编译成功后启动./lnovpn程序，自动创建SSTP证书、C/S验证码、用户核心信息、服务端程序、多个客户端程序，执行服务端程序后将客户端程序下发，运行客户端向服务端请求允许当前IP访问指定端口

## 安全性
* 服务端采用非周知TCP端口  
* 客户端与服务端使用随机生成的消息来进行验证，验证通过才会核实帐号和密码，否则终止连接且不返回任何消息  
* 核实帐号和密码是否正确，正确的话会检测当前帐号状态，分别是：验证成功、有效期过期、帐号或密码错误  
* 验证通过后(保持当前程序运行)，放行客户端IP允许访问服务器指定端口
* 终止程序后客户端向服务端发送断开连接指令，服务端清除当前客户端允许访问的临时策略

# 使用教程
## 环境安装
openssl、go(1.21.6版本以上)、gcc、make
### debian环境安装教程
```bash
# 安装软件包
sudo apt install gcc make openssl wget
```
### centos环境安装教程
```bash
# 安装软件包
sudo yum install gcc make openssl wget
```
## 安装
```bash
./tool/envInstall
```

## 恢复初始化环境(删除生成的文件)
```bash
./tool/init.sh
```
## 后台运行程序
```bash
# 启动服务端
sudo nohub ./server &
```

# 付费搭建加速器环境(适用于个人及团体)
邮箱联系likeniorg@gmail.com