#!/bin/bash

# 设置一些变量
MOSQUITTO_CONF_DIR="/my/local/mosquitto/data/"
MOSQUITTO_PASS_FILE="$MOSQUITTO_CONF_DIR/passwordfile"
MOSQUITTO_CONF_FILE="$MOSQUITTO_CONF_DIR/mosquitto.conf"
DOCKER_IMAGE="eclipse-mosquitto"
CONTAINER_NAME="mosquitto"
USER="mqtt"
PASSWORD="123.com"

# 检查是否已经存在同名容器，如果存在先停止并删除
if [ "$(docker ps -a | grep $CONTAINER_NAME)" ]; then
    echo "停止并删除已有的 MySQL 容器..."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Step 2: 拉取Mosquitto镜像
echo "拉取Mosquitto Docker镜像..."
docker pull $DOCKER_IMAGE

# Step 3: 创建配置目录和配置文件
echo "创建Mosquitto配置目录和文件..."
mkdir -p $MOSQUITTO_CONF_DIR

# 创建mosquitto.conf文件并写入配置
cat > $MOSQUITTO_CONF_FILE <<EOL
listener 1883
allow_anonymous false
password_file /mosquitto/config/passwordfile
EOL

# Step 4: 创建用户名密码文件
echo "创建Mosquitto用户密码文件..."
touch $MOSQUITTO_PASS_FILE

# 使用mosquitto_passwd工具生成密码文件
if ! command -v mosquitto_passwd &> /dev/null
then
    echo "mosquitto_passwd未安装，正在安装..."
    sudo yum install -y mosquitto
fi

mosquitto_passwd -b $MOSQUITTO_PASS_FILE $USER $PASSWORD

# Step 5: 运行Mosquitto容器
echo "启动Mosquitto Docker容器..."
docker run -d --name $CONTAINER_NAME -p 1883:1883 \
  -v $MOSQUITTO_CONF_DIR:/mosquitto/config \
  $DOCKER_IMAGE

# 检查容器是否启动成功
if [ $? -eq 0 ]; then
    echo "Mosquitto已成功启动并配置了用户名密码认证"
else
    echo "Mosquitto启动失败，请检查日志"
fi
