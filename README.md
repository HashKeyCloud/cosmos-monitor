# cosmos-monitor

The program used to cosmos-monitor the status of specific validators on the Cosmos Hub. 
The program will send an alarm by email in the following cases:
- If the validator is jailed, `receiver1` and `receiver2` will be notified by email;
- If the validator is inactive, `receiver1` and `receiver2` will be notified by email;
- If the validator has not signed five blocks in a row or the unsigned rate of the last `n` blocks reaches `m`, `receiver1` and `receiver2` will be notified by email;
- If there is a new proposal, the `receiver1` will be notified by email

> receiver1: receive all alert emails;
> 
> receiver2: only receive emergency alert emails, like jailed/inactive, etc.
## 1、Roadmap
- Alarm channels will be added, like Telegram, Discord;
- Monitoring of stake income will be added;
- Program will support other chains developed using the Cosmos SDK;

## 2. Hardware requirements
We recommennd the following hardware resources:
- CPU: 4 cores
- Memory: 8GM RAM
- Database: PostgreSQL

## 3. configuration file
- `vim /data/config/conf.yaml`
- Configuration example
```yaml
# database config
postgres:
  user: "database_user"
  password: "database_password"
  name: "database_name"
  host: "database_host_ip"
  port: "database_host_port"

# cosmos gRPC
cosmos:
  ip: "cosmos_grpc_host_ip"
  gRPCport: "cosmos_grpc_host_grpc_port"
  httpPort: "cosmos_grpc_host_http_port"

alert:
  operatorAddr: "cosmos_validator_operator_address" # The address starts with "cosmosvaloper" And multiple validators are divided by ",". 
  timeInterval: 600  # Monitoring time interval. 600 means the cosmos-monitor runs every 600 seconds
  blockInterval: 100 # Used to calculate the recent signature rate. 100 means to count the signatures rate of the last 100 blocks
  proportion: 0.05 # Signature rate
  startingBlockHeight: 12495683 # The starting block height of the cosmos-monitor program

mail:
  host: "mail_host"
  port: "mail_port"
  username: "mail_username"
  password: "mail_password"
  sender: "sender_mail_address"
  receiver1: "receiver1_mail_address" # Receive all alert emails
  receiver2: "receiver2_mail_address" # Only receive emergency alert emails, like jailed/inactive, etc.

log:
  path: "log_save_location"
  level: "log_level" # Info, Error
  eventlogpath: "event_log_save_location"
```
## 4. Deployment monitoring
- 4.1 Install go
> TIP
> 
> go 1.18++ is required for building and installing the monitor software.

[Golang Official website](https://go.dev/doc/install)

- 4.2 Build monitor binary file
```shell
git clone https://github.com/HashQuark-Research1/cosmos-monitor.git
go build -o cosmos-monitor
```

- 4.3 Edit configuration file
```shell
vim config/conf.yaml
```

- 4.4 Automating your monitor with systemd
```bash
vim /etc/systemd/system/cosmos-monitor.servic
```
```shell
[Unit]
Description=Cosmos Monitor Daemon
After=network.target
[Service]
User=ubuntu
ExecStart=/home/ubuntu/cosmosmonitor/cosmos-monitor -c /data/config/conf.yaml
KillSignal=SIGINT
Restart=on-failure
RestartSec=5s
LimitNOFILE=1000000
[Install]
WantedBy=multi-user.target
[Manager]
DefaultLimitNOFILE=1000000
```

```shell script
systemctl start cosmos-monitor.servic
```