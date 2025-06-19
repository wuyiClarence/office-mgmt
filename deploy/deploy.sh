#!/bin/bash
set -euo pipefail

if [ -f ".env" ]; then
  docker-compose down
    echo ".env 文件已存在，跳过生成"
fi



INSTALL_DIR="${1:-$PWD/install}"
MYSQL_DATA_DIR="$INSTALL_DIR/mysql"
MQTT_DATA_DIR="$INSTALL_DIR/mqtt"
BACKEND_DATA_DIR="$INSTALL_DIR/backend"
FRONTEND_DATA_DIR="$INSTALL_DIR/frontend"
MDNS_DATA_DIR="$INSTALL_DIR/mdns"
WAKEONLAN_DATA_DIR="$INSTALL_DIR/wakeonlan"

mkdir -p "$MYSQL_DATA_DIR"
mkdir -p "$MQTT_DATA_DIR"
mkdir -p "$BACKEND_DATA_DIR"
mkdir -p "$FRONTEND_DATA_DIR"
mkdir -p "$MDNS_DATA_DIR"
mkdir -p "$WAKEONLAN_DATA_DIR"

export FRONTEND_DATA_DIR
export MQTT_DATA_DIR
export BACKEND_DATA_DIR
export MDNS_DATA_DIR
export WAKEONLAN_DATA_DIR


USER="officemgmt"
PASSWORD="officemgmtuser123"

# 其他数据库相关的环境变量
MYSQL_ROOT_PASSWORD="officemgmt123abc"
MYSQL_DATABASE="office-mgmt"
MYSQL_USER=$USER
MYSQL_PASSWORD=$PASSWORD

echo "使用以下数据库配置进行部署："
echo "MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD"
echo "MYSQL_DATABASE=$MYSQL_DATABASE"
echo "MYSQL_USER=$MYSQL_USER"
echo "MYSQL_PASSWORD=$MYSQL_PASSWORD"
echo "MYSQL_DATA_DIR=$MYSQL_DATA_DIR"

export MYSQL_ROOT_PASSWORD MYSQL_DATABASE MYSQL_USER MYSQL_PASSWORD MYSQL_DATA_DIR


mkdir -p $MQTT_DATA_DIR/config $MQTT_DATA_DIR/data $MQTT_DATA_DIR/log




# 创建空的密码文件
touch $MQTT_DATA_DIR/config/passwd

# 启动容器以便创建 Mosquitto 用户
docker-compose up -d mysql

echo "等待 MySQL 服务启动完成..."

# 检查 mysql 是否 ready
until docker exec officemgmt-mysql mysql -h127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD -e "SELECT 1;" &> /dev/null
do
    echo "MySQL 尚未就绪，等待 3 秒..."
    sleep 3
done
echo "MySQL 启动完成，配置远程访问权限..."

docker exec -i officemgmt-mysql mysql -h127.0.0.1 -u root -p$MYSQL_ROOT_PASSWORD <<EOF
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '$MYSQL_ROOT_PASSWORD';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
CREATE USER IF NOT EXISTS '$MYSQL_USER'@'%' IDENTIFIED WITH mysql_native_password BY '$MYSQL_PASSWORD';
GRANT ALL PRIVILEGES ON *.* TO '$MYSQL_USER'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EOF

echo "MySQL 远程访问配置完成"


cat > $MQTT_DATA_DIR/config/mosquitto.conf << EOF
# 监听端口
listener 1883

# 持久化设置
persistence true
persistence_location /mosquitto/data/

# 日志设置
log_dest file /mosquitto/log/mosquitto.log

# 禁止匿名访问
allow_anonymous false

# 密码文件路径
password_file /mosquitto/config/passwd
EOF

# 启动容器以便创建 Mosquitto 用户
docker-compose up -d mosquitto
# 等待 Mosquitto 服务启动
sleep 5

MOSQUITTO_USERNAME=$USER
MOSQUITTO_PASSWORD=$PASSWORD
docker exec officemgmt-mosquitto mosquitto_passwd -b /mosquitto/config/passwd $MOSQUITTO_USERNAME $MOSQUITTO_PASSWORD

# 重启 Mosquitto 服务以应用新的用户配置
docker-compose restart mosquitto


mkdir -p $BACKEND_DATA_DIR/config $BACKEND_DATA_DIR/runtime
cat > $BACKEND_DATA_DIR/config/config.yaml << EOF
App:
  Name: platform-backend
  Address: ":8080"
  Mod: debug
  LogExpire: 7
Mysql:
  Dsn: root:$MYSQL_ROOT_PASSWORD@tcp(officemgmt-mysql:3306)/office-mgmt?parseTime=true&rejectReadOnly=true
  MaxIdle: 3
  MaxOpen: 3
  Name:
  Debug: true
AutoMigration: true
CronGapTime: 10
MqttConfig:
  Broker: "tcp://officemgmt-mosquitto:1883"
  ClientID: "devicemgmt-server"
  UserName: "$MOSQUITTO_USERNAME"
  Password: "$MOSQUITTO_PASSWORD"
  KeepAlive: 20
  ConnectTimeOut: 100
  MaxReconnectInterval: 60
PasswordAuthKey: 973e53c7
EOF

cp ../platform-backend/config/menu_permission.json $BACKEND_DATA_DIR/config 


mkdir -p $MDNS_DATA_DIR/config $MDNS_DATA_DIR/runtime
cat > $MDNS_DATA_DIR/config/config.yaml << EOF
App:
  Name: platform-mdns
  Address: ":8088"
  Mod: debug
  LogExpire: 7
MqttConfig:
  Broker: ""
  UserName: "$MOSQUITTO_USERNAME"
  Password: "$MOSQUITTO_PASSWORD"
  Port: 21883
MdnsConfig:
  Domain: "local."
  Host: "officemgmt"
  Key: "zwxlink"
  ServiceName: "_officemgmt._tcp"
  Port: 80
EOF

mkdir -p $WAKEONLAN_DATA_DIR/config $WAKEONLAN_DATA_DIR/runtime
cat > $WAKEONLAN_DATA_DIR/config/config.yaml << EOF
App:
  Name: platform-wakeonlan
  Mod: debug
  LogExpire: 7
MqttConfig:
  Broker: "tcp://127.0.0.1:21883"
  UserName: "$MOSQUITTO_USERNAME"
  Password: "$MOSQUITTO_PASSWORD"
  Port: 21883
  KeepAlive: 20
  ConnectTimeOut: 100
  MaxReconnectInterval: 60
EOF

mkdir -p $FRONTEND_DATA_DIR/config $BACKEND_DATA_DIR/web

cat > $FRONTEND_DATA_DIR/config/nginx.conf << EOF
server {
    listen 80;
    server_name localhost;

    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files \$uri /index.html;
    }

    location /api/ {
        proxy_pass http://officemgmt-backend:8080/api/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

echo "构建前端 UI..."
docker-compose build frontend

cat > .env << EOF
MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD
MYSQL_DATABASE=$MYSQL_DATABASE
MYSQL_USER=$MYSQL_USER
MYSQL_PASSWORD=$MYSQL_PASSWORD
INSTALL_DIR=$INSTALL_DIR
MYSQL_DATA_DIR=$MYSQL_DATA_DIR
MQTT_DATA_DIR=$MQTT_DATA_DIR
BACKEND_DATA_DIR=$BACKEND_DATA_DIR
FRONTEND_DATA_DIR=$FRONTEND_DATA_DIR
MDNS_DATA_DIR=$MDNS_DATA_DIR
WAKEONLAN_DATA_DIR=$WAKEONLAN_DATA_DIR
EOF

docker-compose down
docker-compose build
docker-compose up -d