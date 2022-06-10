## Install

```
curl -fsSL https://raw.githubusercontent.com/GeeScot/postgresql-check/main/install.sh | sudo bash
```

## Configure postgresql-check

```
sudo nano /lib/systemd/system/postgresql-check.service
```

Available Environment Variable Options

- PGUSER (default: postgres)
- PGPASS (default: no-password)
- PGHOST (default: localhost)
- PGPORT (default: 5433)

## Restart

```
sudo systemctl daemon-reload
```

```
sudo systemctl restart postgresql-check
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

## Test

```
curl -v http://localhost:26726
```

## Load Testing (macOS)

```
brew install k6
```

```
k6 run --vus 10 --duration 30s script.js
```
