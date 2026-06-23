`Goph` 是一款封装简化了官方 `golang.org/x/crypto/ssh` 的 `Go` 第三方 `SSH` `工具库，GitHub` 开源，主打极简 `API`、开箱即用，降低原生 `SSH` 包复杂的连接、会话、命令执行、文件传输成本。

## 安装

```bash
go get github.com/melbahja/goph
```

## 快速开始

使用账号密码登陆执行 `ls` 命令，并将结果打印出来。

```golang
package main

import (
	"log"

	"github.com/melbahja/goph"
)

func main() {

	// 创建 ssh 客户端i
	client, err := goph.New(
		// ssh 用户名
		"root",
		// ssh 服务器地址
		"10.160.162.42",
		// 使用密码方式鉴权
		goph.Password("dev123"))

	if err != nil {
		log.Fatal(err)
	}
	// 使用defer 关闭网络连接
	defer client.Close()

	// 执行 ls 命令
	out, err := client.Run("ls")
	if err != nil {
		log.Fatal(err)
	}

	// 打印执行结果
	print(string(out))
}

```

## 鉴权

`goph` 支持两种鉴权方式，密码和密钥文件。对应鉴权对象 `goph.Auth` 类型（`ssh.AuthMethod` 切片）。

使用将创建好的 `auth` 传入即可：

```golang
    client, err := goph.New("root", "192.1.1.3", auth)
    ...

    // 或者
    client, err = goph.NewConn(&goph.Config{
		User:     user,
		Addr:     addr,
		Port:     port,
		Auth:     auth,
		Callback: VerifyHost,
	})
```

### 密码

```golang
    goph.Password("密码")
```

### 密钥文件

密钥文件鉴权支持密码加密文件，将密码传入 `you_passphrase_here` 位置既可以。如果使用的是为加密的密钥文件，则 `you_passphrase_here` 处填入空字符串即可。

```golang
    auth, err := goph.Key("/home/mohamed/.ssh/id_rsa", "you_passphrase_here")
    if err != nil {
        // handle error
    }
```

## 创建客户端

常用的有两种，都能得到同一个类型的客户端对象：

> 这里仅介绍的方法仅仅允许连接已授信的 主机,当然也可以使用 `NewUnknown`，或配置 `ssh.InsecureIgnoreHostKey` 实现允许访问未授信的主机。但是这里并不推荐这样做。

### 极简模式

除了用户名、主机地址、鉴权信息外，其他配置均采用默认值。

- 端口：22
- 超时时间：20秒
- 主机验证：仅允许已授信主机

```golang
    ...
    client, err := goph.New("root", "192.1.1.3", goph.Password("you_password_here"))
    ...
```

### 自定配置

使用 `goph.NewConn` 可获得自定义空间，除了之前提到过的端口外，我们还可以自定义主机地址校验，剔除不希望访问的地址。

当然，如果没有特殊需求，主机验证方法可直接使用 `goph.DefaultKnownHosts` 或 `ssh.InsecureIgnoreHostKey`（不推荐）。

```golang
    ...
    client, err = goph.NewConn(&goph.Config{
        User:     user,
        Addr:     addr,
        Port:     port,
        Auth:     auth,
        Callback: VerifyHost,
    })
    ...
```

## 执行命令

```golang
	// 普通 返回字节切片
	out, err = client.Run("ls")
	// 带context的，可以进行超时控制、主动停止 返回字节切片
	out, err := client.RunContext(ctx, cmd)
	// command 模式 返回 *Cmd 对象
	cmd, err := client.Command(name, args...)
	// 带 context 的 command 模式 返回 *Cmd 对象
	cmd, err := client.CommandContext(ctx, name，args...)
```

## 文件相关

### 上传/下载

```golang
	// 上传 本地 到 远端
	err := client.Upload("/path/to/local/file", "/path/to/remote/file")
	// 下载 远端 到 本地
	err := client.Download("/path/to/remote/file", "/path/to/local/file")
```

### 文件操作

`goph` 不直接提供文件操作能力，而是提供了 `ssh` 客户端到 `sftp` 客户端的转化方法，通过 `sftp` 库实现文件的操作，包括读写、修改权限等等。

```golang
	...
	ftp, err := client.NewSftp()
	...

	// mkdir -p 创建文件夹
	err = ftp.MkdirAll("/path/to/remote/file")
```

`sftp` 提供的一些其他常用功能：

```golang
	func (c *sftp.Client) Chmod(path string, mode os.FileMode) error
	func (c *sftp.Client) Chown(path string, uid int, gid int) error
	func (c *sftp.Client) MkdirAll(path string) error
	func (c *sftp.Client) Remove(path string) error
	func (c *sftp.Client) Rename(oldname string, newname string) error
	func (c *sftp.Client) Truncate(path string, size int64) error
```

## 参考资料

[goph 仓库](github.com/melbahja/goph)
