# 如何部署

首先创建一个目录，将deploy，platform-ui, platform-backend, platform-mdns四个工程放在这个目录下面。

然后切换到deploy目录，直接执行 ./deploy.sh



## 安装可能遇到的问题

项目依赖docker-compose，需要确定安装的主机是否支持docker-compse命令。如果不支持则去网上下载

```
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

要安装其他版本的 Compose，请替换 1.24.1。

将可执行权限应用于二进制文件：

```
$ sudo chmod +x /usr/local/bin/docker-compose
```

创建软链：

```
$ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```

测试是否安装成功：

```
$ docker-compose --version
cker-compose version 1.24.1, build 4667896b
```

docker镜像的拉取可能失败，需要手动设置代理

1. 创建 dockerd 相关的 systemd 目录，这个目录下的配置将覆盖 dockerd 的默认配置

```shell
$ sudo mkdir -p /etc/systemd/system/docker.service.d
```

1. 新建配置文件 `/etc/systemd/system/docker.service.d/http-proxy.conf`，这个文件中将包含环境变量

```ini
[Service]
Environment="HTTP_PROXY=http://proxy.example.com"
Environment="HTTPS_PROXY=http://proxy.example.com"
```

1. 如果你自己建了私有的镜像仓库，需要 dockerd 绕过代理服务器直连，那么配置 `NO_PROXY` 变量：

```ini
[Service]
Environment="HTTP_PROXY=http://proxy.example.com"
Environment="HTTPS_PROXY=http://proxy.example.com"
Environment="NO_PROXY=your-registry.com,10.10.10.10,*.example.com"
```

多个 `NO_PROXY` 变量的值用逗号分隔，而且可以使用通配符（*），极端情况下，如果 `NO_PROXY=*`，那么所有请求都将不通过代理服务器。

1. 重新加载配置文件，重启 dockerd

```shell
$ sudo systemctl daemon-reload
$ sudo systemctl restart docker
```

1. 检查确认环境变量已经正确配置：

```shell
$ sudo systemctl show --property=Environment docker
```



docker-compose拉取镜像失败，可以尝试使用docker pull xxx 拉取镜像，拉取成功后，再执行./deploy.sh