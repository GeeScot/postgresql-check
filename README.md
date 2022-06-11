## Install

```console
curl -fsSL https://gee.dev/install/postgresql-check | sudo bash
```

## Configure postgresql-check

```console
sudo nano /etc/postgresql-check/config.json
```

```json
{
  "postgres": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "password"
  },
  "port": 26726
}
```

## Restart

```console
sudo systemctl daemon-reload
```

```console
sudo systemctl restart postgresql-check
```

## Run

```console
go build
```

```console
sudo systemctl start postgresql-check.service
```

## Test

```console
curl -v http://localhost:26726
```

## Load Testing (macOS)

```console
brew install k6
```

```console
k6 run --vus 10 --duration 30s script.js
```

## Configure HAProxy (/etc/haproxy/haproxy.cfg)

```
global
    maxconn 500

defaults
    log global
    mode tcp
    retries 2
    timeout client 30m
    timeout connect 4s
    timeout server 30m
    timeout check 5s

listen stats
    mode http
    bind *:7000
    stats enable
    stats uri /

listen postgresql
    bind *:6432
    mode tcp
    option httpchk
    http-check expect status 200
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server postgresql0 primary-server-ip:6432 check port 26726
    server postgresql1 replica-server-ip:6432 check port 26726

listen postgresql-readonly
    bind *:6433
    mode tcp
    option httpchk
    http-check expect status 206
    default-server inter 3s fall 3 rise 2 on-marked-down shutdown-sessions
    server postgresql0 primary-server-ip:6432 check port 26726
    server postgresql1 replica-server-ip:6432 check port 26726
```
