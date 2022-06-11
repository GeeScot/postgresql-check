#!/bin/bash

# Stop the service if it is already running
sudo systemctl stop postgresql-check.service > /dev/null

# Install unzip silently (required for unarchiving download)
sudo apt install -y unzip > /dev/null

TARGET_FILE=postgresql-check-linux-amd64.zip

# Download the latest version
wget -q --show-progress --progress=bar:force:noscroll https://github.com/GeeScot/postgresql-check/releases/download/latest/$TARGET_FILE

# Unzip and remove zip file
unzip $TARGET_FILE
rm $TARGET_FILE

# Move binary to /usr/local/bin/ folder and fix permissions
sudo mv -f postgresql-check /usr/local/bin/
sudo chown root:root /usr/local/bin/postgresql-check

# File paths
APP_PATH=/etc/postgresql-check
CONFIG_FILE="${APP_PATH}/config.json"
SERVICE_FILE=/lib/systemd/system/postgresql-check.service

# Write config file if not exists
if [[ ! -f "$CONFIG_FILE" ]]; then
  sudo mkdir -p $APP_PATH
  sudo tee $CONFIG_FILE > /dev/null <<EOT
  {
    "postgres": {
      "host": "localhost",
      "port": 5432,
      "username": "postgres",
      "password": "password"
    },
    "port": 26726
  }
EOT
fi

# Write service file if not exists
if [[ -f "$SERVICE_FILE" ]]; then
    sudo systemctl start postgresql-check.service
    exit 0
fi

sudo tee $SERVICE_FILE > /dev/null <<EOT
[Unit]
Description=postgresql-check

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/local/bin/postgresql-check

[Install]
WantedBy=multi-user.target
EOT

sudo systemctl enable postgresql-check.service
sudo systemctl start postgresql-check.service
