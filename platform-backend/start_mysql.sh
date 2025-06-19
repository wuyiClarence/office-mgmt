#!/bin/bash

# MySQL 容器名称
CONTAINER_NAME="mysql"

# MySQL root 用户密码
MYSQL_ROOT_PASSWORD="123.com"

# 挂载的本地 MySQL 数据存储路径
MYSQL_DATA_PATH="/my/local/mysql/data"

mkdir -p $MYSQL_DATA_PATH

# 检查是否已经存在同名容器，如果存在先停止并删除
if [ "$(docker ps -a | grep $CONTAINER_NAME)" ]; then
    echo "停止并删除已有的 MySQL 容器..."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# 启动 MySQL 8.0 容器
echo "启动 MySQL 8.0 容器..."
docker run -d \
    --name $CONTAINER_NAME \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -v $MYSQL_DATA_PATH:/var/lib/mysql \
    -p 3306:3306 \
    mysql:8.0 \
    --default-authentication-plugin=mysql_native_password 

# 等待 MySQL 服务启动完成
echo "等待 MySQL 启动..."
sleep 20  # 可以根据实际情况调整等待时间

# 设置 root 用户允许远程访问
echo "设置 root 用户允许远程访问..."
docker exec -i $CONTAINER_NAME mysql -u root -p$MYSQL_ROOT_PASSWORD <<EOF
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '$MYSQL_ROOT_PASSWORD';
FLUSH PRIVILEGES;
EOF

# 重启 MySQL 容器
echo "重启 MySQL 容器..."
docker restart $CONTAINER_NAME

echo "MySQL 容器启动并已配置完成。"
