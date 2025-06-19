# 如何编译



## 编译x86-64系统

​	sh build.sh 



ARM设备上编译

​       sh build.sh --target=linux-aarch64

## 创建 RPM 和 DEB 安装包

要创建 RPM 和 DEB 包，可以使用 `fpm` 工具来简化流程。

yum install rpm-build -y

### 安装 `fpm`

如果尚未安装 `fpm`，可以使用以下命令来安装：

```
bash


复制代码
# 安装 Ruby（用于安装 fpm）
sudo apt-get install ruby ruby-dev build-essential   # 在 Ubuntu 上
sudo yum install ruby ruby-devel make gcc            # 在 CentOS 上

# 安装 fpm
sudo gem install --no-document fpm
```



### 创建安装包

到源码目录的apps/device_mgmt目录，执行 sh build_package.sh或者执行下面的命令生成

### 创建 RPM 包（用于 CentOS）

在 CentOS 上生成 RPM 包的命令如下：

```
bash


复制代码
fpm -s dir -t rpm -n devicemgmt -v 1.0 \
--after-install <(echo "systemctl enable devicemgmt && systemctl start devicemgmt") \
--after-remove <(echo "systemctl stop devicemgmt && systemctl disable devicemgmt") \
--description "Device Management Service" \
--license "MIT" \
devicemgmt=/usr/local/bin/devicemgmt \
devicemgmt.service=/etc/systemd/system/devicemgmt.service \
config.json=/etc/devicemgmt/config.json
```

此命令将生成一个名为 `devicemgmt-1.0-1.x86_64.rpm` 的 RPM 包。
`--after-install` 和 `--after-remove` 脚本用于在安装后启用并启动服务，卸载后停止并禁用服务。

### 创建 DEB 包（用于 Ubuntu）

在 Ubuntu 上生成 DEB 包的命令如下：

```
bash


复制代码
fpm -s dir -t deb -n devicemgmt -v 1.0 \
--after-install <(echo "systemctl enable devicemgmt && systemctl start devicemgmt") \
--after-remove <(echo "systemctl stop devicemgmt && systemctl disable devicemgmt") \
--description "Device Management Service" \
--license "MIT" \
devicemgmt=/usr/local/bin/devicemgmt \
devicemgmt.service=/etc/systemd/system/devicemgmt.service \
config.json=/etc/devicemgmt/config.json
```

此命令将生成一个名为 `devicemgmt_1.0_amd64.deb` 的 DEB 包。

## 安装和测试包

- 在 CentOS 上安装生成的 RPM 包：

  ```
  bash
  
  
  复制代码
  sudo rpm -ivh devicemgmt-1.0-1.x86_64.rpm
  ```

- 在 Ubuntu 上安装生成的 DEB 包：

  ```
  bash
  
  
  复制代码
  sudo dpkg -i devicemgmt_1.0_amd64.deb
  ```

## 验证服务状态

安装完成后，验证服务是否正确启动，并检查是否能在退出后自动拉起：

```
bash


复制代码
# 启动服务（如果未自动启动）
sudo systemctl start devicemgmt

# 查看服务状态
sudo systemctl status devicemgmt

# 测试自动拉起
sudo pkill devicemgmt
sleep 5
sudo systemctl status devicemgmt  # 确认服务重新启动
```

##  卸载服务

如需卸载服务，可以使用以下命令：

- 在 CentOS 上卸载：

  ```
  bash
  
  
  复制代码
  sudo rpm -e devicemgmt
  ```

- 在 Ubuntu 上卸载：

  ```
  bash
  
  
  复制代码
  sudo dpkg -r devicemgmt
  ```