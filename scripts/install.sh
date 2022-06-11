#!/bin/bash

sudo apt install unzip
wget https://github.com/GeeScot/postgresql-check/releases/download/latest/postgresql-check-linux-amd64.zip

unzip postgresql-check-linux-amd64.zip
rm postgresql-check-linux-amd64.zip

sudo mv postgresql-check /usr/local/bin/
sudo chown root:root /usr/local/bin/postgresql-check

sudo tee /lib/systemd/system/postgresql-check.service > /dev/null <<EOT
[Unit]
Description=postgresql-check

[Service]
Type=simple
Restart=always
RestartSec=5s
Environment="PGUSER=postgres"
Environment="PGPASS="
Environment="PGHOST=localhost"
Environment="PGPORT=5432"
ExecStart=/usr/local/bin/postgresql-check

[Install]
WantedBy=multi-user.target
EOT

sudo systemctl enable postgresql-check.service
sudo systemctl start postgresql-check.service
