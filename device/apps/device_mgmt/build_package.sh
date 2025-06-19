#!/bin/bash

rm -rf *.rpm
rm -rf *.deb

after_install_script=$(mktemp)
before_remove_script=$(mktemp)


echo "systemctl enable devicemgmt" > "$after_install_script"
echo "systemctl start devicemgmt" >> "$after_install_script"
echo "systemctl restart devicemgmt" >> "$after_install_script"  # 自动重启服务

echo "systemctl stop devicemgmt" > "$before_remove_script"
echo "systemctl disable devicemgmt" >> "$before_remove_script"

# 使用 fpm 创建 RPM 包
fpm -s dir -t rpm -n devicemgmt -v 1.0 \
--after-install "$after_install_script" \
--before-remove "$before_remove_script" \
--description "Device Management Service" \
--license "MIT" \
--config-files /etc/devicemgmt/config.json \
devicemgmt=/usr/local/bin/devicemgmt \
devicemgmt.service=/etc/systemd/system/devicemgmt.service \
config.json=/etc/devicemgmt/config.json

fpm -s dir -t deb -n devicemgmt -v 1.0 \
--deb-no-default-config-files \
--after-install "$after_install_script" \
--before-remove "$before_remove_script" \
--description "Device Management Service" \
--license "MIT" \
--config-files /etc/devicemgmt/config.json \
devicemgmt=/usr/local/bin/devicemgmt \
devicemgmt.service=/etc/systemd/system/devicemgmt.service \
config.json=/etc/devicemgmt/config.json

rm -f "$after_install_script" "$before_remove_script"